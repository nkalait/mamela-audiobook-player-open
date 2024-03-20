package audio

import (
	"fmt"
	"mamela/types"
	"time"

	bass "github.com/pteich/gobass"
)

var pos = make(chan int)
var end = make(chan bool)
var input = make(chan string)

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadPlugins() []uint32 {
	// pluginLibbass, err := bass.PluginLoad("libbass.dylib", bass.StreamDecode)
	// panicError(err)
	pluginLibbassAac, err := bass.PluginLoad("libbass_aac.dylib", bass.StreamDecode)
	panicError(err)

	plugins := make([]uint32, 2)
	// plugins = append(plugins, pluginLibbass)
	plugins = append(plugins, pluginLibbassAac)

	return plugins
}

func LoadAndPlay(playingBook types.PlayingBook, updateNowPlayingChannel chan types.PlayingBook) {
	plugins := loadPlugins()
	defer func() {
		for _, p := range plugins {
			bass.PluginFree(p)
		}
	}()

	err := bass.Init(-1, 44100, bass.DeviceStereo, 0, 0)
	if err != nil {
		panic(err)
	}

	// c, err := bass.StreamCreateURL("http://music.myradio.ua:8000/PopRock_news128.mp3", bass.DeviceStereo)
	c, err := bass.StreamCreateFile("songs/t.m4b", 0, bass.AsyncFile)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		err := bass.Channel.Play(c, false)
		if err != nil {
			panic(err)
		}
		println("here")
		for {
			select {
			case <-time.After(time.Second):
				pb, err := c.GetPosition(bass.POS_BYTE)
				if err != nil {
					panic(err)
				}
				p, err := bass.Channel.Bytes2Seconds(c, pb)
				if err != nil {
					panic(err)
				}
				println(fmt.Sprint(p))
			}
		}
	}()

	bass.SetVolume(50)
	<-end
	println("exit")
}

// func getInput(c bass.Channel) {
// 	var in string
// 	for {
// 		fmt.Scan(&in)
// 		// println(in)
// 		switch in {
// 		case "n":
// 			{
// 				newPos, err := bass.Channel.Seconds2Bytes(c, 60)
// 				if err != nil {
// 					panic(err)
// 				}
// 				bass.Channel.SetPosition(c, newPos, bass.POS_BYTE)
// 			}
// 		case "p":
// 			{
// 				c.Pause()
// 			}
// 		case "g":
// 			{
// 				c.Play(false)
// 			}
// 		}
// 	}
// }
