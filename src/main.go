package main

import (
	"encoding/json"
	"fmt"
	"os"

	bindata "github.com/jteeuwen/go-bindata"
	"github.com/zserge/webview"
)

var isDebugging = false
var goui GoUI

func main() {
	//Hack to keep the dependency github.com/jteeuwen/go-bindata in vendor folder
	var _ = bindata.NewConfig

	isDebugging = (len(os.Args) >= 2 && os.Args[1] == "debug")

	goui = newGoUI(webview.Settings{
		Title:     "Youtube Downloader", // + uiFrameworkName,
		Resizable: true,
		Debug:     isDebugging,
		Height:    768,
		Width:     1024,
	})

	if isDebugging {
		go func() {
			goui.StartDevServer(3030)
		}()
	}

	goui.StartApplication(func() {
		// Register ui libraries here (js + css)
		goui.PrependAsset("lib/bootstrap/bootstrap.min.css", AssetTypeCSS)
		goui.PrependAsset("lib/vue/vue.js", AssetTypeJS)

		// Register application specific css assets here
		goui.PrependAsset("src/ui/styles.css", AssetTypeCSS)

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
}
