package audio

import (
	"fmt"
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
	pluginLibbassOpus, e := bass.PluginLoad("lib/mac/libbassopus.dylib", bass.StreamDecode)
	err.PanicError(e)

	plugins := make([]uint32, 2)
	plugins = append(plugins, pluginLibbassAac)
	plugins = append(plugins, pluginLibbassOpus)

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
				updateUICurrentlyPlayingInfo()
			case <-exitListener:
				break RoutineLoop
			}
		}
		exitAudio <- true
		time.Sleep(time.Second * 2)
		exitApp <- true
	}()
}

// Pad number below 10 with a zero
func pad(i int) string {
	s := fmt.Sprint(i)
	if i < 10 {
		s = "0" + fmt.Sprint(i)
	}
	return s
}

// Convert seconds to time in hh : mm : ss
func SecondsToTimeText(seconds time.Duration) string {
	var h int = int(seconds.Seconds()) / 3600
	var m int = int(seconds.Seconds()) / 60
	var s int = int(seconds.Seconds()) % 60

	sh := pad(h)
	sm := pad(m)
	ss := pad(s)

	return sh + " : " + sm + " : " + ss
}

// Get the full length in seconds of the currently playing audiobook
func getFullBookLengthSeconds(c bass.Channel) float64 {
	length, e := c.GetLength(bass.POS_BYTE)
	if e != nil {
		err.ShowError("Cannot get full length of audio book", e)
	} else {
		t, e := player.channel.Bytes2Seconds(length)
		if e != nil {
			err.ShowError("Cannot get full length of audio book", e)
		} else {
			return t
		}
	}
	return 0
}

// Update the currently playing audiobook information on the UI
func updateUICurrentlyPlayingInfo() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		err.PanicError(e)
		if active == bass.ACTIVE_PLAYING {
			bytePosition, e := player.channel.GetPosition(bass.POS_BYTE)
			err.PanicError(e)

			p, e := player.channel.Bytes2Seconds(bytePosition)
			err.PanicError(e)
			var d time.Duration = time.Duration(math.Round(p * 1000000000))
			player.currentBook.Position = time.Duration(d)
		} else if active == bass.ACTIVE_STOPPED {
			player.currentBook.Position = 0
		}
		player.updater <- player.currentBook
	}
}

// Start playing a selected audiobook
func LoadAndPlay(playingBook types.PlayingBook) {
	// c, e := bass.StreamCreateURL("http://music.myradio.ua:8000/PopRock_news128.mp3", bass.DeviceStereo)
	player.currentBook = playingBook
	stopPlayingIfPlaying(player.channel, player)

	chapter := player.currentBook.CurrentChapter
	e := loadAudioBookFile(player.currentBook.FullPath + "/" + player.currentBook.Chapters[chapter])
	if e == nil {
		startPlaying()
	}

	updateUICurrentlyPlayingInfo()
}

func stopPlayingIfPlaying(c bass.Channel, p Player) {
	if c != 0 {
		a, e := c.IsActive()
		err.PanicError(e)
		if a == bass.ACTIVE_PLAYING || a == bass.ACTIVE_PAUSED {
			p.stop()
		}
	}
}

func loadAudioBookFile(fullPath string) error {
	var e error = nil
	player.channel, e = bass.StreamCreateFile(fullPath, 0, bass.AsyncFile)
	if e != nil {
		err.ShowError("There seems to be a problem loading the the audiobook file(s)", e)
	}
	return e
}

func startPlaying() error {
	e := player.channel.SetPosition(0, bass.POS_BYTE)
	if e != nil {
		err.ShowError("There seems to be a problem playing the the audiobook file(s)", e)
	} else {
		player.currentBook.FullLengthSeconds = getFullBookLengthSeconds(player.channel)
		player.play()
	}
	return e
}
