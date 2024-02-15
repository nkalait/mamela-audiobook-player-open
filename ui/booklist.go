package ui

import (
	"image/color"
	"mamela/types"
	"os"
	"path/filepath"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var updateBookListChannel = make(chan bool)

func initBookList(updateNowPlayingChannel chan types.PlayingBook) *fyne.Container {
	bookListVBox := container.New(layout.NewVBoxLayout())

	bookListContainer := initBookPane(bookListVBox)
	updateBookList(bookListVBox, updateNowPlayingChannel)

	go func() {
		for update := range updateBookListChannel {
			if update {
				bookListVBox.Objects = bookListVBox.Objects[:0]
				updateBookList(bookListVBox, updateNowPlayingChannel)
			}
		}
	}()
	return bookListContainer
}

func setBookListHeader() string {
	return "Loaded Books"
}

func generateBookListContainerTop(window fyne.Window) *fyne.Container {
	bookListHeaderTxt := canvas.NewText(setBookListHeader(), textColour)
	bookListHeaderTxt.TextSize = 24
	bookListHeaderTxt.TextStyle.Bold = true
	spacer := canvas.NewText("    ", color.Transparent)
	top := container.NewHBox(bookListHeaderTxt, spacer, container.NewVBox(createFileDialogButton(window)))
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

func createFileDialogButton(w fyne.Window) *widget.Button {
	icon := theme.FolderOpenIcon()
	button := widget.NewButtonWithIcon("", icon, func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if dir != nil {
				rootPath = dir.Path()
				updateBookListChannel <- true
			}
		}, w)
	})
	return button
}

func updateBookList(bookListVBox *fyne.Container, updateNowPlayingChannel chan types.PlayingBook) {
	books, err := getAudioBooks()
	if err == nil {
		for _, v := range books {
			bookTileLayout := NewMyListItemWidget(v, updateNowPlayingChannel)
			bookListVBox.Add(bookTileLayout)
		}
	}
}

func getAudioBooks() ([]types.Book, error) {
	var bookList = []types.Book{}
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
								a := types.Book{
									Title:    b.Name(),
									FullPath: bookFullPath + "/" + i.Name(),
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
