package util

import (
	"fmt"
	"os"

	"github.com/bingoohuang/statiq/fs"
)

// Ipo initializes the ctl.sh and cnf.tpl.toml.
func Ipo(ipo bool) {
	if !ipo {
		return
	}

	if err := ipoInit(); err != nil {
		fmt.Println(err)
	}

	os.Exit(0)
}

func ipoInit() error {
	sfs, err := fs.New()

	if err != nil {
		return err
	}

	if err = InitCtl(sfs, "/ctl.tpl.sh", "./ctl"); err != nil {
		return err
	}

	if err = InitCfgFile(sfs, "/cnf.tpl.toml", "./cnf.toml"); err != nil {
		return err
	}

	return nil
}
