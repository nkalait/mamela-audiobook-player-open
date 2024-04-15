//go:build prod_mac

package storage

import (
	"log"
	"mamela/buildconstraints"
	"mamela/err"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	ex, e := os.Executable()
	if e != nil {
		err.ShowError("", e)
	}
	exPath := filepath.Dir(ex)
	exPathArr := strings.Split(exPath, buildconstraints.PathSeparator)
	exPathArr = exPathArr[0 : len(exPathArr)-1]

	dir := strings.Join(exPathArr, buildconstraints.PathSeparator) + buildconstraints.PathSeparator + "db"

	storageFile = dir + buildconstraints.PathSeparator + "data.json"
	log.Println("initiated storage: " + storageFile)
}
