package helpers

import (
	"errors"
	"io/fs"
	"log"
	"mamela/buildconstraints"
	"mamela/filetypes"
	"mamela/types"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"

	bass "github.com/pteich/gobass"
)

func getFullBookLengthSeconds(chapters []types.Chapter) float64 {
	length := float64(0)
	for i := 0; i < len(chapters); i++ {
		length = length + chapters[i].LengthInSeconds
	}
	return length
}

func getChapterLengthInSeconds(fullPath string, fileName string) float64 {
	length := float64(0)
	c, err := bass.StreamCreateFile(fullPath+buildconstraints.PathSeparator+fileName, 0, bass.DeviceMono)

	if err != nil {
		log.Println("Error loading media file to get length: ", err.Error())
	} else {
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

func getBookContents(bookFolderName string, bookFullPath string, bookFolderContents []fs.DirEntry) types.Book {
	highestQuality := int64(0)
	folderArt := ""
	var book types.Book
	for _, bookFile := range bookFolderContents {
		var fileInfo fs.FileInfo
		fileInfo, err := bookFile.Info()
		if err == nil {
			if fileInfo.Mode().IsRegular() {
				c, err := getBookChapter(fileInfo, bookFullPath)
				if err == nil {
					book.Chapters = append(book.Chapters, c)
				}
				tmpArt, tmpQuality := getFolderArtFile(fileInfo, highestQuality)
				if tmpArt != "" {
					folderArt = tmpArt
					highestQuality = tmpQuality
				}
			}
		}
	}
	book.Title = bookFolderName
	book.Path = bookFolderName
	book.FolderArt = folderArt
	book.FullLengthSeconds = getFullBookLengthSeconds(book.Chapters)

	return book
}

func getBookChapter(fileInfo fs.FileInfo, bookFullPath string) (types.Chapter, error) {
	var err error
	var bookChapter types.Chapter
	name := strings.ToLower(fileInfo.Name())
	if slices.Contains(filetypes.AllowedFileTypes, filepath.Ext(name)) {
		bookChapter = types.Chapter{
			FileName:        fileInfo.Name(),
			LengthInSeconds: getChapterLengthInSeconds(bookFullPath, fileInfo.Name()),
		}
	} else {
		err = errors.New("file is not and allowed audio book media file type: " + fileInfo.Name())
	}
	return bookChapter, err
}

func getFolderArtFile(bookFolderFile fs.FileInfo, foldArtSize int64) (string, int64) {
	folderArt := ""
	if slices.Contains(filetypes.BookArtFileTypes, filepath.Ext(strings.ToLower(bookFolderFile.Name()))) {
		if bookFolderFile.Size() > foldArtSize {
			foldArtSize = bookFolderFile.Size()
			folderArt = bookFolderFile.Name()
		}
	}
	return folderArt, foldArtSize
}

// Reading a folder's files seems to place 11 before 2, eg, 0,1,11,2,3,...
func sortBookChaptersByNumerical(bookList []types.Book) []types.Book {
	for _, book := range bookList {
		sort.Slice(book.Chapters, func(i, j int) bool {
			tempi := ""
			for _, letteri := range strings.Split(book.Chapters[i].FileName, "") {
				_, err := strconv.ParseInt(letteri, 10, 64)
				if err == nil {
					tempi += letteri
				} else {
					break
				}
			}
			tempj := ""
			for _, letterj := range strings.Split(book.Chapters[j].FileName, "") {
				_, err := strconv.ParseInt(letterj, 10, 64)
				if err == nil {
					tempj += letterj
				} else {
					break
				}
			}

			numi, err := strconv.ParseInt(tempi, 10, 64)
			if err == nil {
				numj, err := strconv.ParseInt(tempj, 10, 64)
				if err == nil {
					return numi < numj
				}
			}

			return false
		})
	}
	return bookList
}

// Reading a directory seems to place lowercase names after uppercase names, eg, Aa, Bb, aa, and not aa, Aa, Bb
func sortBooks(bookList []types.Book) []types.Book {
	sort.Slice(bookList, func(i, j int) bool {
		return strings.ToLower(bookList[i].Path) < strings.ToLower(bookList[j].Path)
	})
	return bookList
}
