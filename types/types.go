package types

import "time"

type Book struct {
	Title    string
	FullPath string
}

type PlayingBook struct {
	Book
	Position time.Duration
}
