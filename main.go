package main

import (
	"mamela/audio"
	"mamela/storage"
	"mamela/ui"
)

// Basically, this app plays audio books.
// There is one root folder in which there are other folders.
// Each folder under the root folder represents a single audio book.
// Each audio book folder should contain playable audio files, each file would be a chapter.

// The name of the app
const appLabel = "Mamela"

// Listens to exit app event
// var exitApp = make(chan bool)

// // Handle quit, kill, cancel type of signals
// func handleSignals() {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)
// 	<-c
// 	onExit()
// 	// exitApp <- true
// }

// Notify the audio package that it is type to exit
// func onExit() {
// 	audio.ExitListener <- true
// }

func main() {
	// Load the storage file from disk, we are dong it here so that
	// if there is any error then we can show it on the UI
	storage.LoadStorageFile()

	// Listen for the exit app event
	// go func() {
	// 	<-exitApp
	// 	onExit()
	// 	ui.MainWindow.Close()
	// }()

	go audio.StartChannelListener()

	// go handleSignals()

	// Wait for Bass to be initialised before moving on
	<-audio.BassInitiatedChan

	// There is a blocking call somewhere in there
	ui.BuildUI(appLabel)

	// audio.ExitListener <- true
	// <-exitApp
}
