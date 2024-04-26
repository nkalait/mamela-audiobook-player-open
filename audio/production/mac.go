//go:build prod_mac

package production

import (
	"mamela/audio"
	"mamela/buildConstraints"
	"mamela/merror"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		merror.ShowError("", e)
	}
	exPath := filepath.Dir(ex)
	exPathArr := strings.Split(exPath, buildConstraints.PathSeparator)
	exPathArr = exPathArr[0 : len(exPathArr)-1]
	audio.LibDir = strings.Join(exPathArr, buildConstraints.PathSeparator) + buildConstraints.PathSeparator + "lib" + buildConstraints.PathSeparator + "mac"
}
