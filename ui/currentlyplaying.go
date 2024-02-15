package ui

import (
	"mamela/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var updateNowPlayingChannel = make(chan types.PlayingBook)
var bookTitle *canvas.Text

func createPlayingLayout() *fyne.Container {
	initUI()
	playingVBox := container.NewVBox(bookTitle)

	go func() {
		for playingBook := range updateNowPlayingChannel {
			updatePlaying(playingBook)
		}
	}()
	return playingVBox
}

func initUI() {
	initTitle()
}

func initTitle() {
	bookTitle = canvas.NewText("", textColour)
	bookTitle.TextSize = 32
	bookTitle.TextStyle.Bold = true
	bookTitle.Alignment = fyne.TextAlignCenter
}

func updatePlaying(p types.PlayingBook) {
	updateTitle(p.Title)
}

func updateTitle(title string) {
	bookTitle.Text = cases.Title(language.English).String(title)
	bookTitle.Refresh()
}
