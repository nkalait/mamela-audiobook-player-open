//go:build prod_win

package audio

import (
	"mamela/buildConstraints"
	"mamela/merror"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		merror.ShowError("", err)
	}
	exPath := filepath.Dir(ex)
	exPathArr := strings.Split(exPath, buildConstraints.PathSeparator)
	exPathArr = exPathArr[0 : len(exPathArr)-1]
	LibDir = strings.Join(exPathArr, buildConstraints.PathSeparator) + buildConstraints.PathSeparator + "lib"
	LibExt = ".dll"
}
