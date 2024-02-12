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
	window := mamelaApp.NewWindow(appLabel)

	bookListContainer := initBookList(updateBookListChannel)

	currentlyPlayingContainer := createPlayingLayout()
	bodyParts := container.NewBorder(generateBookListContainerTop(window, updateBookListChannel), nil, bookListContainer, nil, currentlyPlayingContainer)
	bodyBg := canvas.NewRectangle(BgColour)
	body := container.NewStack(bodyBg, bodyParts)
	main := container.NewGridWithColumns(1, body)
	window.SetContent(main)

	window.Resize(fyne.NewSize(600, 300))
	window.ShowAndRun()
}

func initBookList(updateChannel chan bool) *fyne.Container {
	bookListVBox := container.New(layout.NewVBoxLayout())

	bookListContainer := initBookPane(bookListVBox)
	updateBookList(bookListVBox)

	go func() {
		for update := range updateChannel {
			if update {
				bookListVBox.Objects = bookListVBox.Objects[:0]
				updateBookList(bookListVBox)
			}
		}
	}()
	return bookListContainer
}

func setBookListHeader() string {
	return "Loaded Books"
}

func generateBookListContainerTop(window fyne.Window, updateChannel chan bool) *fyne.Container {
	bookListHeaderTxt := canvas.NewText(setBookListHeader(), textColour)
	bookListHeaderTxt.TextSize = 24
	bookListHeaderTxt.TextStyle.Bold = true
	spacer := canvas.NewText("    ", color.Transparent)
	top := container.NewHBox(bookListHeaderTxt, spacer, container.NewVBox(createFileDialogButton(window, updateChannel)))
	return top
}

func initBookPane(bookListVBox *fyne.Container) *fyne.Container {
	// TODO the dots below are just to give the scroller the desired width, NEED TO FIND A WAY TO DO THIS BETTER!!
	dots := canvas.NewText("..........................................................", color.Transparent)
	bookListScroller := container.NewVScroll(bookListVBox)
	bookListVBoxContainerPadded := container.NewPadded(dots, bookListScroller)
	bookListContainer := container.NewStack(canvas.NewRectangle(BgColourLight))
	bookListContainer.Add(bookListVBoxContainerPadded)
	return bookListContainer
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

func createPlayingLayout() *fyne.Container {
	playingVBox := container.NewVBox()
	playingVBox.Add(canvas.NewText("now playing", colourDarkThemeWhite))
	return playingVBox
}
