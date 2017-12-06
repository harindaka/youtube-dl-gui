package main

//import "github.com/zserge/webview"
import (
	"fmt"
	"html"
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	bindata "github.com/jteeuwen/go-bindata"
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
	//Hack to keep the dependency github.com/jteeuwen/go-bindata in vendor folder
	var _ = bindata.NewConfig

	go func() {

		http.Handle("/",
			http.FileServer(
				&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}))

		http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})

		fmt.Println("Starting http server...")
		err := http.ListenAndServe(":9090", nil) // set listen port
		if err != nil {
			fmt.Println("ListenAndServe:", err)
		} else {
			fmt.Println("Http server started.")
		}
	}()

	w := webview.New(webview.Settings{
		Title: "Click counter: ", // + uiFrameworkName,
	})
	defer w.Exit()

	w.Dispatch(func() {

		// Register ui libraries here (js + css)
		w.InjectCSS(string(MustAsset("lib/bootstrap/bootstrap.min.css")))
		w.Eval(string(MustAsset("lib/vue/vue.js")))

		// Register application specific css assets here
		w.InjectCSS(string(MustAsset("src/ui/styles.css")))

		// Register application specific utils here
		w.Bind("counter", &Counter{})

		// Register application specific initialization module last
		w.Eval(string(MustAsset("src/ui/app.js")))
	})
	w.Run()

}
