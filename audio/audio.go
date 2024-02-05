package audio

import (
	"log"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type State struct {
	paused  bool
	playing bool
}

type audioPanel struct {
	format   beep.Format
	streamer beep.StreamSeeker
	ctrl     *beep.Ctrl
	player   State
}

func newAudioPanel(format beep.Format, streamer beep.StreamSeeker) *audioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, streamer)}
	player := State{paused: false, playing: false}
	return &audioPanel{format, streamer, ctrl, player}
}

func (ap *audioPanel) play(resampled *beep.Resampler) {
	speaker.Clear()
	ap.streamer.Seek(0)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		log.Println("done playing")
		done <- true
	})))
	ap.player.playing = true

}

func initialiseSpeaker(format beep.Format) {
	if !speakerInitialised {
		err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err == nil {
			speakerInitialised = true
		}
	}
}

// open media file and start playing it at given start position
func openFileForPlaying(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	var resampled *beep.Resampler = beep.ResampleRatio(6, 1, streamer)
	initialiseSpeaker(format)
	if ap == nil {
		ap = newAudioPanel(format, streamer)
	} else {
		resampled = beep.Resample(6, format.SampleRate, ap.format.SampleRate, streamer)
		ap.streamer = streamer
	}
	startPlaying(2, resampled)
}
