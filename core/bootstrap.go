package core

import (
	"os"
	"path/filepath"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/permission"
	"github.com/xc/digimaker/core/util"
)

func Bootstrap(homePath string) bool {

	if _, err := os.Stat(homePath); os.IsNotExist(err) {
		log.Fatal("Folder " + homePath + " doesn't exist.")
		return false
	}

	abs, _ := filepath.Abs(homePath)
	log.Info("Starting from " + abs)

	util.InitHomePath(homePath)
	err := contenttype.LoadDefinition()
	if err != nil {
		return false
	}

	err = permission.LoadPolicies()
	if err != nil {
		return false
	}
	return true
}

//Initialize db
func InitDB() bool {
	return true
}

func Reload() {

}
