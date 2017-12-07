package main

//import "github.com/zserge/webview"
import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

var webviewTask = make(chan interface{})
var isDebugging = false
var am = newAssetManager()

func main() {
	//Hack to keep the dependency github.com/jteeuwen/go-bindata in vendor folder
	var _ = bindata.NewConfig

	if len(os.Args) >= 2 && os.Args[1] == "debug" {
		isDebugging = true
	}

	var fileServerPort uint
	fileServerPort = 3030

	if isDebugging && len(os.Args) >= 3 {
		fileServerPortStr := os.Args[2]
		port, err := strconv.ParseUint(fileServerPortStr, 0, 64)
		if err != nil || port < 0 || port > 65535 {
			panic(fmt.Sprintf("Invalid file server port specified %s", fileServerPortStr))
		}

		fileServerPort = uint(port)
	}

	if len(os.Args) >= 2 && os.Args[1] == "debug" {

		go func() {
			launchFileServer(fileServerPort)
		}()

		go func() {
			launchWebview()
			close(webviewTask)
		}()

		<-webviewTask
	} else {
		launchWebview()
	}
}

func launchWebview() {
	w := webview.New(webview.Settings{
		Title:     "Youtube Downloader", // + uiFrameworkName,
		Resizable: true,
		Debug:     isDebugging,
		Height:    768,
		Width:     1024,
	})
	defer w.Exit()

	w.Dispatch(func() {

		// Register ui libraries here (js + css)
		//w.InjectCSS(string(MustAsset("lib/bootstrap/bootstrap.min.css")))
		am.addCSS(w, "lib/bootstrap/bootstrap.min.css")
		//w.Eval(string(MustAsset("lib/vue/vue.js")))
		am.addJS(w, "lib/vue/vue.js")

		// Register application specific css assets here
		//w.InjectCSS(string(MustAsset("src/ui/styles.css")))
		am.addCSS(w, "src/ui/styles.css")

		// Register application specific utils here
		w.Bind("counter", &Counter{})

		// Register application specific initialization module last
		//w.Eval(string(MustAsset("src/ui/app.js")))
		am.addJS(w, "src/ui/app.js")
	})
	w.Run()
}

func launchFileServer(port uint) {
	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}))

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/src/ui/debug.html", 302)
	})

	fileServerHostAddress := fmt.Sprintf(":%d", port)
	fmt.Printf("File server listening on http://localhost%s", fileServerHostAddress)
	err := http.ListenAndServe(fileServerHostAddress, nil) // set listen port
	if err != nil {
		fmt.Printf("Unable to start file server due to error: %s", err)
	}
}
