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

func main() {
	// Load the storage file from disk, we are dong it here so that
	// if there is any error then we can show it on the UI
	storage.LoadStorageFile()

	go audio.StartChannelListener()

	// Wait for Bass to be initialised before moving on
	<-audio.BassInitiatedChan

	// There is a blocking call somewhere in there
	ui.BuildUI(appLabel)
}
