//go:build prod_mac

package production

import (
	"mamela/audio"
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
	audio.LibDir = strings.Join(exPathArr, buildconstraints.PathSeparator) + buildconstraints.PathSeparator + "lib" + buildconstraints.PathSeparator + "mac"
}
