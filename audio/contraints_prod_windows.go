//go:build prod_win

package audio

import (
	"mamela/buildconstraints"
	"mamela/merror"
	"os"
	"path/filepath"
)

func init() {
	exPath, err := filepath.Abs(filepath.Dir(os.Args[0])) //get the current working directory
	if err != nil {
		merror.ShowError("Cannot get executable file path", err)
	}
	LibDir = exPath + buildconstraints.PathSeparator + "lib"
	LibExt = ".dll"
	NotifyInitReady <- true
}
