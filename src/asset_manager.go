package main

import (
	"github.com/zserge/webview"
)

type assetManager struct {
	assets map[string]string
}

func newAssetManager() assetManager {
	return assetManager{
		assets: make(map[string]string),
	}
}

func (am *assetManager) addCSS(w webview.WebView, assetPath string) {
	w.InjectCSS(string(MustAsset(assetPath)))
	am.assets[assetPath] = AssetTypeCSS
}

func (am *assetManager) addJS(w webview.WebView, assetPath string) {
	w.Eval(string(MustAsset(assetPath)))
	am.assets[assetPath] = AssetTypeJS
}
