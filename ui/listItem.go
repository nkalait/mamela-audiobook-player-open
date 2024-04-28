package ui

import (
	"bytes"
	"image"
	"mamela/audio"
	"mamela/buildConstraints"
	"mamela/storage"
	"mamela/types"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type MyListItemWidget struct {
	widget.BaseWidget
	Icon   *canvas.Image
	Title  *widget.Label
	Button *widget.Button
}

func getBookImage(book types.Book) *tag.Picture {
	var pic *tag.Picture = nil
	if book.Metadata != nil && book.Metadata.Picture() != nil {
		pic = book.Metadata.Picture()
	} else if book.FolderArt != "" {
		fileBytes, err := os.ReadFile(storage.Data.Root + buildConstraints.PathSeparator + book.Path + buildConstraints.PathSeparator + book.FolderArt)
		if err == nil {
			pic = &tag.Picture{Data: fileBytes}
		}
	}

	return pic
}

// Update the image shown when loading an audio book
var funcChanFolderArtUpdaterCallBack = func(playingBook types.PlayingBook) {
	pic := getBookImage(playingBook.Book)
	updateBookArt(pic)
}

// TODO find a better way to do the button icon and text placement
func NewMyListItemWidget(b types.Book) *MyListItemWidget {
	title := cases.Title(language.English).String(b.Title)
	var button *widget.Button
	bookImage := getBookImage(b)
	var img image.Image

	if bookImage != nil {
		img, _, _ = image.Decode(bytes.NewReader(bookImage.Data))
	}

	var res = &fyne.StaticResource{
		StaticName:    "",
		StaticContent: []byte{},
	}

	callback := func() {
		// var playingBook types.PlayingBook = types.PlayingBook{Book: b, CurrentChapter: 0, Finished: false, Position: 0}
		var playingBook types.PlayingBook = types.PlayingBook{Book: b, CurrentChapter: 0, Finished: false}
		audio.LoadAndPlay(playingBook, true, funcChanFolderArtUpdaterCallBack)
		// Pause the audio book immediately after loading it, ...just to be user friendly
		audio.Pause()
	}

	if bookImage != nil {
		res.StaticContent = bookImage.Data
		button = widget.NewButtonWithIcon("", res, callback)
		title = "       " + cases.Title(language.English).String(b.Title)
	} else {
		button = widget.NewButton("", callback)
	}
	button.Alignment = widget.ButtonAlignLeading

	item := &MyListItemWidget{
		Icon:   canvas.NewImageFromImage(img),
		Title:  widget.NewLabel(title),
		Button: button,
	}
	item.Title.Truncation = fyne.TextTruncateEllipsis
	item.ExtendBaseWidget(item)

	return item
}

func (item *MyListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	stack := container.NewStack(item.Button, item.Title)
	return widget.NewSimpleRenderer(stack)
}
