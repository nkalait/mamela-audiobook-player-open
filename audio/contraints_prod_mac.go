//go:build prod_mac

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
	LibDir = strings.Join(exPathArr, buildconstraints.PathSeparator) + buildconstraints.PathSeparator + "lib" + buildconstraints.PathSeparator + "mac"
	LibExt = ".dylib"
	NotifyInitReady <- true
}
