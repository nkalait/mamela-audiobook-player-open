package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type book struct {
	title    string
	fullPath string
}

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
	updateBookListChannel := make(chan bool)
	mamelaApp := app.New()
	mamelaApp.Settings().SetTheme(mamelaAppTheme())
	window := mamelaApp.NewWindow(appLabel)

	bookListVBox := container.New(layout.NewVBoxLayout())
	bookListScroller := container.NewVScroll(bookListVBox)
	// bookListVBox.Resize(fyne.NewSize(600, 600))
	// bookListScroller.Resize(fyne.NewSize(600, 600))
	bookListVBoxContainerOuter := container.NewStack(canvas.NewRectangle(BgColourLight))
	initBookPane(window, bookListScroller, bookListVBoxContainerOuter, updateBookListChannel)
	updateBookList(bookListVBox)

	main := container.New(layout.NewHBoxLayout(), bookListVBoxContainerOuter)
	window.SetContent(main)

	go func() {
		for update := range updateBookListChannel {
			if update {
				bookListVBox.Objects = bookListVBox.Objects[:0]
				updateBookList(bookListVBox)
			}
		}
	}()
	window.Resize(fyne.NewSize(600, 300))
	window.ShowAndRun()
}

func setBookListHeader() string {
	return "Loaded Books"
}

func generateBookListContainerTop(window fyne.Window, updateChannel chan bool) *fyne.Container {
	bookListHeaderTxt := canvas.NewText(setBookListHeader(), textColour)
	bookListHeaderTxt.TextSize = 24
	bookListHeaderTxt.TextStyle.Bold = true
	top := container.NewHBox(bookListHeaderTxt, container.NewVBox(createFileDialogButton(window, updateChannel)))
	return top
}

func initBookPane(window fyne.Window, bookListScroller fyne.Widget, bookListContainerOuter *fyne.Container, updateChannel chan bool) {
	bookListVBoxContainerTop := generateBookListContainerTop(window, updateChannel)
	// bookListVBoxContainerTop.Resize(fyne.NewSize(250, 40))
	bookListScroller.Resize(fyne.NewSize(210, 400))
	bookListScroller.Move(fyne.NewPos(0, 40))
	bookListVBoxContainerContent := container.NewWithoutLayout(container.NewHBox(bookListVBoxContainerTop), bookListScroller)
	bookListVBoxContainerPadded := container.New(layout.NewPaddedLayout(), bookListVBoxContainerContent)
	bookListContainerOuter.Add(bookListVBoxContainerPadded)
}

func createFileDialogButton(w fyne.Window, updateChannel chan bool) *widget.Button {
	icon := theme.FolderOpenIcon()
	button := widget.NewButtonWithIcon("", icon, func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if dir != nil {
				rootPath = dir.Path()
				updateChannel <- true
			}
		}, w)
	})
	return button
}

// func createPlayingLayout() fyne.Layout {
// 	playingVBox := layout.NewVBoxLayout()

// 	return playingVBox
// }
