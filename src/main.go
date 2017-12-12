package main

//import "github.com/zserge/webview"
import (
	"encoding/json"
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
var goui GoUI

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

	goui = newGoUI(w)

	w.Dispatch(func() {
		w.Bind("goui", &goui)
		goui.Eval(goui.GetGoUIJS())

		// Register ui libraries here (js + css)
		goui.PrependAsset("lib/bootstrap/bootstrap.min.css", AssetTypeCSS)
		goui.PrependAsset("lib/vue/vue.js", AssetTypeJS)

		// Register application specific css assets here
		//w.InjectCSS(string(MustAsset("src/ui/styles.css")))
		goui.PrependAsset("src/ui/styles.css", AssetTypeCSS)

		// Register application specific utils here
		//w.Bind("counter", &Counter{})
		//counter := plugins.NewCounter()
		//w.Bind("counter", &Counter{})

		goui.OnMessage("add", func(message []byte) {
			var args map[string]uint
			json.Unmarshal(message, &args)

			val1 := args["val1"]
			val2 := args["val2"]

			goui.Send("add", val1+val2)
		})

		goui.OnMessage("getIncText", func(message []byte) {
			var args map[string]uint
			json.Unmarshal(message, &args)

			val1 := args["val1"]
			val2 := args["val2"]
			result := fmt.Sprintf("Incremented %d by %d", val1, val2)

			goui.Send("getIncText", result)
		})

		// Register application specific initialization module last
		goui.AppendAsset("src/ui/app.js", AssetTypeJS)
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

		goui.ForEachPrependAsset(func(assetPath string, assetType string) {
			markup := fmt.Sprintf(assetType, assetPath)
			markup = fmt.Sprintf("%s\n", markup)
			w.Write([]byte(markup))
		})

		w.Write([]byte("<title>Page Title</title>\n"))
		w.Write([]byte("</head>\n"))
		w.Write([]byte("<body>\n"))
		w.Write([]byte("<div id=\"app\"></div>\n"))

		goui.ForEachAppendAsset(func(assetPath string, assetType string) {
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
