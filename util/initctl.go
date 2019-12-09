package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/bingoohuang/statiq/fs"
)

// InitCtl initializes the ctl file.
func InitCtl(sfs *fs.StatiqFS, ctlTplName, ctlFilename string) error {
	exists, err := FileStat(ctlFilename)
	if err != nil {
		return err
	}

	if exists == Exists {
		fmt.Println(ctlFilename + " already exists, ignored!")
		return nil
	}

	ctl := string(sfs.Files[ctlTplName].Data)
	tpl, err := template.New(ctlTplName).Parse(ctl)

	if err != nil {
		return err
	}

	binArgs := argsExcludeInit()

	m := map[string]string{"BinName": os.Args[0], "BinArgs": strings.Join(binArgs, " ")}

	var content bytes.Buffer
	if err := tpl.Execute(&content, m); err != nil {
		return err
	}

	// 0755->即用户具有读/写/执行权限，组用户和其它用户具有读写权限；
	if err = ioutil.WriteFile(ctlFilename, content.Bytes(), 0755); err != nil {
		return err
	}

	fmt.Println(ctlFilename + " created!")

	return nil
}

func argsExcludeInit() []string {
	binArgs := make([]string, 0, len(os.Args)-2)

	for i, arg := range os.Args {
		if i == 0 {
			continue
		}

		if strings.Index(arg, "-i") == 0 || strings.Index(arg, "--init") == 0 {
			continue
		}

		if strings.Index(arg, "-") != 0 {
			arg = strconv.Quote(arg)
		}

		binArgs = append(binArgs, arg)
	}

	return binArgs
}
