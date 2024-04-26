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
	if _, e := os.Stat(StorageFile); e == nil {
		fileExisted = true
	} else if errors.Is(e, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		_, e := os.OpenFile(StorageFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if e != nil {
			log.Println(e)
		}
	}
	go func() {
		time.Sleep(time.Second * 3)
		storageFileIsValid()
	}()
	return fileExisted
}

func storageFileIsValid() bool {
	file, e := os.ReadFile(StorageFile)
	if e != nil {
		merror.ShowError("There is a problem with the storage file", e)
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
	file, e := os.ReadFile(StorageFile)
	merror.ShowError("Problem reading storage file", e)
	merror.PanicError(e)
	json.Unmarshal(file, &d)
	Data.Root = d.Root
}

func SaveDataToStorageFile() {
	jsonString, _ := json.Marshal(Data)
	e := os.WriteFile(StorageFile, jsonString, os.ModePerm)
	merror.ShowError("Problem writing to storage file", e)
	merror.PanicError(e)
}
