package audio

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

var done = make(chan bool)
var playing = false

var speakerInitialised = false
var streamer beep.StreamSeekCloser
var format beep.Format

type audioPanel struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
}

func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker) *audioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &audioPanel{sampleRate, streamer, ctrl, resampler, volume}
}

func (ap *audioPanel) play() {
	// speaker.Play(ap.volume)

	// speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// 	log.Println("done playing")
	// 	done <- true
	// })))

	speaker.Play(beep.Seq(ap.streamer, beep.Callback(func() {
		log.Println("done playing")
		done <- true
	})))
}

func LoadAndPlay(path string) {
	openFileForPlaying(path, &streamer, &format)
	ap := newAudioPanel(format.SampleRate, streamer)

	// ******************************
	// set that start position of the stream. should be done before play() otherwise it speeds up
	speaker.Lock()
	newPos := ap.streamer.Position()
	newPos += ap.sampleRate.N(time.Minute * 2)
	if newPos < 0 {
		newPos = 0
	}
	if newPos >= ap.streamer.Len() {
		newPos = ap.streamer.Len() - 1
	}
	if err := ap.streamer.Seek(newPos); err != nil {
		log.Fatal(err)
	}
	speaker.Unlock()
	// ******************************
	ap.play()
	playing = true

	// streamer.Seek(0)

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
