// Package storage provides a crude data storage for the app runtime

package storage

import (
	"encoding/json"
	"errors"
	"log"
	"mamela/err"
	"os"
	"time"
)

var storageFile = "data.json"

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
	if _, e := os.Stat(storageFile); e == nil {
		fileExisted = true
	} else if errors.Is(e, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		_, e := os.OpenFile(storageFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if e != nil {
			log.Panicln(e)
		}
	}
	go func() {
		time.Sleep(time.Second * 5)
		storageFileIsValid()
	}()
	return fileExisted
}

func storageFileIsValid() bool {
	file, e := os.ReadFile(storageFile)
	if e != nil {
		err.ShowError("There is a problem with the storage file", e)
		return false
	}
	if !json.Valid(file) {
		err.ShowError("The storage file does not seem to be valid", errors.New("invalid json"))
		return false
	}
	return true
}

func readJSONToken() {
	var d Store
	file, e := os.ReadFile(storageFile)
	err.ShowError("Problem reading storage file", e)
	err.PanicError(e)
	json.Unmarshal(file, &d)
	Data.Root = d.Root
	log.Printf("read storage file at %s\n", storageFile)
}

func SaveDataToStorageFile() {
	jsonString, _ := json.Marshal(Data)
	e := os.WriteFile(storageFile, jsonString, os.ModePerm)
	err.ShowError("Problem writing to storage file", e)
	err.PanicError(e)
	log.Printf("saved storage file: %x\n", jsonString)
}
