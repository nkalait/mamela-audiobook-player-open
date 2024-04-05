package ui

import (
	"image/color"
	"mamela/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Path to folder where audiobooks are stored
var rootPath = ""

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

var (
	MamelaApp  fyne.App
	MainWindow fyne.Window
)

func BuildUI(appLabel string, rootP string, updateNowPlayingChannel chan types.PlayingBook) {
	rootPath = rootP
	MamelaApp = app.New()
	MainWindow = MamelaApp.NewWindow(appLabel)

	// UI to show list of audiobooks
	leftPane := container.NewBorder(generateBookListContainerTop(), nil, nil, nil, initBookList())

	// Part of UI to place UI elements pertaining to currently playing audiobook
	currentlyPlayingContainer := createPlayingLayout(updateNowPlayingChannel)

	// Place all parts UI parts together
	bodyParts := container.NewBorder(nil, nil, leftPane, nil, currentlyPlayingContainer)

	bodyBg := canvas.NewRectangle(BgColour)
	body := container.NewStack(bodyBg, bodyParts)
	main := container.NewGridWithColumns(1, body)
	MainWindow.SetContent(main)
	MainWindow.Resize(fyne.NewSize(800, 400))
	MainWindow.ShowAndRun()
}
