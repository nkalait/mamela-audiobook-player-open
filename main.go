package main

import (
	"mamela/audio"
	"mamela/storage"
	"mamela/ui"
	"os"
	"os/signal"
	"syscall"
)

// Basically, this app plays audio books.
// There is one root folder in which there are other folders.
// Each folder under the root folder respresents a single audio book.
// Each audio book folder should contain playable audio files, each file would be a chapter.

// The name of the app
const appLabel = "Mamela"

// The root folder where audio books will be place, this variable here is for testing only
var rootPath string = ""

// Carries info about currently playing audiobook

// Listens to exit app event
var exitApp = make(chan bool)

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)
	<-c
	onExit()
	// exitApp <- true
}

func onExit() {
	audio.ExitListener <- true
}

func main() {
	// defer onExit()
	storage.LoadStorageFile()
	rootPath = storage.Data.Root
	go func() {
		<-exitApp
		onExit()
		ui.MainWindow.Close()
	}()
	go func() {
		audio.StartChannelListener(exitApp)
	}()
	go handleSignals()
	<-audio.BassInitiatedChan
	ui.BuildUI(appLabel, rootPath)
	audio.ExitListener <- true
	<-exitApp
}
