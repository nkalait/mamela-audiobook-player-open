package ui

import (
	"bytes"
	"image"
	"image/color"
	"mamela/audio"
	"mamela/buildconstraints"
	"mamela/bundled"
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

func getBookImage(book types.Book) []byte {
	var picBytes []byte
	var pic *tag.Picture = nil
	if book.Missing {
		return bundled.ResourceIconPng.StaticContent
	}
	if book.Metadata != nil && book.Metadata.Picture() != nil {
		pic = book.Metadata.Picture()
		picBytes = pic.Data
	} else if book.FolderArt != "" {
		fileBytes, err := os.ReadFile(storage.Data.Root + buildconstraints.PathSeparator + book.Path + buildconstraints.PathSeparator + book.FolderArt)
		if err == nil {
			pic = &tag.Picture{Data: fileBytes}
			picBytes = pic.Data
		}
	} else {
		picBytes = bundled.ResourceIconPng.StaticContent
	}

	return picBytes
}

// Update the image shown when loading an audio book
var funcChanFolderArtUpdaterCallBack = func(playingBook types.PlayingBook) {
	pic := getBookImage(playingBook.Book)
	updateBookArt(pic)
}

// TODO find a better way to do the button icon and text placement
func NewMyListItemWidget(b types.Book) *MyListItemWidget {
	item := &MyListItemWidget{}

	title := cases.Title(language.English).String(b.Title)
	if b.Missing {
		title += " (missing)"
	}
	var button *widget.Button
	bookImage := getBookImage(b)
	var img image.Image

	img, _, _ = image.Decode(bytes.NewReader(bookImage))

	// var res = &fyne.StaticResource{
	// 	StaticName:    "",
	// 	StaticContent: []byte{},
	// }

	onTapped := func() {}
	if !b.Missing {
		onTapped = func() {
			var playingBook types.PlayingBook = types.PlayingBook{Book: b, CurrentChapter: 0, Finished: false}
			audio.LoadAndPlay(playingBook, false, true, funcChanFolderArtUpdaterCallBack)
			audio.NotifyNewBookLoaded <- audio.GetCurrentBookFullLength()
		}
	}

	// res.StaticContent = bookImage
	button = widget.NewButton("", onTapped)
	if b.Missing {
		button.Disable()
	}

	item = &MyListItemWidget{
		Icon:   canvas.NewImageFromImage(img),
		Title:  widget.NewLabel(title),
		Button: button,
	}

	// item.Title.Truncation = fyne.TextTruncateEllipsis
	item.Title.Wrapping = fyne.TextWrapWord
	item.ExtendBaseWidget(item)

	return item
}

func (item *MyListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	c := fyne.NewContainer(item.Icon) //, item.Title)
	label := widget.NewLabel("\n\n")  // ensure minimum height for the AllItemsStack
	allItemsStack := container.NewStack(
		item.Button,
		c,
		container.NewBorder(nil, nil, canvas.NewText("              ", color.Opaque), nil, item.Title), // TODO: find a better way of doing this
		label,
	)

	item.Icon.Resize(fyne.Size{Width: 50, Height: 50})

	item.Title.Move(fyne.Position{X: 60, Y: 0})
	item.Icon.Move(fyne.Position{X: 6, Y: 11})

	return widget.NewSimpleRenderer(allItemsStack)
}
