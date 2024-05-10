package ui

import (
	"mamela/audio"
	"mamela/storage"
	"mamela/ui/customtheme"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/sqweek/dialog"
)

var (
	mamelaApp  fyne.App
	MainWindow fyne.Window
)

func BuildUI(appLabel string) {
	mamelaApp = app.New()
	setupTheming()
	mainContainer := arrangeUI()
	prepareMainWindow(appLabel, mainContainer)
}

func prepareMainWindow(label string, c *fyne.Container) {
	MainWindow = mamelaApp.NewWindow(label)
	MainWindow.SetContent(c)
	MainWindow.Resize(fyne.NewSize(800, 600))
	MainWindow.CenterOnScreen()
	MainWindow.SetMainMenu(makeMainMenu())
	MainWindow.SetCloseIntercept(func() {
		audio.ExitListener <- true
		time.Sleep(1 * time.Second)
		MainWindow.Close()
	})
	MainWindow.ShowAndRun()
}

func makeMainMenu() *fyne.MainMenu {
	menuItem := fyne.NewMenuItem("Root Folder", func() {
		openSelectRootFolderDialog()
	})
	menu := fyne.NewMenu("File", menuItem)
	mainMenu := fyne.NewMainMenu(menu)
	return mainMenu
}

func openSelectRootFolderDialog() {
	path, err := dialog.Directory().Title("Open root folder").Browse()
	if err != nil {
		dialog.Message(err.Error())
	} else if path != "" {
		audio.ClearCurrentlyPlaying()
		storage.Data.Root = path
		refreshBookList()
		storage.SaveDataToStorageFile()
	}
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
