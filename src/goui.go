package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zserge/webview"
)

//GoUI plugin
type GoUI struct {
	messageHandlers map[string]func([]byte)
	wv              webview.WebView

	prependAssets      map[string]string
	prependAssetsIndex []string
	appendAssets       map[string]string
	appendAssetsIndex  []string
}

//NewNative creates a new Counter plugin
func newGoUI(wv webview.WebView) GoUI {
	return GoUI{
		messageHandlers: make(map[string]func([]byte)),
		wv:              wv,

		prependAssets: make(map[string]string),
		appendAssets:  make(map[string]string),
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

//ForEachPrependAsset allows iteration of prepended assets
func (g *GoUI) ForEachPrependAsset(f func(string, string)) {
	for _, assetPath := range g.prependAssetsIndex {
		assetType := g.prependAssets[assetPath]
		f(assetPath, assetType)
	}
}

//ForEachAppendAsset allows iteration of appended assets
func (g *GoUI) ForEachAppendAsset(f func(string, string)) {
	for _, assetPath := range g.appendAssetsIndex {
		assetType := g.appendAssets[assetPath]
		f(assetPath, assetType)
	}
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
