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

// // D returns the duration of n samples.
// func (sr PlayingBook) D(n int) time.Duration {
// 	return time.Second * time.Duration(n) / time.Duration(sr)
// }

// // N returns the number of samples that last for d duration.
// func (sr PlayingBook) N(d time.Duration) int {
// 	return int(d * time.Duration(sr) / time.Second)
// }
