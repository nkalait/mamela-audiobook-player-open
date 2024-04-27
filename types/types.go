package types

import (
	"time"

	"github.com/dhowden/tag"
)

// Chapter data
type Chapter struct {
	LengthInSeconds float64
	FileName        string
}

// An audio book
type Book struct {
	Title             string
	Path              string
	Chapters          []Chapter // Each file name
	FolderArt         string
	FullLengthSeconds float64
	Metadata          tag.Metadata
	Position          time.Duration // Position in the currently playing file
}

// An audio book that is currently playing
type PlayingBook struct {
	Book
	CurrentChapter int  // Currently playing chapter
	Finished       bool // Has the audio book finished playing or not
}
