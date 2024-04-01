package ui

import (
	"mamela/audio"
	"mamela/types"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var bookTitle *canvas.Text
var bookFullLength *canvas.Text
var playingPosition *canvas.Text
var playerButtonPlay *widget.Button
var playerButtonPause *widget.Button
var playerButtonStop *widget.Button
var playerButtonFastRewind *widget.Button
var playerButtonFastForward *widget.Button
var playerButtonSkipNext *widget.Button
var playerButtonSkipPrevious *widget.Button

func createPlayingLayout(updateNowPlayingChannel chan types.PlayingBook) *fyne.Container {
	initUI()
	playingVBox := container.NewVBox(
		bookTitle,
		playingPosition,
		layoutPlayerButtons(),
		bookFullLength,
	)

	go func() {
		for playingBook := range updateNowPlayingChannel {
			updatePlaying(playingBook)
		}
	}()
	return playingVBox
}

func initUI() {
	initTitle()
	initPlayingPosition()
	initPlayerButtons()
	initFullBookLength()
}

func initTitle() {
	bookTitle = canvas.NewText("", textColour)
	bookTitle.TextSize = 32
	bookTitle.TextStyle.Bold = true
	bookTitle.Alignment = fyne.TextAlignCenter
}

func initFullBookLength() {
	bookFullLength = canvas.NewText("", textColour)
	bookFullLength.TextSize = 32
	bookFullLength.TextStyle.Bold = true
	bookFullLength.Alignment = fyne.TextAlignCenter
}

func initPlayingPosition() {
	playingPosition = canvas.NewText("", textColour)
	playingPosition.TextSize = 24
	playingPosition.Alignment = fyne.TextAlignCenter
}

func initPlayerButtons() {
	playerButtonPlay = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		audio.Play()
	})
	playerButtonPause = widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
		audio.Pause()
	})
	playerButtonStop = widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
		audio.Stop()
	})
	playerButtonFastRewind = widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		audio.FastRewind()
	})
	playerButtonFastForward = widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		audio.FastForward()
	})
	playerButtonSkipNext = widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
		// audio.SkipNext()
	})
	playerButtonSkipPrevious = widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		// audio.SkipPrevious()
	})
}

func layoutPlayerButtons() *fyne.Container {
	layout := container.NewHBox(
		playerButtonSkipPrevious,
		playerButtonFastRewind,
		playerButtonPause,
		playerButtonStop,
		playerButtonPlay,
		playerButtonFastForward,
		playerButtonSkipNext,
	)
	return container.NewCenter(layout)
}

func updatePlaying(p types.PlayingBook) {
	updateTitle(p.Title)
	updatePlayingPosition(p.Position)

	d := time.Duration(math.Round(p.FullLengthSeconds * 1000000000))
	updateBookFullLength(audio.SecondsToTimeText(d))
}

func updateTitle(title string) {
	bookTitle.Text = cases.Title(language.English).String(title)
	bookTitle.Refresh()
}

func updateBookFullLength(bookLength string) {
	bookFullLength.Text = bookLength
	bookFullLength.Refresh()
}

func updatePlayingPosition(p time.Duration) {
	// var h int = int(p.Seconds()) / 3600
	// var m int = int(p.Seconds()) / 60
	// var s int = int(p.Seconds()) % 60

	// sh := pad(h)
	// sm := pad(m)
	// ss := pad(s)

	// playingPosition.Text = sh + " : " + sm + " : " + ss
	playingPosition.Text = audio.SecondsToTimeText(p)
	playingPosition.Refresh()
}
