package main

import "mamela/ui"

// Basically, this app plays audio books.
// There is one root folder in which there are other folders.
// Each folder under the root folder respresents a single audio book.
// Each audio book folder should contain playable audio files, each file would be a chapter.

// The name of the app
const appLabel = "Mamela"

// The root folder where audio books will be place, this variable here is for testing only
const rootPath = "/Users/nada/Dev/mamela/audio books"

func main() {
	ui.BuildUI(appLabel, rootPath)
}
