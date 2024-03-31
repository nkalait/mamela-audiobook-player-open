package ui

import (
	"fmt"
	"mamela/audio"
	"mamela/types"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var bookTitle *canvas.Text
var playingPosition *canvas.Text
var playerButtonPlay *widget.Button
var playerButtonPause *widget.Button
var playerButtonStop *widget.Button

func createPlayingLayout(updateNowPlayingChannel chan types.PlayingBook) *fyne.Container {
	initUI()
	playingVBox := container.NewVBox(bookTitle, playingPosition, playerButtonPlay, playerButtonPause, playerButtonStop)

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
}

func initTitle() {
	bookTitle = canvas.NewText("", textColour)
	bookTitle.TextSize = 32
	bookTitle.TextStyle.Bold = true
	bookTitle.Alignment = fyne.TextAlignCenter
}

func initPlayingPosition() {
	playingPosition = canvas.NewText("", textColour)
	playingPosition.TextSize = 24
	playingPosition.Alignment = fyne.TextAlignCenter
}

func initPlayerButtons() {
	playerButtonPlay = widget.NewButton("Play", func() {
		audio.Play()
	})
	playerButtonPause = widget.NewButton("Pause", func() {
		audio.Pause()
	})
	playerButtonStop = widget.NewButton("Stop", func() {
		audio.Stop()
	})
}

func updatePlaying(p types.PlayingBook) {
	updateTitle(p.Title)
	updatePlayingPosition(p.Position)
}

func updateTitle(title string) {
	bookTitle.Text = cases.Title(language.English).String(title)
	bookTitle.Refresh()
}

func updatePlayingPosition(p time.Duration) {
	var h int = int(p.Seconds()) / 3600
	var m int = int(p.Seconds()) / 60
	var s int = int(p.Seconds()) % 60

	sh := pad(h)
	sm := pad(m)
	ss := pad(s)

	playingPosition.Text = sh + " : " + sm + " : " + ss
	playingPosition.Refresh()
}

func pad(i int) string {
	s := fmt.Sprint(i)
	if i < 10 {
		s = "0" + fmt.Sprint(i)
	}
	return s
}
