package types

import (
	"time"

	"github.com/dhowden/tag"
)

// for file by file based chapters, lentgh
type Chapter struct {
	LengthInSeconds float64
	FileName        string
}

// An audiobook
type Book struct {
	Title             string
	FullPath          string
	Chapters          []Chapter // Each file name
	FolderArt         string
	FullLengthSeconds float64
	Metadata          tag.Metadata
}

// An audiobook that is currently playing
type PlayingBook struct {
	Book
	CurrentChapter int
	Finished       bool
	Position       time.Duration
}
