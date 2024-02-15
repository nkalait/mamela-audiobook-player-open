package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// type Book struct {
// 	title    string
// 	fullPath string
// }

// type PlayingBook struct {
// 	Book
// 	position uint16
// }

var rootPath = ""

var allowedFileTypes = []string{".mp3"}

var (
	colourDarkThemeBlack        = color.RGBA{33, 33, 33, 255}
	colourDarkThemeBlackLight   = color.RGBA{38, 38, 38, 255}
	colourDarkThemeBlackLighter = color.RGBA{51, 51, 51, 255}
	colourDarkThemeWhite        = color.RGBA{205, 205, 205, 255}
)

var (
	textColour      = colourDarkThemeWhite
	BgColour        = colourDarkThemeBlack
	BgColourLight   = colourDarkThemeBlackLight
	BgColourLighter = colourDarkThemeBlackLighter
)

func BuildUI(appLabel string, rootP string) {
	rootPath = rootP
	mamelaApp := app.New()
	window := mamelaApp.NewWindow(appLabel)

	bookListContainer := initBookList()

	currentlyPlayingContainer := createPlayingLayout()
	bodyParts := container.NewBorder(generateBookListContainerTop(window), nil, bookListContainer, nil, currentlyPlayingContainer)
	bodyBg := canvas.NewRectangle(BgColour)
	body := container.NewStack(bodyBg, bodyParts)
	main := container.NewGridWithColumns(1, body)
	window.SetContent(main)

	window.Resize(fyne.NewSize(600, 300))
	window.ShowAndRun()
}
