package audio

import (
	"log"
	"mamela/types"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

var ap *audioPanel = nil
var done = make(chan bool)
var speakerInitialised = false

func startPlaying(startPos int, resampled *beep.Resampler, newFormat beep.Format) {
	speaker.Lock()
	newPos := 0
	newPos += ap.format.SampleRate.N(time.Second * time.Duration(startPos))
	// newPos += newFormat.SampleRate.N(time.Minute * time.Duration(startPos))
	log.Println(newPos)
	if newPos < 0 {
		newPos = 0
	}
	if newPos >= ap.streamer.Len() {
		newPos = ap.streamer.Len() - 1
	}
	err := ap.streamer.Seek(newPos)
	if err != nil {
		log.Fatal(err)
	}
	speaker.Unlock()
	ap.play(resampled)
}

func stop() {
	if ap.player.playing {
		done <- true
	}
}

func LoadAndPlay(playingBook types.PlayingBook, updateNowPlayingChannel chan types.PlayingBook) {
	if ap != nil {
		stop()
	}
	go func() {
		for {
			select {
			case <-done:
				ap.player.playing = false
				log.Println("stop select")
				return
			case <-time.After(time.Second):
				speaker.Lock()
				if ap != nil {
					playingBook.Position = ap.format.SampleRate.D(ap.streamer.Position()).Round(time.Second)
					updateNowPlayingChannel <- playingBook
					// log.Println(playingBook.Position)
				}
				speaker.Unlock()
			}

		}
	}()
	openFileForPlaying(playingBook.FullPath)
}
