package main

import (
	"fmt"
	"log"

	"github.com/bingoohuang/gg/pkg/emb"
	"github.com/bingoohuang/gg/pkg/flagparse"
	"github.com/bingoohuang/gg/pkg/jsoni"
	"github.com/bingoohuang/gg/pkg/logx"
	"github.com/bingoohuang/gg/pkg/sigx"
	"github.com/bingoohuang/golog"
	"github.com/bingoohuang/gostarter"
	"github.com/bingoohuang/gostarter/pkg/app"
	"github.com/bingoohuang/gostarter/pkg/conf"
	"github.com/bingoohuang/gostarter/pkg/ui"
)

func main() {
	go app.StartWeb()

	// 可以使用 conf.Conf 获取配置参数
	log.Printf("example config duration: %s", conf.Conf.Duration)

	// 演示打印告警级别日志
	err := fmt.Errorf("demo error")
	log.Printf("W! do something failed: %v", err)
	// W! 告警级别
	// D! 调试级别
	// E! 错误级别

	select {}
}

func init() {
	var err error
	ui.WebFS, ui.WebTemplate, err = emb.ParseUI(gostarter.Web, "web", ui.WebFuncMap, true)
	logx.Fatalf(err, "sub web failed: %v", err)

	flagparse.Parse(conf.Conf,
		flagparse.AutoLoadYaml("conf", "conf.yml"),
		flagparse.ProcessInit(&gostarter.InitAssets))
	golog.Setup(golog.Spec(conf.Conf.Log.Spec), golog.Layout(conf.Conf.Log.Layout))
	sigx.RegisterSignalProfile()

	// export GOLOG_STDOUT=true to check the log
	log.Printf("parse configs: %j", jsoni.AsClearJSON(conf.Conf))
}
