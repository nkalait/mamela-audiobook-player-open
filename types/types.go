package types

import (
	"time"

	"github.com/dhowden/tag"
)

// An audiobook
type Book struct {
	Title     string
	FullPath  string
	Chapters  []string // Each file name
	FolderArt string
	Metadata  tag.Metadata
}

// An audiobook that is currently playing
type PlayingBook struct {
	Book
	CurrentChapter    int
	Position          time.Duration
	FullLengthSeconds float64
}
