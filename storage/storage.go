// Package storage provides a crude data storage for the app runtime

package storage

import (
	"encoding/json"
	"errors"
	"log"
	"mamela/merror"
	"mamela/types"
	"os"
	"time"
)

var StorageFile = "data.json"

type Store struct {
	Root              string       `json:"root"` // Rest of the fields should go here.
	BookList          []types.Book `json:"books"`
	CurrentBookFolder string       `json:"current_book_folder"`
}

var Data Store = Store{}

func LoadStorageFile() {
	if checkStorageFile() {
		readJSONFile()
	}
}

func checkStorageFile() bool {
	fileExisted := false
	if _, err := os.Stat(StorageFile); err == nil {
		fileExisted = true
	} else if errors.Is(err, os.ErrNotExist) {
		_, err := os.OpenFile(StorageFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
	}
	go func() {
		time.Sleep(time.Second * 3)
		storageFileIsValid()
	}()
	return fileExisted
}

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

func readJSONFile() {
	file, err := os.ReadFile(StorageFile)
	merror.ShowError("Problem reading storage file", err)
	merror.PanicError(err)
	json.Unmarshal(file, &Data)
}

func SaveDataToStorageFile() {
	jsonString, err := json.Marshal(Data)
	if err != nil {
		merror.ShowError("Internal error, could not marshal data", err)
		return
	}
	err = os.WriteFile(StorageFile, jsonString, os.ModePerm)
	merror.ShowError("Problem writing to storage file", err)
}

func SaveBookListToStorageFile(bookList []types.Book) {
	Data.BookList = bookList
	SaveDataToStorageFile()
}

func UpdateCurrentBook(bookPath string) {
	Data.CurrentBookFolder = bookPath
	SaveDataToStorageFile()
}
