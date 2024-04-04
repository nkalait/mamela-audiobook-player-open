package main

// import (
// 	_ "unsafe"

// 	/*
// 	   #cgo darwin LDFLAGS: -L${SRCDIR}/lib/mac -llibbass
// 	*/
// 	"C"
// )

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
const rootPath string = ""

// Carries info about currently playing audiobook
var updateNowPlayingChannel = make(chan types.PlayingBook)

// Listens to exit app event
var exitApp = make(chan bool)

func main() {
	go func() {
		<-exitApp
		ui.MainWindow.Close()

	}()
	go func() {
		audio.StartChannelListener(updateNowPlayingChannel, exitApp)
	}()
	ui.BuildUI(appLabel, rootPath, updateNowPlayingChannel)
}
