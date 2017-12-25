package main

//Types of assets
const (
	AssetTypeJS           = `<script src="%s"></script>`
	AssetTypeCSS          = `<link rel="stylesheet" href="%s">`
	assetTypeHTMLTemplate = `<script type="text/x-template" id="%s">%s</script>`
)

//Paths
const (
	FilePathDebugHTML = "src/ui/debug.html"
	URLPathDebug      = "/debug"
)

//Start mode
const (
	StartModeDevServer   = false
	StartModeApplication = true
)
