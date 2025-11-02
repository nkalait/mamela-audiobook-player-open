//go:build windows

package audio

import "mamela/buildconstraints"

var LibDir = "lib" + buildconstraints.PathSeparator + "win32"
var LibExt = ".dll"
