package ui

import (
	"os"
	"path/filepath"
	"slices"

	"fyne.io/fyne/v2"
)

func updateBookList(bookListVBox *fyne.Container) {
	books, err := getAudioBooks()
	if err == nil {
		for _, v := range books {
			bookTileLayout := NewMyListItemWidget(v)
			bookListVBox.Add(bookTileLayout)
		}
	}
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
