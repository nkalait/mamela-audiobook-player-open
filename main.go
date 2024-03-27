package main

import (
	"mamela/audio"
	"mamela/types"
	"mamela/ui"
)

// Basically, this app plays audio books.
// There is one root folder in which there are other folders.
// Each folder under the root folder respresents a single audio book.
// Each audio book folder should contain playable audio files, each file would be a chapter.

// The name of the app
const appLabel = "Mamela"

// The root folder where audio books will be place, this variable here is for testing only
const rootPath = "/some_path"

var updateNowPlayingChannel = make(chan types.PlayingBook)

func main() {
	go func() {
		audio.Initiate()
	}()
	ui.BuildUI(appLabel, rootPath, updateNowPlayingChannel)
}
