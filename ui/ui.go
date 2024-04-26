package ui

import (
	"image/color"
	"mamela/ui/customtheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Path to folder where audio books are stored
var rootPath = ""

var (
	// colourDarkThemeBlackLight   = color.RGBA{38, 38, 38, 255}
	colourDarkThemeBlackLighter = color.RGBA{51, 51, 51, 255}
)

var (
	// BgColourLight   = colourDarkThemeBlackLight
	BgColourLighter = colourDarkThemeBlackLighter
)

var (
	MamelaApp  fyne.App
	MainWindow fyne.Window
)

func BuildUI(appLabel string, rootP string) {
	rootPath = rootP
	MamelaApp = app.New()
	// MamelaApp.Settings().SetTheme(theme.DarkTheme())

	// if customtheme.IsDark(fyne.CurrentApp().Settings().Theme()) {
	if customtheme.IsLight() {
		MamelaApp.Settings().SetTheme(customtheme.LightTheme())
	} else {
		MamelaApp.Settings().SetTheme(customtheme.DarkTheme())
	}

	MainWindow = MamelaApp.NewWindow(appLabel)

	// UI to show list of audio books
	leftPane := container.NewStack(
		canvas.NewRectangle(customtheme.GetColour(customtheme.ColourNameBackgroundLight)),
		container.NewPadded(container.NewBorder(generateBookListContainerTop(), nil, nil, nil, initBookList())),
	)

	// Part of UI to place UI elements pertaining to currently playing audio book
	currentlyPlayingContainer := createPlayingLayout()

	// Place all parts UI parts together
	bodyParts := container.NewBorder(nil, nil, leftPane, nil, currentlyPlayingContainer)

	// bodyBg := canvas.NewRectangle(theme.BackgroundColor())
	// body := container.NewStack(bodyBg, bodyParts)
	main := container.NewGridWithColumns(1, bodyParts)
	MainWindow.SetContent(main)
	MainWindow.Resize(fyne.NewSize(800, 400))
	MainWindow.ShowAndRun()
}
