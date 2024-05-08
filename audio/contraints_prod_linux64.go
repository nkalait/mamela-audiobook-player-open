//go:build prod_linux64

package audio

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
		merror.ShowError("", err)
	}
	exPath := filepath.Dir(ex)
	exPathArr := strings.Split(exPath, buildconstraints.PathSeparator)
	exPathArr = exPathArr[0 : len(exPathArr)-1]
	LibDir = strings.Join(exPathArr, buildconstraints.PathSeparator) + buildconstraints.PathSeparator + "lib"
	LibExt = ".so"
}
