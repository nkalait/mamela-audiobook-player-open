package ui

import (
	"mamela/audio"
	"mamela/types"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type MyListItemWidget struct {
	widget.BaseWidget
	Title  *widget.Label
	Button *widget.Button
}

// Update the image shown when loading an audio book
var funcChanFolderArtUpdaterCallBack = func(playingBook types.PlayingBook) {
	// Priority for showing an audio book image is such that we
	// first check if the file being played has an image
	// embedded it it, if not then we check if there is an
	// image file inside the audio book folder
	if playingBook.Metadata != nil && playingBook.Metadata.Picture() != nil {
		updateBookArt(playingBook.Metadata.Picture())
	} else if playingBook.FolderArt != "" {
		fileBytes, e := os.ReadFile(playingBook.FullPath + "/" + playingBook.FolderArt)
		if e == nil {
			pic := tag.Picture{Data: fileBytes}
			updateBookArt(&pic)
		} else {
			updateBookArt(nil)
		}
	} else {
		updateBookArt(nil)
	}
}

func NewMyListItemWidget(b types.Book) *MyListItemWidget {
	item := &MyListItemWidget{
		Title: widget.NewLabel(cases.Title(language.English).String(b.Title)),
		Button: widget.NewButton("", func() {
			var playingBook types.PlayingBook = types.PlayingBook{Book: b, CurrentChapter: 0, Position: 0, FullLengthSeconds: 0}
			audio.LoadAndPlay(playingBook, funcChanFolderArtUpdaterCallBack)
		}),
	}
	item.Title.Truncation = fyne.TextTruncateEllipsis
	item.ExtendBaseWidget(item)
	return item
}

func (item *MyListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	stack := container.NewStack(item.Button, item.Title)
	return widget.NewSimpleRenderer(stack)
}
