// Package storage provides a crude data storage for the app runtime

package storage

import (
	"encoding/json"
	"errors"
	"log"
	"mamela/merror"
	"os"
	"time"
)

var StorageFile = "data.json"

type Store struct {
	Root string `json:"root"` // Rest of the fields should go here.
}

var Data Store = Store{}

func LoadStorageFile() {
	if checkStorageFile() {
		readJSONToken()
	}
}

func checkStorageFile() bool {
	fileExisted := false
	if _, err := os.Stat(StorageFile); err == nil {
		fileExisted = true
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
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

func readJSONToken() {
	var d Store
	file, err := os.ReadFile(StorageFile)
	merror.ShowError("Problem reading storage file", err)
	merror.PanicError(err)
	json.Unmarshal(file, &d)
	Data.Root = d.Root
}

func SaveDataToStorageFile() {
	jsonString, _ := json.Marshal(Data)
	err := os.WriteFile(StorageFile, jsonString, os.ModePerm)
	merror.ShowError("Problem writing to storage file", err)
	merror.PanicError(err)
}
