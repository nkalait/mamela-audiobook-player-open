//go:build linux

package audio

import "mamela/buildconstraints"

var LibDir = "lib" + buildconstraints.PathSeparator + "linux64"
var LibExt = ".so"
