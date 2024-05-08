//go:build prod_linux64

package audio

import (
	"fmt"
	"mamela/buildconstraints"
	"mamela/merror"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	// exPath, err := os.Getwd() //get the current directory using the built-in function
	exPath, err := filepath.Abs(filepath.Dir(os.Args[0])) //get the current working directory
	if err != nil {
		fmt.Println(err) //print the error if obtained
	}
	fmt.Println("Current working directory:", exPath)
	exPathArr := strings.Split(exPath, buildconstraints.PathSeparator)
	if err != nil {
		merror.ShowError("", err)
	}
	exPathArr = exPathArr[0 : len(exPathArr)-1]
	// LibDir = strings.Join(exPathArr, buildconstraints.PathSeparator) + buildconstraints.PathSeparator + "lib"
	LibDir = exPath + buildconstraints.PathSeparator + "lib"
	LibExt = ".so"
	NotifyInitReady <- true
}
