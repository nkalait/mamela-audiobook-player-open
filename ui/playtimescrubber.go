package ui

import (
	"mamela/audio"

	"fyne.io/fyne/v2/widget"
)

var playTimeScrubber *widget.Slider

func init() {
	playTimeScrubber = widget.NewSlider(0, 0)
}

func initPlayTimeScrubberSlider() {
	playTimeScrubber.Orientation = 0 // horizontal

	playTimeScrubber.SetValue(0)
	onScrubberDrag()

	go func() {
		for length := range audio.NotifyNewBookLoaded {
			playTimeScrubber.Max = length
		}
	}()
	go func() {
		for pos := range audio.NotifyBookPlayTime {
			playTimeScrubber.Value = pos.Seconds()
		}
	}()
}

func onScrubberDrag() {
	playTimeScrubber.OnChangeEnded = func(f float64) {
		audio.NotifyBookPlayTimeSliderDragged <- f
	}
}

func hidePlayTimeScrubber() {
	playTimeScrubber.Hide()
}

func showPlayTimeScrubber() {
	playTimeScrubber.Show()
}
