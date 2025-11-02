//go:build darwin

package audio

import "mamela/buildconstraints"

var LibDir = "lib" + buildconstraints.PathSeparator + "mac"
var LibExt = ".dylib"
