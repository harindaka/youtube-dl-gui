package main

import (
	"encoding/json"
	"fmt"

	bindata "github.com/jteeuwen/go-bindata"
)

func main() {
	//Hack to keep the dependency github.com/jteeuwen/go-bindata in vendor folder
	var _ = bindata.NewConfig

	goui := NewGoUIApplication(WindowSettings{
		Title:     "Youtube Downloader", // + uiFrameworkName,
		Resizable: true,
		Height:    768,
		Width:     1024,
	})

	registerMessageHandlers(goui)

	goui.DevServerPort = 3030
	goui.Start(StartModeDevServer, registerAssets)
}

func registerAssets(goui *GoUI) {
	// Register ui libraries here (js + css)
	goui.PrependAsset("lib/bootstrap/bootstrap.min.css", AssetTypeCSS)
	goui.PrependAsset("lib/vue/vue.js", AssetTypeJS)

	// Register application specific css assets here
	goui.PrependAsset("src/ui/styles.css", AssetTypeCSS)

	// Register application specific initialization module last
	goui.AppendAsset("src/ui/app.js", AssetTypeJS)
}

func registerMessageHandlers(goui GoUI) {
	goui.OnMessage("add", func(message []byte, send func(string, interface{})) {
		var args map[string]uint
		json.Unmarshal(message, &args)

		val1 := args["val1"]
		val2 := args["val2"]

		send("add", val1+val2)
	})

	goui.OnMessage("getIncText", func(message []byte, send func(string, interface{})) {
		var args map[string]uint
		json.Unmarshal(message, &args)

		val1 := args["val1"]
		val2 := args["val2"]
		result := fmt.Sprintf("Incremented %d by %d", val1, val2)

		send("getIncText", result)
	})
}
