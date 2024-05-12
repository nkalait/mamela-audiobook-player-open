package ui

import (
	"mamela/audio"
	"mamela/storage"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var volumeSlider *widget.Slider

func init() {
	volumeSlider = widget.NewSlider(0, 10000) // 10000 is bass's Stream global volume max level
}

func createVolumeSlider() *fyne.Container {
	volumeSlider.Orientation = 1 // vertical
	// vol := audio.GetVolume()
	onDragVolumeSlider(volumeSlider)
	hBox := container.NewHBox(volumeSlider, container.NewPadded())
	return container.NewPadded(hBox)
}

func onDragVolumeSlider(slider *widget.Slider) {
	slider.OnChangeEnded = func(f float64) {
		audio.NotifyVolumeSliderDragged <- f
		storage.SaveVolumeLevel(f)
	}
}

func hideVolumeSlider() {
	volumeSlider.Hide()
}

func showVolumeSlider() {
	volumeSlider.Show()
	volumeSlider.SetValue(storage.GetVolumeLevel())
	audio.NotifyVolumeSliderDragged <- storage.GetVolumeLevel()
}
