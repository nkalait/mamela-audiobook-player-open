package ui

import (
	"mamela/ui/customtheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var (
	mamelaApp  fyne.App
	MainWindow fyne.Window
)

func BuildUI(appLabel string) {
	mamelaApp = app.New()
	setupTheming()
	arrangeUI()
	mainContainer := arrangeUI()
	prepareMainWindow(appLabel, mainContainer)
}

func prepareMainWindow(label string, c *fyne.Container) {
	MainWindow = mamelaApp.NewWindow(label)
	MainWindow.SetContent(c)
	MainWindow.Resize(fyne.NewSize(800, 600))
	MainWindow.CenterOnScreen()
	MainWindow.ShowAndRun()
}

func setupTheming() {
	if customtheme.IsLight() {
		mamelaApp.Settings().SetTheme(customtheme.LightTheme())
	} else {
		mamelaApp.Settings().SetTheme(customtheme.DarkTheme())
	}
}

func arrangeUI() *fyne.Container {
	// UI to show list of audio books
	leftPane := container.NewStack(
		canvas.NewRectangle(customtheme.GetColour(customtheme.ColourNameBackgroundLight)),
		container.NewPadded(container.NewBorder(
			generateBookListContainerTop(),
			nil, nil, nil,
			initBookList(),
		)),
	)

	// Part of UI to place UI elements pertaining to currently playing audio book
	currentlyPlayingContainer := createPlayingLayout()

	// Place all parts UI parts together
	bodyParts := container.NewBorder(
		nil,
		nil,
		leftPane,
		nil,
		currentlyPlayingContainer,
	)

	return container.NewGridWithColumns(1, bodyParts)
}
