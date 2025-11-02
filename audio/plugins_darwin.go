//go:build darwin

package audio

import (
	"fmt"
	bass "github.com/pteich/gobass"
	"mamela/merror"
	"mamela/buildconstraints"
)

// Load plugins needed by Bass
func loadPlugins() []uint32 {
	opusPath := LibDir + buildconstraints.PathSeparator + "libbassopus" + LibExt
	fmt.Println(opusPath)
	pluginLibbassOpus, err := bass.PluginLoad(opusPath, bass.StreamDecode)
	merror.PanicError(err)

	plugins := make([]uint32, 1)
	plugins = append(plugins, pluginLibbassOpus)

	return plugins
}
