package main

//import "github.com/zserge/webview"
import (
	"fmt"

	"github.com/zserge/webview"
)

// func main() {
// 	// Open wikipedia in a 800x600 resizable window
// 	webview.Open("Minimal webview example",
// 		"https://en.m.wikipedia.org/wiki/Main_Page", 800, 600, true)
// }

// Counter is a simple example of automatic Go-to-JS data binding
type Counter struct {
	Value int `json:"value"`
}

// Add increases the value of a counter by n
func (c *Counter) Add(n int) {
	c.Value = c.Value + int(n)
}

// Reset sets the value of a counter back to zero
func (c *Counter) Reset() {
	c.Value = 0
}

func main() {
	w := webview.New(webview.Settings{
		Title: "Click counter: ", // + uiFrameworkName,
	})
	defer w.Exit()

	w.Dispatch(func() {
		// Inject controller
		w.Bind("counter", &Counter{})

		// Inject CSS
		w.InjectCSS(string(MustAsset("src/ui/styles.css")))

		// Inject VueJS
		w.Eval(string(MustAsset("lib/vue/vue.js")))

		// Inject app code
		w.Eval(string(MustAsset("src/ui/app.js")))
	})
	w.Run()

	fmt.Println("Hello")
	asset, err := Asset("src/ui/app.js")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(asset))
}
