package ui

import (
	"image/color"
	"mamela/err"
	"mamela/filetypes"
	"mamela/types"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
	"github.com/sqweek/dialog"
)

// Listens to events about changes to audiobooks root folder
var updateBookListChannel = make(chan bool)

// Initialise part of the UI that lists audiobooks and listen to update events
func initBookList() *fyne.Container {
	bookListVBox := container.New(layout.NewVBoxLayout())
	bookListContainer := initBookPane(bookListVBox)
	updateBookList(bookListVBox)

	// Listen to book list update events
	go func() {
		for update := range updateBookListChannel {
			if update {
				updateBookList(bookListVBox)
			}
		}
	}()

	return bookListContainer
}

func setBookListHeader() string {
	if rootPath != "" {
		books, e := getAudioBooks()
		if e != nil {
			err.ShowError("An error has occured", e)
		}
		if len(books) > 0 {
			return "Loaded Books"
		}
	}
	// TODO find a better we of doing the padding below
	return "Load Books     "
}

func generateBookListContainerTop() *fyne.Container {
	bookListHeaderTxt := canvas.NewText(setBookListHeader(), textColour)
	bookListHeaderTxt.TextSize = 24
	bookListHeaderTxt.TextStyle.Bold = true
	spacer := canvas.NewText("    ", color.Transparent)
	top := container.NewHBox(bookListHeaderTxt, spacer, container.NewVBox(createFileDialogButton()))
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

func createFileDialogButton() *widget.Button {
	icon := theme.FolderOpenIcon()
	button := widget.NewButtonWithIcon("", icon, func() {
		path, e := dialog.Directory().Title("Open root folder").Browse()
		if e != nil {
			dialog.Message(e.Error())
		} else if path != "" {
			rootPath = path
			updateBookListChannel <- true
		}

	})
	return button
}

// Update the part of UI showing list of audiobooks
func updateBookList(bookListVBox *fyne.Container) {
	if rootPath != "" {
		bookListVBox.Objects = bookListVBox.Objects[:0]
		books, e := getAudioBooks()
		if e == nil {
			for _, v := range books {
				bookTileLayout := NewMyListItemWidget(v)
				bookListVBox.Add(bookTileLayout)
			}
		}
	}
}

func getAudioBooks() ([]types.Book, error) {
	var bookList = []types.Book{}
	rootFolderEntries, e := os.ReadDir(rootPath)
	if e != nil {
		return nil, e
	}

	for _, b := range rootFolderEntries {
		isAValidAudioBook := false
		if b.IsDir() {
			var bookFullPath = rootPath + "/" + b.Name()
			bookFolder, e := os.ReadDir(bookFullPath)
			if e == nil {
				highestQuality := int64(0)
				folderArt := ""
				var book types.Book
				for _, bookFile := range bookFolder {
					i, e := bookFile.Info()
					if e == nil {
						if i.Mode().IsRegular() {
							name := strings.ToLower(i.Name())
							if slices.Contains(filetypes.AllowedFileTypes, filepath.Ext(name)) {
								isAValidAudioBook = true
								book.Chapters = append(book.Chapters, i.Name())
							} else if slices.Contains(filetypes.BookArtFileTypes, filepath.Ext(name)) {
								// If the folder contains an image file, get the one of best quality
								if i.Size() > highestQuality {
									highestQuality = i.Size()
									folderArt = i.Name()
								}
							}
						}
					}
				}

				if isAValidAudioBook {
					book.Title = b.Name()
					book.FullPath = bookFullPath
					book.FolderArt = folderArt
					book.Metadata = getFileTag(book)
					bookList = append(bookList, book)
				}
			}
		}
	}
	return bookList, nil
}

func getBookFile(b types.Book) *os.File {
	path := b.FullPath + "/" + b.Chapters[0]
	f, _ := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	return f
}

func getFileTag(b types.Book) tag.Metadata {
	f := getBookFile(b)
	var meta tag.Metadata = nil
	if f != nil {
		meta, _ = tag.ReadFrom(f)
		f.Close()
	}
	return meta
}
