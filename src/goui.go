package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/websocket"
	"github.com/zserge/webview"
)

type message struct {
	messageType string
	message     interface{}
}

//DebugHTMLTemplateModel is the debug html template model
type DebugHTMLTemplateModel struct {
	Port          uint
	PrependAssets []UIAsset
	AppendAssets  []UIAsset
	UIJS          string
}

//UIAsset is the model for a ui asset
type UIAsset struct {
	AssetPath string
	AssetLink string
}

//GoUI plugin
type GoUI struct {
	messageHandlers map[string]func([]byte)
	wv              webview.WebView
	wvSettings      webview.Settings

	prependAssets      map[string]string
	prependAssetsIndex []string
	appendAssets       map[string]string
	appendAssetsIndex  []string
}

//NewNative creates a new Counter plugin
func newGoUI(s webview.Settings) GoUI {
	wv := webview.New(s)
	defer wv.Exit()

	return GoUI{
		messageHandlers: make(map[string]func([]byte)),
		wv:              wv,
		wvSettings:      s,

		prependAssets: make(map[string]string),
		appendAssets:  make(map[string]string),
	}
}

//GetWebView returns the WebView used in this GoUI object
func (g *GoUI) GetWebView(dispatch func()) webview.WebView {
	return g.wv
}

//GetWebViewSettings returns the settings used to create this GoUI object
func (g *GoUI) GetWebViewSettings(dispatch func()) webview.Settings {
	return g.wvSettings
}

//StartApplication starts a GoUI application
func (g *GoUI) StartApplication(dispatch func()) {
	g.wv.Dispatch(func() {
		g.wv.Bind("goui", g)
		g.wv.Eval(g.getGoUIJS())

		dispatch()
	})
	g.wv.Run()
}

func parseTemplate(templatePath string, model interface{}) string {
	templateContent := string(MustAsset(templatePath))
	t := template.New(templatePath)
	t.Parse(templateContent)

	var parsedBytes bytes.Buffer
	if err := t.Execute(&parsedBytes, model); err != nil {
		panic(err)
	}

	return parsedBytes.String()
}

func (g *GoUI) generateDebugHTML(port uint) string {

	model := DebugHTMLTemplateModel{
		Port:          port,
		AppendAssets:  assetsToArray(g.appendAssets, g.appendAssetsIndex),
		PrependAssets: assetsToArray(g.prependAssets, g.prependAssetsIndex),
		UIJS:          g.getGoUIJS(),
	}

	return parseTemplate("templates/goui/debug.html", model)
}

//StartDevServer runs a dev server on specified port
func (g *GoUI) StartDevServer(port uint) {
	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}))

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		debugHTML := g.generateDebugHTML(port)
		w.Write([]byte(debugHTML))
	})

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != "http://"+r.Host {
			http.Error(w, "Origin not allowed", 403)
			return
		}
		con, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		}

		go listenWS(con)
	})

	fileServerHostAddress := fmt.Sprintf(":%d", port)
	fmt.Printf("Debug server listening on http://localhost%s%s\n", fileServerHostAddress, URLPathDebug)
	err := http.ListenAndServe(fileServerHostAddress, nil) // set listen port
	if err != nil {
		fmt.Printf("Unable to start file server due to error: %s", err)
	}
}

func listenWS(con *websocket.Conn) {
	for {
		messageType, p, err := con.ReadMessage()
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(messageType)
			fmt.Println(string(p))
		}

		if err := con.WriteMessage(messageType, p); err != nil {
			log.Println(err)
		}
	}
}

//PrependAsset prepends an asset in the HTML header
func (g *GoUI) PrependAsset(assetPath string, assetType string) {
	switch assetType {
	case AssetTypeJS:
		g.wv.Eval(string(MustAsset(assetPath)))
	case AssetTypeCSS:
		g.wv.InjectCSS(string(MustAsset(assetPath)))
	default:
		panic(fmt.Sprintf("Unsupported asset type specified: %s", assetType))
	}

	g.prependAssets[assetPath] = assetType
	g.prependAssetsIndex = append(g.prependAssetsIndex, assetPath)
}

