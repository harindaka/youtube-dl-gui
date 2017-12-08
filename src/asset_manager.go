package main

import (
	"fmt"

	"github.com/zserge/webview"
)

type assetManager struct {
	prependAssets      map[string]string
	prependAssetsIndex []string

	appendAssets      map[string]string
	appendAssetsIndex []string
}

func newAssetManager() assetManager {
	return assetManager{
		prependAssets: make(map[string]string),
		appendAssets:  make(map[string]string),
	}
}

func (am *assetManager) prependAsset(w webview.WebView, assetPath string, assetType string) {
	switch assetType {
	case AssetTypeJS:
		w.Eval(string(MustAsset(assetPath)))
	case AssetTypeCSS:
		w.InjectCSS(string(MustAsset(assetPath)))
	default:
		panic(fmt.Sprintf("Unsupported asset type specified: %s", assetType))
	}

	am.prependAssets[assetPath] = assetType
	am.prependAssetsIndex = append(am.prependAssetsIndex, assetPath)
}

func (am *assetManager) appendAsset(w webview.WebView, assetPath string, assetType string) {
	switch assetType {
	case AssetTypeJS:
		w.Eval(string(MustAsset(assetPath)))
	case AssetTypeCSS:
		w.InjectCSS(string(MustAsset(assetPath)))
	default:
		panic(fmt.Sprintf("Unsupported asset type specified: %s", assetType))
	}

	am.appendAssets[assetPath] = assetType
	am.appendAssetsIndex = append(am.appendAssetsIndex, assetPath)
}

func (am *assetManager) forEachPrependAsset(f func(string, string)) {
	for _, assetPath := range am.prependAssetsIndex {
		assetType := am.prependAssets[assetPath]
		f(assetPath, assetType)
	}
}

func (am *assetManager) forEachAppendAsset(f func(string, string)) {
	for _, assetPath := range am.appendAssetsIndex {
		assetType := am.appendAssets[assetPath]
		f(assetPath, assetType)
	}
}
