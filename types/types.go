package types

type Book struct {
	Title    string
	FullPath string
}

type PlayingBook struct {
	Book
	Position uint16
}
