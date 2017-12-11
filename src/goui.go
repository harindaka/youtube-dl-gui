package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/zserge/webview"
)

//GoUI plugin
type GoUI struct {
	messageHandlers map[string]func([]byte)
}

//NewNative creates a new Counter plugin
func newGoUI() GoUI {
	return GoUI{
		messageHandlers: make(map[string]func([]byte)),
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
func (g *GoUI) Send(wv webview.WebView, messageType string, message interface{}) error {
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

	js := fmt.Sprintf("goui.invokeJsMessageHandler('%s', '%s');", messageType, string(serializedMessage))
	wv.Eval(js)

	return nil
}

//nativeResult sends a message
func nativeResult(result interface{}) {
	jsMethodName := toLowerCamelCase(getCallingFunctionName())

	var js string
	stringResult, isString := result.(string)

	if isString {
		stringResult = strings.Replace(stringResult, "\\", "\\\\", -1)
		stringResult = strings.Replace(stringResult, "'", "\\'", -1)
		js = fmt.Sprintf("goui.invokeJsMessageHandler('%s', %s);", jsMethodName, fmt.Sprintf("'%s'", stringResult))
	} else if reflect.TypeOf(result).Kind() == reflect.Struct {
		js = fmt.Sprintf("goui.invokeJsMessageHandler('%s', %v);", jsMethodName, result)
	} else {
		js = fmt.Sprintf("goui.invokeJsMessageHandler('%s', %v);", jsMethodName, result)
	}

	w.Eval(js)
}

func toLowerCamelCase(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func getCallingFunctionName() string {

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		panic("Failed to determine the calling function")
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		panic("Failed to obtain details of calling function")
	}

	// return its name
	callingFuncName := fun.Name()
	lastDotCharIndex := strings.LastIndex(callingFuncName, ".")
	if lastDotCharIndex >= 0 {
		return callingFuncName[lastDotCharIndex+1 : len(callingFuncName)]
	}

	return callingFuncName
}
