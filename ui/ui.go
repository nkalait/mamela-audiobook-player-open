package ui

import (
	"image/color"
	"os"
	"path/filepath"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var rootPath = ""

var allowedFileTypes = []string{".mp3"}

var (
	colourDarkThemeBlack      = color.RGBA{44, 44, 44, 255}
	colourDarkThemeBlackLight = color.RGBA{67, 67, 67, 255}
	colourDarkThemeWhite      = color.RGBA{214, 214, 214, 255}
)

var (
	textColour    = colourDarkThemeWhite
	BgColour      = colourDarkThemeBlack
	BgColourLight = colourDarkThemeBlackLight
)

func BuildUI(appLabel string, rootP string) {
	rootPath = rootP
	updateBookListChannel := make(chan bool)
	mamelaApp := app.New()
	window := mamelaApp.NewWindow(appLabel)
	bookListVBox := container.New(layout.NewVBoxLayout())
	initBookListPane(bookListVBox, window, updateBookListChannel)
	updateBookList(bookListVBox)
	main := container.New(layout.NewHBoxLayout(), bookListVBox)
	window.SetContent(main)

	go func() {
		for update := range updateBookListChannel {
			if update {
				bookListVBox.Objects = bookListVBox.Objects[:0]
				initBookListPane(bookListVBox, window, updateBookListChannel)
				updateBookList(bookListVBox)
			}
		}
	}()
	window.Resize(fyne.NewSize(600, 600))
	window.ShowAndRun()
}

func initBookListPane(bookListVBox *fyne.Container, window fyne.Window, updateChannel chan bool) {
	bookListHeaderTxt := canvas.NewText(setBookListHeader(), textColour)
	bookListHeaderTxt.TextSize = 24
	bookListHeaderTxt.TextStyle.Bold = true
	top := container.New(layout.NewHBoxLayout(), bookListHeaderTxt, layout.NewSpacer(), layout.NewSpacer(), createFileDialogButton(window, updateChannel))
	bookListVBox.Add(top)
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

type book struct {
	title    string
	fullPath string
}

func setBookListHeader() string {
	return "Loaded Books"
}

func getAudioBooks() ([]book, error) {
	var bookList = []book{}
	rootFolderEntries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, b := range rootFolderEntries {
		if b.IsDir() {
			var bookFullPath = rootPath + "/" + b.Name()
			bookFolder, err := os.ReadDir(bookFullPath)
			if err == nil {
				for _, bookFile := range bookFolder {
					i, err := bookFile.Info()
					if err == nil {
						if i.Mode().IsRegular() {
							if slices.Contains(allowedFileTypes, filepath.Ext(i.Name())) {
								a := book{
									title:    b.Name(),
									fullPath: bookFullPath + "/" + i.Name(),
								}
								bookList = append(bookList, a)
								break
							}
						}
					}
				}
			}
		}
	}
	return bookList, nil
}

func updateBookList(bookListVBox *fyne.Container) {
	books, err := getAudioBooks()
	if err == nil {
		for _, v := range books {
			bookTileLayout := NewMyListItemWidget(v)
			bookListVBox.Add(bookTileLayout)
		}
	}
}

// func createPlayingLayout() fyne.Layout {
// 	playingVBox := layout.NewVBoxLayout()

// 	return playingVBox
// }
