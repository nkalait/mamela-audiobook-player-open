package ui

import (
	"fmt"
	"mamela/types"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var bookTitle *canvas.Text
var playingPosition *canvas.Text

func createPlayingLayout(updateNowPlayingChannel chan types.PlayingBook) *fyne.Container {
	initUI()
	playingVBox := container.NewVBox(bookTitle, playingPosition)

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

	sh := fmt.Sprint(h)
	if h < 10 {
		sh = "0" + fmt.Sprint(h)
	}
	sm := fmt.Sprint(m)
	if m < 10 {
		sm = "0" + fmt.Sprint(m)
	}
	ss := fmt.Sprint(s)
	if s < 10 {
		ss = "0" + fmt.Sprint(s)
	}
	playingPosition.Text = sh + " : " + sm + " : " + ss
	playingPosition.Refresh()
}
