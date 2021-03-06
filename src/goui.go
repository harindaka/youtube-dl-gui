package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/websocket"
	"github.com/zserge/webview"
)

//WSMessage represents a Websocket Message
type WSMessage struct {
	MessageType        string `json:"messageType"`
	StringifiedMessage string `json:"stringifiedMessage"`
	CallbackID         string `json:"callbackId"`
}

//DebugHTMLTemplateModel is the debug html template model
type DebugHTMLTemplateModel struct {
	DevServerPort uint
	PrependAssets []UIAsset
	AppendAssets  []UIAsset
	UIJS          string
}

//UIAsset is the model for a ui asset
type UIAsset struct {
	AssetPath      string
	FormattedAsset string
}

//WindowSettings represents the settings for the goui window
type WindowSettings struct {
	// WebView main window title
	Title string
	// URL to open in a webview
	URL string
	// Window width in pixels
	Width int
	// Window height in pixels
	Height int
	// Allows/disallows window resizing
	Resizable bool
}

//GoUI plugin
type GoUI struct {
	DevServerPort uint

	messageHandlers map[string]func([]byte, func(string, interface{}))

	wv         webview.WebView
	wvSettings webview.Settings
	startMode  bool

	prependAssets      map[string]string
	prependAssetsIndex []string
	appendAssets       map[string]string
	appendAssetsIndex  []string
}

//NewGoUIApplication creates a new Counter plugin
func NewGoUIApplication(w WindowSettings) GoUI {
	return GoUI{
		DevServerPort: 3030,

		messageHandlers: make(map[string]func([]byte, func(string, interface{}))),

		wvSettings: webview.Settings{
			Title:     w.Title,
			URL:       w.URL,
			Width:     w.Width,
			Height:    w.Height,
			Resizable: w.Resizable,
		},
		startMode: StartModeApplication,

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
			g.wv.Bind("goui", &(*g))
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

	g.listenHTTP()
}

func (g *GoUI) listenHTTP() {
	fileServerHostAddress := fmt.Sprintf(":%d", g.DevServerPort)
	fmt.Printf("Debug server listening on http://localhost%s%s\n", fileServerHostAddress, URLPathDebug)
	err := http.ListenAndServe(fileServerHostAddress, nil)
	if err != nil {
		panic(fmt.Sprintf("Unable to start file server due to error: %s", err))
	}
}

func (g *GoUI) listenWS(con *websocket.Conn) {
	for {
		wsMessageType, messageBuffer, err := con.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		} else {
			if wsMessageType == websocket.TextMessage {
				var w WSMessage
				json.Unmarshal(messageBuffer, &w)

				if w.MessageType == "" {
					panic(fmt.Sprintf("Invalid messageType field in received websocket message: %s", string(messageBuffer)))
				}

				if w.CallbackID == "" {
					panic(fmt.Sprintf("No callbackId field found in received websocket message: %s", string(messageBuffer)))
				}

				if w.StringifiedMessage == "" {
					panic(fmt.Sprintf("No stringifiedMessage field found in received websocket message: %s", string(messageBuffer)))
				}

				fmt.Println("Received message:", w.StringifiedMessage)

				g.jsToGoCore(w.MessageType, w.StringifiedMessage, w.CallbackID, con)
			}
		}
	}
}

//PrependAsset prepends an asset in the HTML header
func (g *GoUI) PrependAsset(assetPath string, assetType string) {
	formattedAsset := ""

	switch assetType {
	case AssetTypeJS:
		formattedAsset = g.injectJSAsset(assetPath, assetType)
	case AssetTypeCSS:
		formattedAsset = g.injectCSSAsset(assetPath, assetType)
	default:
		panic(fmt.Sprintf("Unsupported asset type specified: %s", assetType))
	}

	g.prependAssets[assetPath] = formattedAsset
	g.prependAssetsIndex = append(g.prependAssetsIndex, assetPath)
}

