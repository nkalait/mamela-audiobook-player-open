// Package storage provides a crude data storage for the app runtime

package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mamela/merror"
	"mamela/types"
	"os"
	"sort"
	"strings"
	"time"
)

var StorageFile = "data.json"

const jsonStrVolumeLevel = "volume_level"
const defaultVolumeLevel = 5000 // this is bass's 50%

type Store struct {
	Root              string       `json:"root"`                // Folder where audio book folders can be found
	BookList          []types.Book `json:"books"`               // Audio books in the root folder
	CurrentBookFolder string       `json:"current_book_folder"` // Currently playing audio book
	VolumeLevel       float64      `json:"volume_level"`        // The last set volume level
}

var Data Store = Store{}

// Load storage file data
func LoadStorageFile() {
	if checkStorageFile() {
		if !readJSONFile() {
			ClearAll()
		}
	}
}

func ClearBooks() {
	for i := range Data.BookList {
		Data.BookList[i].Missing = true
	}
	Data.CurrentBookFolder = ""
	SaveDataToStorageFile()
}

func ClearAll() bool {
	Data.BookList = []types.Book{}
	Data.Root = ""
	Data.CurrentBookFolder = ""
	saved := SaveDataToStorageFile()
	return saved
}

// Check if storage file exists and is valid
func checkStorageFile() bool {
	fileExists := false
	if _, err := os.Stat(StorageFile); err == nil {
		fileExists = true
		wait := make(chan bool)
		go func() {
			time.Sleep(time.Second * 3)
			if !storageFileIsValid() {
				ClearAll()
			}
			wait <- true
		}()
		<-wait

	} else if errors.Is(err, os.ErrNotExist) {
		f, err := os.OpenFile(StorageFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			go func() {
				time.Sleep(time.Second * 10)
				merror.ShowError("Error creating storage file", err)
			}()
		} else {
			_, err = f.WriteString("{\"" + jsonStrVolumeLevel + "\":" + fmt.Sprint(defaultVolumeLevel) + "}")
			if err != nil {
				go func() {
					time.Sleep(time.Second * 10)
					merror.ShowError("Error writing to storage file", err)
				}()
			}
			Data.VolumeLevel = defaultVolumeLevel
			f.Close()
		}
	}

	return fileExists
}

// Check if the storage file is a valid JSON file
func storageFileIsValid() bool {
	file, err := os.ReadFile(StorageFile)
	if err != nil {
		log.Println("There is a problem with the storage file, a new storage file has been created", err.Error())
		go func() {
			time.Sleep(time.Second * 10)
			merror.ShowError("There is a problem with the storage file, a new storage file has been created", errors.New("invalid json"))
		}()
		return false
	}
	if !json.Valid(file) {
		log.Println("The storage file does not seem to be valid, a new storage file has been created", errors.New("invalid json"))
		go func() {
			time.Sleep(time.Second * 10)
			merror.ShowError("The storage file does not seem to be valid, a new storage file has been created", errors.New("invalid json"))
		}()
		return false
	}
	return true
}

// Read storage JSON file data into our Data variable
func readJSONFile() bool {
	file, err := os.ReadFile(StorageFile)
	if err != nil {
		log.Println("Problem reading storage file: ", err.Error())
		go func() {
			time.Sleep(time.Second * 10)
			merror.ShowError("Problem reading storage file", err)
		}()
		return false
	} else {
		err = json.Unmarshal(file, &Data)
		if err != nil {
			log.Println("Problem unpacking storage file: ", err.Error())
			go func() {
				time.Sleep(time.Second * 10)
				merror.ShowError("Problem unpacking storage file", err)
			}()
			return false
		}
		return true
	}
}

// Save data in Data variable to file on disk
func SaveDataToStorageFile() bool {
	jsonString, err := json.Marshal(Data)
	if err != nil {
		go func() {
			time.Sleep(time.Second * 10)
			merror.ShowError("Internal error, could not marshal data", err)
		}()
		return false
	}
	err = os.WriteFile(StorageFile, jsonString, os.ModePerm)
	if err != nil {
		go func() {
			time.Sleep(time.Second * 10)
			merror.ShowError("Error writing to storage file", err)
		}()
	}
	return true
}

// Update list of books in storage file
func SaveBookListToStorageFile(bookList []types.Book) {
	existing := make(map[string]types.Book, len(Data.BookList))
	for _, book := range Data.BookList {
		existing[strings.ToLower(book.Path)] = book
	}

	merged := make([]types.Book, 0, len(bookList)+len(Data.BookList))
	added := make(map[string]struct{}, len(bookList))

	for _, incoming := range bookList {
		key := strings.ToLower(incoming.Path)
		if current, ok := existing[key]; ok {
			merged = append(merged, mergeBookData(current, incoming))
		} else {
			incoming.Missing = false
			merged = append(merged, incoming)
		}
		added[key] = struct{}{}
	}

	for _, book := range Data.BookList {
		key := strings.ToLower(book.Path)
		if _, ok := added[key]; ok {
			continue
		}
		book.Missing = true
		merged = append(merged, book)
	}

	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].Missing != merged[j].Missing {
			return !merged[i].Missing && merged[j].Missing
		}
		return strings.ToLower(merged[i].Path) < strings.ToLower(merged[j].Path)
	})

	Data.BookList = merged
	SaveDataToStorageFile()
}

// DeleteBook removes a book entry from storage by its path and persists the change.
func DeleteBook(path string) bool {
	if path == "" {
		return false
	}
	lowerPath := strings.ToLower(path)
	filtered := make([]types.Book, 0, len(Data.BookList))
	removed := false
	for _, book := range Data.BookList {
		if strings.ToLower(book.Path) == lowerPath {
			removed = true
			continue
		}
		filtered = append(filtered, book)
	}
	if !removed {
		return false
	}
	Data.BookList = filtered
	if strings.EqualFold(Data.CurrentBookFolder, path) {
		Data.CurrentBookFolder = ""
	}
	SaveDataToStorageFile()
	return true
}

func mergeBookData(existing, incoming types.Book) types.Book {
	merged := existing
	merged.Title = incoming.Title
	merged.Path = incoming.Path
	merged.Chapters = incoming.Chapters
	merged.FolderArt = incoming.FolderArt
	merged.FullLengthSeconds = incoming.FullLengthSeconds
	merged.Missing = false
	return merged
}

// Update the currently playing audio book data in storage file
func UpdateCurrentBook(bookPath string) {
	Data.CurrentBookFolder = bookPath
	SaveDataToStorageFile()
}

func SaveVolumeLevel(vol float64) {
	Data.VolumeLevel = vol
	SaveDataToStorageFile()
}

func GetVolumeLevel() float64 {
	return Data.VolumeLevel
}
