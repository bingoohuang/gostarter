package main

import (
	"github.com/bingoohuang/gg/pkg/ctl"
	"github.com/bingoohuang/gg/pkg/fla9"
	"github.com/bingoohuang/gg/pkg/sigx"
	"github.com/bingoohuang/golog"
	"github.com/bingoohuang/gostarter"
)

func main() {
	pInit := fla9.Bool("init", false, "Create initial ctl and exit")
	pVersion := fla9.Bool("version", false, "Create initial ctl and exit")
	fla9.Parse()
	ctl.Config{Initing: *pInit, PrintVersion: *pVersion}.ProcessInit()
	golog.Setup()
	sigx.RegisterSignalProfile()

	go gostarter.StartWeb()
	select {}
}
