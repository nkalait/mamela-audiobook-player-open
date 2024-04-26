package ui

import (
	"image/color"
	"mamela/buildConstraints"
	"mamela/filetypes"
	"mamela/merror"
	"mamela/storage"
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
	bass "github.com/pteich/gobass"
	"github.com/sqweek/dialog"
)

// Listens to events about changes to audio books root folder
var updateBookListChannel = make(chan bool)

// Initialise part of the UI that lists audio books and listen to update events
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
		books, err := getAudioBooks()
		if err != nil {
			merror.ShowError("An error has occurred", err)
		}
		if len(books) > 0 {
			return "Loaded Books"
		}
	}
	// TODO find a better we of doing the padding below
	return "Load Books     "
}

func generateBookListContainerTop() *fyne.Container {
	bookListHeaderTxt := canvas.NewText(setBookListHeader(), theme.ForegroundColor())
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
	hBox := container.NewStack(dots, bookListScroller)
	return hBox
}

func createFileDialogButton() *widget.Button {
	icon := theme.FolderOpenIcon()
	button := widget.NewButtonWithIcon("", icon, func() {
		path, err := dialog.Directory().Title("Open root folder").Browse()
		if err != nil {
			dialog.Message(err.Error())
		} else if path != "" {
			storage.Data.Root = path
			rootPath = storage.Data.Root
			storage.SaveDataToStorageFile()
			updateBookListChannel <- true
		}

	})
	return button
}

// Update the part of UI showing list of audio books
func updateBookList(bookListVBox *fyne.Container) {
	if rootPath != "" {
		bookListVBox.Objects = bookListVBox.Objects[:0]
		books, err := getAudioBooks()
		if err == nil {
			for _, v := range books {
				bookTileLayout := NewMyListItemWidget(v)
				bookListVBox.Add(bookTileLayout)
			}
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
		isAValidAudioBook := false
		if b.IsDir() {
			var bookFullPath = rootPath + buildConstraints.PathSeparator + b.Name()
			bookFolder, err := os.ReadDir(bookFullPath)
			if err == nil {
				highestQuality := int64(0)
				folderArt := ""
				var book types.Book
				for _, bookFile := range bookFolder {
					i, err := bookFile.Info()
					if err == nil {
						if i.Mode().IsRegular() {
							name := strings.ToLower(i.Name())
							if slices.Contains(filetypes.AllowedFileTypes, filepath.Ext(name)) {
								isAValidAudioBook = true
								chapter := types.Chapter{
									FileName:        i.Name(),
									LengthInSeconds: getChapterLengthInSeconds(bookFullPath, i.Name()),
								}
								book.Chapters = append(book.Chapters, chapter)
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
					book.FullLengthSeconds = getFullBookLengthSeconds(book.Chapters)
					book.Metadata = getFileTag(book)
					bookList = append(bookList, book)
				}
			}
		}
	}
	return bookList, nil
}

func getBookFile(b types.Book) *os.File {
	path := b.FullPath + buildConstraints.PathSeparator + b.Chapters[0].FileName
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

func getChapterLengthInSeconds(fullPath string, fileName string) float64 {
	length := float64(0)
	c, err := bass.StreamCreateFile(fullPath+buildConstraints.PathSeparator+fileName, 0, bass.AsyncFile)
	if err == nil {
		bytesLen, err := c.GetLength(bass.POS_BYTE)
		if err == nil {
			t, err := c.Bytes2Seconds(bytesLen)
			if err == nil {
				length = t
			}
		}
	}
	c.Free()
	return length
}

func getFullBookLengthSeconds(chapters []types.Chapter) float64 {
	length := float64(0)
	for i := 0; i < len(chapters); i++ {
		length = length + chapters[i].LengthInSeconds
	}
	return length
}
