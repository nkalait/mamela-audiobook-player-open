//go:build linux

package audio

import (
	"fmt"
	bass "github.com/pteich/gobass"
	"mamela/merror"
	"mamela/buildconstraints"
)

// Load plugins needed by Bass
func loadPlugins() []uint32 {
	aacPath := LibDir + buildconstraints.PathSeparator + "libbass_aac" + LibExt
	opusPath := LibDir + buildconstraints.PathSeparator + "libbassopus" + LibExt
	fmt.Println(aacPath)
	fmt.Println(opusPath)
	pluginLibbassAac, err := bass.PluginLoad(aacPath, bass.StreamDecode)
	merror.PanicError(err)
	pluginLibbassOpus, err := bass.PluginLoad(opusPath, bass.StreamDecode)
	merror.PanicError(err)

	plugins := make([]uint32, 2)
	plugins = append(plugins, pluginLibbassAac)
	plugins = append(plugins, pluginLibbassOpus)

	return plugins
}
