package types

import "time"

// An audiobook
type Book struct {
	Title    string
	FullPath string
	Chapters []string // Each file name
}

// An audiobook that is currently playing
type PlayingBook struct {
	Book
	CurrentChapter    int
	Position          time.Duration
	FullLengthSeconds float64
}
