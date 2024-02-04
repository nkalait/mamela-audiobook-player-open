package audio

import (
	"log"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
)

var ap *audioPanel = nil
var done = make(chan bool)
var speakerInitialised = false

func startPlaying(startPos int, resampled *beep.Resampler) {
	ap.play(resampled)
	// speaker.Lock()
	// newPos := 0
	// newPos += ap.format.SampleRate.N(time.Minute * time.Duration(startPos))
	// log.Println(newPos)
	// if newPos < 0 {
	// 	newPos = 0
	// }
	// if newPos >= ap.streamer.Len() {
	// 	newPos = ap.streamer.Len() - 1
	// }
	// err := ap.streamer.Seek(newPos)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// speaker.Unlock()
}

func stop() {
	if ap.player.playing {
		done <- true
	}
}

func LoadAndPlay(path string) {
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
					log.Println(ap.format.SampleRate.D(ap.streamer.Position()).Round(time.Second))
				}
				speaker.Unlock()
			}

		}
	}()
	openFileForPlaying(path)
}
