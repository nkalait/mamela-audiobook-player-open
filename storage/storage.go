// Package storage provides a crude data storage for the app runtime

package storage

import (
	"encoding/json"
	"errors"
	"mamela/merror"
	"mamela/types"
	"os"
	"time"
)

var StorageFile = "data.json"

type Store struct {
	Root              string       `json:"root"`                // Folder where audio book folders can be found
	BookList          []types.Book `json:"books"`               // Audio books in the root folder
	CurrentBookFolder string       `json:"current_book_folder"` // Currently playing audio book
}

var Data Store = Store{}

// Load storage file data
func LoadStorageFile() {
	if checkStorageFile() {
		readJSONFile()
	}
}

// Check if storage file exists and is valid
func checkStorageFile() bool {
	fileExisted := false
	if _, err := os.Stat(StorageFile); err == nil {
		fileExisted = true
		go func() {
			time.Sleep(time.Second * 3)
			storageFileIsValid()
		}()
	} else if errors.Is(err, os.ErrNotExist) {
		f, err := os.OpenFile(StorageFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			go func() {
				time.Sleep(time.Second * 3)
				merror.ShowError("Error creating storage file", err)
			}()
		} else {
			_, err = f.WriteString("{}")
			f.Close()
			if err != nil {
				go func() {
					time.Sleep(time.Second * 3)
					merror.ShowError("Error writing to storage file", err)
				}()
			}
		}
	}

	return fileExisted
}

// Check if the storage file is a valid JSON file
func storageFileIsValid() bool {
	file, err := os.ReadFile(StorageFile)
	if err != nil {
		merror.ShowError("There is a problem with the storage file", err)
		return false
	}
	if !json.Valid(file) {
		merror.ShowError("The storage file does not seem to be valid", errors.New("invalid json"))
		return false
	}
	return true
}

// Read storage JSON file data into our Data variable
func readJSONFile() {
	file, err := os.ReadFile(StorageFile)
	merror.ShowError("Problem reading storage file", err)
	merror.PanicError(err)
	json.Unmarshal(file, &Data)
}

// Save data in Data variable to file on disk
func SaveDataToStorageFile() {
	jsonString, err := json.Marshal(Data)
	if err != nil {
		merror.ShowError("Internal error, could not marshal data", err)
		return
	}
	err = os.WriteFile(StorageFile, jsonString, os.ModePerm)
	merror.ShowError("Problem writing to storage file", err)
}

// Update list of books in storage file
func SaveBookListToStorageFile(bookList []types.Book) {
	Data.BookList = bookList
	SaveDataToStorageFile()
}

// Update the currently playing audio book data in storage file
func UpdateCurrentBook(bookPath string) {
	Data.CurrentBookFolder = bookPath
	SaveDataToStorageFile()
}
