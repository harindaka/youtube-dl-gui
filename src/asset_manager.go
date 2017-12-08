package main

import (
	"fmt"

	"github.com/zserge/webview"
)

type assetManager struct {
	assets      map[string]string
	assetsIndex []string
}

func newAssetManager() assetManager {
	return assetManager{
		assets: make(map[string]string),
	}
}

func (am *assetManager) addAsset(w webview.WebView, assetPath string, assetType string) {
	switch assetType {
	case AssetTypeJS:
		w.Eval(string(MustAsset(assetPath)))
	case AssetTypeCSS:
		w.InjectCSS(string(MustAsset(assetPath)))
	default:
		panic(fmt.Sprintf("Unsupported asset type specified: %s", assetType))
	}

	am.assets[assetPath] = assetType
	am.assetsIndex = append(am.assetsIndex, assetPath)
}

func (am *assetManager) forEachAsset(f func(string, string)) {
	for _, assetPath := range am.assetsIndex {
		assetType := am.assets[assetPath]
		f(assetPath, assetType)
	}
}