//AppendAsset appends an asset in the HTML body
func (g *GoUI) AppendAsset(assetPath string, assetType string) {
	switch assetType {
	case AssetTypeJS:
		g.wv.Eval(string(MustAsset(assetPath)))
	case AssetTypeCSS:
		g.wv.InjectCSS(string(MustAsset(assetPath)))
	default:
		panic(fmt.Sprintf("Unsupported asset type specified: %s", assetType))
	}

	g.appendAssets[assetPath] = assetType
	g.appendAssetsIndex = append(g.appendAssetsIndex, assetPath)
}

func (g *GoUI) getGoUIJS() string {
	return parseTemplate("templates/goui/goui.js", nil)
}

func assetsToArray(assets map[string]string, index []string) []UIAsset {
	var assetsArray []UIAsset

	for _, assetPath := range index {
		assetType := assets[assetPath]
		assetsArray = append(assetsArray, UIAsset{
			AssetPath: assetPath,
			AssetLink: fmt.Sprintf(assetType, assetPath),
		})
	}

	return assetsArray
}

//OnMessage registers a message handler
func (g *GoUI) OnMessage(messageType string, messageHandler func([]byte)) {
	g.messageHandlers[messageType] = messageHandler
}

//InvokeGoMessageHandler triggers the message handler
func (g *GoUI) InvokeGoMessageHandler(messageType string, message string) {
	handler, ok := g.messageHandlers[messageType]
	if ok {
		handler([]byte(message))
	}
}

//Send sends a message
func (g *GoUI) Send(messageType string, message interface{}) error {
	var serializedMessage []byte
	var err error
	if message != nil {
		serializedMessage, err = json.Marshal(message)
		if err != nil {
			return err
		}
	} else {
		serializedMessage = []byte("")
	}

	js := fmt.Sprintf("goui.invokeJsMessageHandler(%s, %s);", toJsString(messageType), toJsString(string(serializedMessage)))
	g.wv.Eval(js)

	return nil
}

func toJsString(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "'", "\\'", -1)
	return fmt.Sprintf("'%s'", s)
}

//nativeResult sends a message
// func nativeResult(result interface{}) {
// 	jsMethodName := toLowerCamelCase(getCallingFunctionName())

// 	var js string
// 	stringResult, isString := result.(string)

// 	if isString {
// 		stringResult = strings.Replace(stringResult, "\\", "\\\\", -1)
// 		stringResult = strings.Replace(stringResult, "'", "\\'", -1)
// 		js = fmt.Sprintf("goui.invokeJsMessageHandler('%s', %s);", jsMethodName, fmt.Sprintf("'%s'", stringResult))
// 	} else if reflect.TypeOf(result).Kind() == reflect.Struct {
// 		js = fmt.Sprintf("goui.invokeJsMessageHandler('%s', %v);", jsMethodName, result)
// 	} else {
// 		js = fmt.Sprintf("goui.invokeJsMessageHandler('%s', %v);", jsMethodName, result)
// 	}

// 	w.Eval(js)
// }

// func toLowerCamelCase(s string) string {
// 	if s == "" {
// 		return ""
// 	}
// 	r, n := utf8.DecodeRuneInString(s)
// 	return string(unicode.ToLower(r)) + s[n:]
// }

// func getCallingFunctionName() string {

// 	// we get the callers as uintptrs - but we just need 1
// 	fpcs := make([]uintptr, 1)

// 	// skip 3 levels to get to the caller of whoever called Caller()
// 	n := runtime.Callers(3, fpcs)
// 	if n == 0 {
// 		panic("Failed to determine the calling function")
// 	}

// 	// get the info of the actual function that's in the pointer
// 	fun := runtime.FuncForPC(fpcs[0] - 1)
// 	if fun == nil {
// 		panic("Failed to obtain details of calling function")
// 	}

// 	// return its name
// 	callingFuncName := fun.Name()
// 	lastDotCharIndex := strings.LastIndex(callingFuncName, ".")
// 	if lastDotCharIndex >= 0 {
// 		return callingFuncName[lastDotCharIndex+1 : len(callingFuncName)]
// 	}

// 	return callingFuncName
// }
