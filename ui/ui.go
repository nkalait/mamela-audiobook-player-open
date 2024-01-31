package ui

import (
	"image/color"
	"log"
	"os"
	"path/filepath"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

const rootPath = "/Users/nada/Dev/mamela/books"

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

func BuildUI(appLabel string) {
	mamelaApp := app.New()
	window := mamelaApp.NewWindow(appLabel)

	// button1 := widget.NewButton("click me", func() {
	// 	go func() {
	// 		audio.Stop()
	// 		log.Println("play file 1")
	// 		audio.LoadAndPlay(audioFile1)
	// 	}()
	// })

	// button2 := widget.NewButton("click me", func() {
	// 	go func() {
	// 		audio.Stop()
	// 		log.Println("play file 2")
	// 		audio.LoadAndPlay(audioFile2)
	// 	}()
	// })

	main := container.New(layout.NewHBoxLayout(), createBookList())

	window.SetContent(main)
	window.ShowAndRun()
}

type book struct {
	title    string
	fullPath string
}

func setBookListHeader() string {
	return "Loaded Books"
}

func getAudioBooks() []book {
	var bookList = []book{}
	// d, err := os.Open(dirname)
	rootFolderEntries, err := os.ReadDir(rootPath)
	if err != nil {
		log.Fatal(err)
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
								println(a.fullPath)
								bookList = append(bookList, a)
								break
							}
						}
					}
				}
			}
		}
	}
	return bookList
}

func createBookList() fyne.CanvasObject {
	bookListHeaderTxt := canvas.NewText(setBookListHeader(), textColour)
	bookListHeaderTxt.TextSize = 24
	bookListHeaderTxt.TextStyle.Bold = true
	// audioFile1 := "/Users/nada/Dev/mamela/gunshot.mp3"
	// audioFile2 := "/Users/nada/Dev/mamela/song.mp3"

	books := getAudioBooks()
	bookListVBox := container.New(layout.NewVBoxLayout(), bookListHeaderTxt)
	for _, v := range books {
		log.Println("laying out " + v.fullPath)
		bookTileLayout := NewMyListItemWidget(v)
		bookListVBox.Add(bookTileLayout)
	}
	return bookListVBox
}

// func createPlayingLayout() fyne.Layout {
// 	playingVBox := layout.NewVBoxLayout()

// 	return playingVBox
// }
