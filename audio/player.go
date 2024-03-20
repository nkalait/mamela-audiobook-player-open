package audio

import (
	"mamela/types"
	"math"
	"time"

	bass "github.com/pteich/gobass"
)

type Player struct {
	channel bass.Channel
}

func (p *Player) Stop() {
	err := p.channel.Stop()
	panicError(err)
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

var exit = make(chan bool)
var stop = make(chan bool)
var player Player

func Initiate() {
	plugins := loadPlugins()
	defer func() {
		for _, p := range plugins {
			bass.PluginFree(p)
		}
		bass.Free()
	}()

	err := bass.Init(-1, 44100, bass.DeviceStereo, 0, 0)
	panicError(err)
	bass.SetVolume(100)
	<-exit
}

func loadPlugins() []uint32 {
	pluginLibbassAac, err := bass.PluginLoad("lib/mac/libbass_aac.dylib", bass.StreamDecode)
	panicError(err)

	plugins := make([]uint32, 2)
	// plugins = append(plugins, pluginLibbass)
	plugins = append(plugins, pluginLibbassAac)

	return plugins
}

func LoadAndPlay(playingBook types.PlayingBook, updateNowPlayingChannel chan types.PlayingBook) {

	// c, err := bass.StreamCreateURL("http://music.myradio.ua:8000/PopRock_news128.mp3", bass.DeviceStereo)
	// c, err := bass.StreamCreateFile("songs/t.m4b", 0, bass.AsyncFile)
	if player.channel != 0 {
		a, e := player.channel.IsActive()
		panicError(e)
		if a == bass.ACTIVE_PLAYING || a == bass.ACTIVE_PAUSED {
			stop <- true
			player.Stop()
		}
	}
	go func() {
		var err error = nil
		player.channel, err = bass.StreamCreateFile(playingBook.FullPath, 0, bass.AsyncFile)
		panicError(err)

		newPos, err := bass.Channel.Seconds2Bytes(player.channel, 0)
		panicError(err)

		bass.Channel.SetPosition(player.channel, newPos, bass.POS_BYTE)
		err = bass.Channel.Play(player.channel, false)
		panicError(err)

	RoutineLoop:
		for {
			select {
			case <-time.After(time.Second):

				active, err := player.channel.IsActive()
				// if err != bass.ERROR_HANDLE {
				panicError(err)
				position, err := player.channel.GetPosition(bass.POS_BYTE)
				panicError(err)
				length, err := player.channel.GetLength(bass.POS_BYTE)
				panicError(err)

				if active == bass.ACTIVE_STOPPED || position == length {
					// println("exit play for loop")
					player.Stop()
					break RoutineLoop
				} else {
					p, err := bass.Channel.Bytes2Seconds(player.channel, position)
					panicError(err)
					// println(fmt.Sprint(p))
					var d time.Duration = time.Duration(math.Round(p * 1000000000))
					// println(fmt.Sprint(d))
					playingBook.Position = time.Duration(d)
					updateNowPlayingChannel <- playingBook
				}
				// }
			case <-stop:
				// println("exit play for loop on <-stop")
				break RoutineLoop
			}
		}
		println("exiting play goroutine")
	}()
}