//AppendAsset appends an asset in the HTML body
func (g *GoUI) AppendAsset(assetPath string, assetType string) {
	formattedAsset := ""

	switch assetType {
	case AssetTypeJS:
		formattedAsset = g.injectJSAsset(assetPath, assetType)
	case AssetTypeCSS:
		formattedAsset = g.injectCSSAsset(assetPath, assetType)
	default:
		panic(fmt.Sprintf("Unsupported asset type specified: %s", assetType))
	}

	g.appendAssets[assetPath] = formattedAsset
	g.appendAssetsIndex = append(g.appendAssetsIndex, assetPath)
}

//AppendHTMLTemplate appends a HTML Template
func (g *GoUI) AppendHTMLTemplate(assetPath string, id string) {
	assetContent := string(MustAsset(assetPath))

	if g.startMode == StartModeApplication {
		js := fmt.Sprintf(`
			(function(){
				var templateHtml = goui.toES5MultilineString(function() {/**
					%s
				**/});
				goui.appendHtmlTemplate(%s, templateHtml);
			})();`, assetContent, toJsString(id))
		g.wv.Eval(js)
	}

	g.appendAssets[assetPath] = fmt.Sprintf(assetTypeHTMLTemplate, id, assetContent)
	g.appendAssetsIndex = append(g.appendAssetsIndex, assetPath)
}

func (g *GoUI) injectJSAsset(assetPath string, assetType string) string {
	if g.startMode == StartModeApplication {
		g.wv.Eval(string(MustAsset(assetPath)))
	}

	return fmt.Sprintf(assetType, assetPath)
}

func (g *GoUI) injectCSSAsset(assetPath string, assetType string) string {
	if g.startMode == StartModeApplication {
		g.wv.InjectCSS(string(MustAsset(assetPath)))
	}

	return fmt.Sprintf(assetType, assetPath)
}

func (g *GoUI) getGoUIJS() string {
	return parseTemplate("templates/goui/goui.js", nil)
}

func assetsToArray(assets map[string]string, index []string) []UIAsset {
	var assetsArray []UIAsset

	for _, assetPath := range index {
		formattedAsset := assets[assetPath]
		assetsArray = append(assetsArray, UIAsset{
			AssetPath:      assetPath,
			FormattedAsset: formattedAsset,
		})
	}

	return assetsArray
}

//OnMessage registers a message handler
func (g *GoUI) OnMessage(messageType string, messageHandler func([]byte, func(string, interface{}))) {
	g.messageHandlers[messageType] = messageHandler
}

//JsToGo triggers the message handler
func (g *GoUI) JsToGo(messageType string, message string, callbackID string) {
	g.jsToGoCore(messageType, message, callbackID, nil)
}

func (g *GoUI) jsToGoCore(messageType string, message string, callbackID string, con *websocket.Conn) {
	handler, ok := g.messageHandlers[messageType]
	if ok {
		handler([]byte(message), func(messageType string, message interface{}) {
			m, err := json.Marshal(message)
			if err != nil {
				fmt.Print(err)
			} else {
				g.send(messageType, string(m), callbackID, con)
			}
		})
	}
}

func (g *GoUI) send(messageType string, stringifiedMessage string, callbackID string, con *websocket.Conn) error {
	if g.startMode == StartModeApplication {
		js := fmt.Sprintf("goui.goToJs(%s, %s, %s);", toJsString(messageType), toJsString(string(stringifiedMessage)), toJsString(callbackID))
		g.wv.Eval(js)
	} else {
		w := WSMessage{
			MessageType:        messageType,
			StringifiedMessage: stringifiedMessage,
			CallbackID:         callbackID,
		}

		messageBuffer, jsonErr := json.Marshal(w)
		if jsonErr != nil {
			panic(jsonErr)
		}

		fmt.Println("Sending message: ", w.StringifiedMessage)
		if wsErr := con.WriteMessage(websocket.TextMessage, messageBuffer); wsErr != nil {
			panic(wsErr)
		}
	}

	return nil
}

func toJsString(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "'", "\\'", -1)
	return fmt.Sprintf("'%s'", s)
}
