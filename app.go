package gostarter

import (
	"embed"
)

//go:embed  web
var Web embed.FS

//go:embed initassets
var InitAssets embed.FS
