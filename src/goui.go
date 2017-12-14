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

//DebugHTMLTemplateModel is the debug html template model
type DebugHTMLTemplateModel struct {
	DevServerPort uint
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
	DevServerPort uint

	messageHandlers map[string]func([]byte, func(string, interface{}))
	wv              webview.WebView
	wvSettings      webview.Settings
	startMode       bool

	prependAssets      map[string]string
	prependAssetsIndex []string
	appendAssets       map[string]string
	appendAssetsIndex  []string
}

//NewNative creates a new Counter plugin
func newGoUI(s webview.Settings) GoUI {
	return GoUI{
		DevServerPort: 3030,

		messageHandlers: make(map[string]func([]byte, func(string, interface{}))),
		wvSettings:      s,
		startMode:       StartModeApplication,

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

//Start starts a GoUI application or dev server
func (g *GoUI) Start(startMode bool, registerAssets func(*GoUI)) {
	g.startMode = startMode
	if startMode == StartModeApplication {
		g.wv = webview.New(g.wvSettings)
		defer g.wv.Exit()

		g.wv.Dispatch(func() {
			g.wv.Bind("goui", g)
			g.wv.Eval(g.getGoUIJS())

			registerAssets(g)
		})

		g.wv.Run()
	} else {
		registerAssets(g)
		g.startDevServer()
	}
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

func (g *GoUI) generateDebugHTML() string {

	model := DebugHTMLTemplateModel{
		DevServerPort: g.DevServerPort,
		AppendAssets:  assetsToArray(g.appendAssets, g.appendAssetsIndex),
		PrependAssets: assetsToArray(g.prependAssets, g.prependAssetsIndex),
		UIJS:          g.getGoUIJS(),
	}

	return parseTemplate("templates/goui/debug.html", model)
}

//StartDevServer runs a dev server on specified port
func (g *GoUI) startDevServer() {
	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}))

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		debugHTML := g.generateDebugHTML()
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

		go g.listenWS(con)
	})

	fileServerHostAddress := fmt.Sprintf(":%d", g.DevServerPort)
	fmt.Printf("Debug server listening on http://localhost%s%s\n", fileServerHostAddress, URLPathDebug)
	err := http.ListenAndServe(fileServerHostAddress, nil) // set listen port
	if err != nil {
		fmt.Printf("Unable to start file server due to error: %s", err)
	}
}

func (g *GoUI) listenWS(con *websocket.Conn) {
	for {
		wsMessageType, messageBuffer, err := con.ReadMessage()
		if err != nil {
			log.Println(err)
		} else {
			if wsMessageType == websocket.TextMessage {
				wsMessage := make(map[string]string)
				json.Unmarshal(messageBuffer, &wsMessage)

				var fieldExists bool
				var messageType string
				var callbackId string
				var stringifiedMessage string
				messageType, fieldExists = wsMessage["messageType"]
				if !fieldExists {
					panic(fmt.Sprintf("No messageType field found in received websocket message: %s", string(messageBuffer)))
				}

				callbackId, fieldExists = wsMessage["callbackId"]
				if !fieldExists {
					panic(fmt.Sprintf("No callbackId field found in received websocket message: %s", string(messageBuffer)))
				}

				stringifiedMessage, fieldExists = wsMessage["stringifiedMessage"]
				if !fieldExists {
					panic(fmt.Sprintf("No message field found in received websocket message: %s", string(messageBuffer)))
				}

				fmt.Println("Received message: " + stringifiedMessage)

				g.InvokeGoMessageHandler(messageType, stringifiedMessage, callbackId)
			}
		}

		// if err := con.WriteMessage(websocket.TextMessage, messageBuffer); err != nil {
		// 	log.Println(err)
		// }
	}
}

//PrependAsset prepends an asset in the HTML header
func (g *GoUI) PrependAsset(assetPath string, assetType string) {
	switch assetType {
	case AssetTypeJS:
		g.evalAsset(assetPath)
	case AssetTypeCSS:
		g.injectCSSAsset(assetPath)
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
		g.evalAsset(assetPath)
	case AssetTypeCSS:
		g.injectCSSAsset(string(MustAsset(assetPath)))
	default:
		panic(fmt.Sprintf("Unsupported asset type specified: %s", assetType))
	}

	g.appendAssets[assetPath] = assetType
	g.appendAssetsIndex = append(g.appendAssetsIndex, assetPath)
}

func (g *GoUI) evalAsset(assetPath string) {
	if g.startMode == StartModeApplication {
		g.wv.Eval(string(MustAsset(assetPath)))
	}
}

func (g *GoUI) injectCSSAsset(assetPath string) {
	if g.startMode == StartModeApplication {
		g.wv.InjectCSS(string(MustAsset(assetPath)))
	}
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
func (g *GoUI) OnMessage(messageType string, messageHandler func([]byte, func(string, interface{}))) {
	g.messageHandlers[messageType] = messageHandler
}

//InvokeGoMessageHandler triggers the message handler
func (g *GoUI) InvokeGoMessageHandler(messageType string, message string, callbackID string) {
	handler, ok := g.messageHandlers[messageType]
	if ok {
		handler([]byte(message), func(messageType string, message interface{}) {
			g.send(messageType, message, callbackID)
		})
	}
}

//Send sends a message
func (g *GoUI) send(messageType string, message interface{}, callbackID string) error {
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

	if g.startMode == StartModeApplication {
		js := fmt.Sprintf("goui.invokeJsMessageHandler(%s, %s, %s);", toJsString(messageType), toJsString(string(serializedMessage)), toJsString(callbackID))
		g.wv.Eval(js)
	} else {
		//todo: send message with messageType and callbackid
	}

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
