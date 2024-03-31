package audio

import (
	"mamela/err"
	"mamela/types"
	"math"
	"time"

	bass "github.com/pteich/gobass"
)

// Event listeners
var (
	exitListener = make(chan bool) // for stopping to listen to playing events
	exitAudio    = make(chan bool) // for unloading audio stuff
)

// Holds data structures important to playing an audiobook
var player Player

// Initiate Bass
func init() {
	go func() {
		plugins := loadPlugins()
		defer func(p []uint32) {
			tearDown(p)
		}(plugins)
		initBass()
		<-exitAudio
	}()
}

// Unload loadded Bass plugins and free all resources used by Bass
func tearDown(plugins []uint32) {
	for _, p := range plugins {
		bass.PluginFree(p)
	}
	bass.Free()
}

// Initialise Bass
func initBass() {
	e := bass.Init(-1, 44100, bass.DeviceStereo, 0, 0)
	err.PanicError(e)
	bass.SetVolume(100)
}

// Load pluggins needed by Bass
func loadPlugins() []uint32 {
	pluginLibbassAac, e := bass.PluginLoad("lib/mac/libbass_aac.dylib", bass.StreamDecode)
	err.PanicError(e)

	plugins := make([]uint32, 2)
	plugins = append(plugins, pluginLibbassAac)

	return plugins
}

// Start listening to audio playing event and exit event
func StartChannelListener(updateNowPlayingChannel chan types.PlayingBook, exitApp chan bool) {
	player.updater = updateNowPlayingChannel
	go func() {
	RoutineLoop:
		for {
			select {
			case <-time.After(time.Second):
				if player.channel != 0 {
					active, e := player.channel.IsActive()
					err.PanicError(e)
					bytePosition, e := player.channel.GetPosition(bass.POS_BYTE)
					err.PanicError(e)

					if active == bass.ACTIVE_PLAYING {
						updateUIPlayingPosition(bytePosition)

					}
				}
			case <-exitListener:
				break RoutineLoop
			}
		}
		exitAudio <- true
		time.Sleep(time.Second * 2)
		exitApp <- true
	}()
}

// Update the currently playing audiobook position
func updateUIPlayingPosition(bytePosition int) {
	p, e := player.channel.Bytes2Seconds(bytePosition)
	err.PanicError(e)
	var d time.Duration = time.Duration(math.Round(p * 1000000000))
	player.currentBook.Position = time.Duration(d)
	player.updater <- player.currentBook
}

// Start playing a selected audiobook
func LoadAndPlay(playingBook types.PlayingBook) {
	// c, e := bass.StreamCreateURL("http://music.myradio.ua:8000/PopRock_news128.mp3", bass.DeviceStereo)
	// c, e := bass.StreamCreateFile("songs/t.m4b", 0, bass.AsyncFile)
	player.currentBook = playingBook

	if player.channel != 0 {
		a, e := player.channel.IsActive()
		err.PanicError(e)
		if a == bass.ACTIVE_PLAYING || a == bass.ACTIVE_PAUSED {
			player.stop()
		}
	}
	var e error = nil
	player.channel, e = bass.StreamCreateFile(player.currentBook.FullPath, 0, bass.AsyncFile)
	err.PanicError(e)

	e = player.channel.SetPosition(0, bass.POS_BYTE)
	err.PanicError(e)

	player.play()
}
