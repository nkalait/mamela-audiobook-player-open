package audio

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

var done = make(chan bool)
var playing = false

var speakerInitialised = false
var streamer beep.StreamSeekCloser
var format beep.Format

func LoadAndPlay(path string) {
	openFileForPlaying(path, &streamer, &format)
	playing = true
	streamer.Seek(0)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		log.Println("done playing")
		done <- true
	})))
	for {
		select {
		case <-done:
			speaker.Clear()
			playing = false
			return
		case <-time.After(time.Second):
			speaker.Lock()
			fmt.Println(format.SampleRate.D(streamer.Position()).Round(time.Second))
			speaker.Unlock()
		}
	}
}

func Stop() {
	if playing {
		done <- true
	}
}

func openFileForPlaying(path string, streamer *beep.StreamSeekCloser, format *beep.Format) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	*streamer, *format, err = mp3.Decode(f)
	if err != nil {
		log.Fatal(err)

	}

	if !speakerInitialised {
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			log.Fatal(err)
		}
		speakerInitialised = true
	}
}
