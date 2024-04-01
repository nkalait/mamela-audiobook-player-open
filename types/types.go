package types

import "time"

// An audiobook
type Book struct {
	Title    string
	FullPath string
}

// An audiobook that is currently playing
type PlayingBook struct {
	Book
	Position          time.Duration
	FullLengthSeconds float64
}
