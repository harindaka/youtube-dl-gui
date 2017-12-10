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

var webviewTask = make(chan interface{})
var isDebugging = false
var am = newAssetManager()
var w webview.WebView

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
	w = webview.New(webview.Settings{
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
		am.prependAsset(w, "lib/bootstrap/bootstrap.min.css", AssetTypeCSS)
		//w.Eval(string(MustAsset("lib/vue/vue.js")))
		am.prependAsset(w, "lib/vue/vue.js", AssetTypeJS)

		// Register application specific css assets here
		//w.InjectCSS(string(MustAsset("src/ui/styles.css")))
		am.prependAsset(w, "src/ui/styles.css", AssetTypeCSS)

		// Register application specific utils here
		//w.Bind("counter", &Counter{})
		//counter := plugins.NewCounter()
		//w.Bind("counter", &Counter{})

		native := newNative()
		w.Bind("native", &native)

		//w.Eval("window.plugins = plugins.data;")

		// Register application specific initialization module last
		//w.Eval(string(MustAsset("src/ui/app.js")))
		am.appendAsset(w, "src/ui/app.js", AssetTypeJS)
	})
	w.Run()
}

func launchFileServer(port uint) {
	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}))

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<!DOCTYPE html>\n"))
		w.Write([]byte("<html>\n"))
		w.Write([]byte("<head>\n"))

		am.forEachPrependAsset(func(assetPath string, assetType string) {
			markup := fmt.Sprintf(assetType, assetPath)
			markup = fmt.Sprintf("%s\n", markup)
			w.Write([]byte(markup))
		})

		w.Write([]byte("<title>Page Title</title>\n"))
		w.Write([]byte("</head>\n"))
		w.Write([]byte("<body>\n"))
		w.Write([]byte("<div id=\"app\"></div>\n"))

		am.forEachAppendAsset(func(assetPath string, assetType string) {
			markup := fmt.Sprintf(assetType, assetPath)
			markup = fmt.Sprintf("%s\n", markup)
			w.Write([]byte(markup))
		})

		w.Write([]byte("</body>\n"))
		w.Write([]byte("</html>"))

		//w.WriteHeader(http.StatusOK)
	})

	fileServerHostAddress := fmt.Sprintf(":%d", port)
	fmt.Printf("Debug server listening on http://localhost%s%s\n", fileServerHostAddress, URLPathDebug)
	err := http.ListenAndServe(fileServerHostAddress, nil) // set listen port
	if err != nil {
		fmt.Printf("Unable to start file server due to error: %s", err)
	}
}
