package ui

import (
	"mamela/audio"

	"fyne.io/fyne/v2/widget"
)

var playTimeScrubber *widget.Slider
var maxLength = float64(0)
var currentScrubberPosition = float64(0)

func init() {
	playTimeScrubber = widget.NewSlider(0, maxLength)
}

func initPlayTimeScrubberSlider() {
	playTimeScrubber.Orientation = 0 // horizontal

	playTimeScrubber.SetValue(currentScrubberPosition)
	onScrubberDrag()

	go func() {
		for l := range audio.NotifyNewBookLoaded {
			maxLength = l
			playTimeScrubber.Max = maxLength
		}
	}()
	go func() {
		for pos := range audio.NotifyBookPlayTime {
			playTimeScrubber.Value = pos.Seconds()
			playTimeScrubber.Refresh()
		}
	}()
}

func onScrubberDrag() {
	playTimeScrubber.OnChangeEnded = func(f float64) {
		currentScrubberPosition = f
		audio.NotifyBookPlayTimeSliderDragged <- currentScrubberPosition
	}
}

func hidePlayTimeScrubber() {
	playTimeScrubber.Hide()
}

func showPlayTimeScrubber() {
	playTimeScrubber.Show()
}

func adjustPlayTimeScrubberOnKeyPress(keyName string) {
	switch keyName {
	case "Right":
		audio.FastForward()
	case "Left":
		audio.FastRewind()
	}
}
