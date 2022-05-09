package gostarter

import (
	"embed"
	"github.com/bingoohuang/gg/pkg/emb"
	"github.com/bingoohuang/gg/pkg/logx"
	"github.com/bingoohuang/gostarter/pkg/app"
	"github.com/bingoohuang/gostarter/pkg/ui"
)

//go:embed  web
var web embed.FS

func init() {
	var err error
	ui.WebFS, ui.WebTemplate, err = emb.ParseUI(web, "web", ui.WebFuncMap, true)
	logx.Fatalf(err, "sub web failed: %v", err)
}

func StartWeb() {
	app.StartWeb()
}
