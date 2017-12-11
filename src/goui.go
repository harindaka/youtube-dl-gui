package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

//GoUI plugin
type GoUI struct{}

//NewNative creates a new Counter plugin
func newGoUI() GoUI {
	return GoUI{}
}

//Add increments value
func (c *GoUI) Add(val1 uint, val2 uint) {
	nativeResult(val1 + val2)
}

//GetIncText gets incremented text
func (c *GoUI) GetIncText(incVal uint, incBy uint) {
	nativeResult(fmt.Sprintf("Incremented %v by %v", incVal, incBy))
}

func nativeResult(result interface{}) {
	jsMethodName := toLowerCamelCase(getCallingFunctionName())

	var js string
	stringResult, isString := result.(string)

	if isString {
		stringResult = strings.Replace(stringResult, "\\", "\\\\", -1)
		stringResult = strings.Replace(stringResult, "'", "\\'", -1)
		js = fmt.Sprintf("goui.onMessage('%s', %s);", jsMethodName, fmt.Sprintf("'%s'", stringResult))
	} else if reflect.TypeOf(result).Kind() == reflect.Struct {
		js = fmt.Sprintf("goui.onMessage('%s', %v);", jsMethodName, result)
	} else {
		js = fmt.Sprintf("goui.onMessage('%s', %v);", jsMethodName, result)
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
