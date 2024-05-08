package ui

import (
	"image/color"
	"mamela/audio"
	"mamela/buildconstraints"
	"mamela/filetypes"
	"mamela/merror"
	"mamela/storage"
	"mamela/types"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
	bass "github.com/pteich/gobass"
)

var bookListVBox *fyne.Container

// Initialise part of the UI that lists audio books and listen to update events
func initBookList() *fyne.Container {
	bookListVBox = container.New(layout.NewVBoxLayout())
	bookListContainer := initBookPane(bookListVBox)
	// Note, passing in true here will overwrite the book store in the storage file on disk,
	// this also means losing saved book last play positions
	updateBookList(false)

	// Listen to book list update events
	go func() {
		for update := range audio.UpdateBookListChannel {
			if update {
				updateBookList(true)
			}
		}
	}()

	return bookListContainer
}

func setBookListHeader() string {
	if len(storage.Data.BookList) > 0 {
		return "Loaded Books"
	}
	// TODO find a better we of doing the padding below
	return "Load Books     "
}

func generateBookListContainerTop() *fyne.Container {
	bookListHeaderTxt := canvas.NewText(setBookListHeader(), theme.ForegroundColor())
	bookListHeaderTxt.TextSize = 24
	bookListHeaderTxt.TextStyle.Bold = true
	spacer := canvas.NewText("    ", color.Transparent)
	top := container.NewHBox(bookListHeaderTxt, spacer, container.NewVBox(createRefreshButton()))
	return top
}

func initBookPane(bookListVBox *fyne.Container) *fyne.Container {
	// TODO the dots below are just to give the scroller the desired width, NEED TO FIND A WAY TO DO THIS BETTER!!
	dots := canvas.NewText("..........................................................", color.Transparent)
	bookListScroller := container.NewVScroll(bookListVBox)
	hBox := container.NewStack(dots, bookListScroller)
	return hBox
}

func createRefreshButton() *widget.Button {
	icon := theme.ViewRefreshIcon()
	button := widget.NewButtonWithIcon("", icon, func() {
		refreshBookList()
	})
	return button
}

func refreshBookList() {
	audio.UpdateBookListChannel <- true
}

// Update the part of UI showing list of audio books,
// if readRootFolder is true then we scan the root folder and save contents
// to the storage file
func updateBookList(readRootFolder bool) {
	if storage.Data.Root != "" {
		if readRootFolder {
			parseRootFolder()
		}
		bookListVBox.Objects = bookListVBox.Objects[:0]
		for _, v := range storage.Data.BookList {
			bookTileLayout := NewMyListItemWidget(v)
			bookListVBox.Add(bookTileLayout)
			loadPreviousBookOnLoad(v.Path, bookTileLayout.Button)
		}
		bookListVBox.Refresh()
	}
}

func loadPreviousBookOnLoad(bookFolderName string, bookLoadButton *widget.Button) {
	if storage.Data.CurrentBookFolder == bookFolderName {
		go func() {
			for {
				time.Sleep(500 * time.Millisecond)
				if MainWindow != nil {
					bookLoadButton.OnTapped()
					audio.Pause()
					break
				}
			}
		}()
	}
}

func parseRootFolder() {
	var bookList = []types.Book{}
	rootFolderEntries, err := os.ReadDir(storage.Data.Root)
	if err != nil {
		merror.ShowError("Could not read root folder", err)
		return
	}

	for _, folder := range rootFolderEntries {
		isAValidAudioBook := false
		if folder.IsDir() {
			bookFullPath := storage.Data.Root + buildconstraints.PathSeparator + folder.Name()
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
					book.Title = folder.Name()
					book.Path = folder.Name()
					book.FolderArt = folderArt
					book.FullLengthSeconds = getFullBookLengthSeconds(book.Chapters)
					book.Metadata = getFileTag(book)
					bookList = append(bookList, book)
				}
			}
		}
	}
	storage.SaveBookListToStorageFile(bookList)
}

func getBookFile(b types.Book) *os.File {
	path := storage.Data.Root + buildconstraints.PathSeparator + b.Path + buildconstraints.PathSeparator + b.Chapters[0].FileName
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
	c, err := bass.StreamCreateFile(fullPath+buildconstraints.PathSeparator+fileName, 0, bass.AsyncFile)
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
