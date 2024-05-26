package helpers

import (
	"io/fs"
	"mamela/buildconstraints"
	"mamela/merror"
	"mamela/storage"
	"mamela/types"
	"os"
)

func ReadRootDirectory() []fs.DirEntry {
	rootFolderEntries, err := os.ReadDir(storage.Data.Root)
	if err != nil {
		merror.ShowError("Could not read root folder", err)
		return []fs.DirEntry{}
	}
	return rootFolderEntries
}

func GetBookList(rootFolderEntries []fs.DirEntry) []types.Book {
	bookList := []types.Book{}
	for _, folder := range rootFolderEntries {
		if folder.IsDir() {
			bookFullPath := storage.Data.Root + buildconstraints.PathSeparator + folder.Name()
			bookFolderContents, err := os.ReadDir(bookFullPath)

			if err == nil {
				book := getBookContents(folder.Name(), bookFullPath, bookFolderContents)
				if len(book.Chapters) > 0 {
					bookList = append(bookList, book)
				}
			}
		}
	}
	bookList = sortBooks(bookList)
	bookList = sortBookChaptersByNumerical(bookList)
	return bookList
}
