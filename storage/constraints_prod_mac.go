//go:build prod_mac

package storage

import (
	"mamela/buildconstraints"
	"mamela/merror"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		merror.ShowError("Cannot get executable path", err)
	}
	exPath := filepath.Dir(ex)
	exPathArr := strings.Split(exPath, buildconstraints.PathSeparator)
	exPathArr = exPathArr[0 : len(exPathArr)-1]

	dir := strings.Join(exPathArr, buildconstraints.PathSeparator) + buildconstraints.PathSeparator + "db"

	StorageFile = dir + buildconstraints.PathSeparator + "data.json"
}
