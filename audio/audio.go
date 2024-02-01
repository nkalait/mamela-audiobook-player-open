package audio

import (
	"fmt"
	"log"
	"os"
	"strconv"
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

func (ap *audioPanel) play(resampled *beep.Resampler) {
	// speaker.Play(ap.volume)

	// speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// 	log.Println("done playing")
	// 	done <- true
	// })))
	if resampled != nil {
		speaker.Play(beep.Seq(resampled, beep.Callback(func() {
			log.Println("done playing")
			done <- true
		})))
	} else {
		speaker.Play(beep.Seq(ap.streamer, beep.Callback(func() {
			log.Println("done playing")
			done <- true
		})))
	}
}

func LoadAndPlay(path string) {
	openFileForPlaying(path, &streamer, &format)

	// streamer.Seek(0)

	for {
		select {
		case <-done:
			speaker.Clear()
			playing = false
			return
		case <-time.After(time.Second):
			speaker.Lock()
			// fmt.Println(format.SampleRate.D(streamer.Position()).Round(time.Second))
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
	fmt.Println(1)
	f, err := os.Open(path)
	fmt.Println(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(3)
	var oldSampleRate beep.SampleRate
	fmt.Println(4)
	var resampled *beep.Resampler = nil
	fmt.Println(5)
	if speakerInitialised {
		fmt.Println(6)
		oldSampleRate = format.SampleRate
	}
	fmt.Println(7)
	*streamer, *format, err = mp3.Decode(f)
	fmt.Println(8)
	if err != nil {
		log.Fatal(err)

	}
	fmt.Println(9)

	if !speakerInitialised {
		fmt.Println(10)
		initialiseSpeaker()
	} else {
		fmt.Println("old sample rate: " + strconv.Itoa(int(oldSampleRate)))
		fmt.Println("new sample rate: " + strconv.Itoa(int(format.SampleRate)))
		resampled = beep.Resample(6, format.SampleRate, oldSampleRate, *streamer)
	}

	ap := newAudioPanel(format.SampleRate, *streamer)
	// ******************************
	// set that start position of the stream. should be done before play() otherwise it speeds up
	speaker.Lock()
	newPos := ap.streamer.Position()
	newPos += ap.sampleRate.N(time.Minute * 0)
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
	ap.play(resampled)
	playing = true
}

func initialiseSpeaker() {
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/5))
	// err := speaker.Init(44100, format.SampleRate.N(time.Second/5))
	if err == nil {
		speakerInitialised = true
	}

}
