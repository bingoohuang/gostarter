package util

import (
	"fmt"
	"io/ioutil"

	"github.com/bingoohuang/statiq/fs"
)

func InitConfigFile(sfs *fs.StatiqFS, configTplFileName, configFileName string) error {
	exists, err := FileStat(configFileName)
	if err != nil {
		return err
	}
	if exists == Exists {
		fmt.Printf("%s already exists, ignored!\n", configFileName)
		return nil
	}

	conf := sfs.Files[configTplFileName].Data
	// 0644->即用户具有读写权限，组用户和其它用户具有只读权限；
	if err := ioutil.WriteFile(configFileName, conf, 0644); err != nil {
		return err
	}

	fmt.Println(configFileName + " created!")

	return nil
}
