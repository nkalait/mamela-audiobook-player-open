package ui

import (
	"mamela/audio"
	"mamela/types"

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

var channelUpdateBookArt = make(chan *tag.Picture)

func NewMyListItemWidget(b types.Book) *MyListItemWidget {
	item := &MyListItemWidget{
		Title: widget.NewLabel(cases.Title(language.English).String(b.Title)),
		Button: widget.NewButton("", func() {
			var playingBook types.PlayingBook = types.PlayingBook{Book: b, CurrentChapter: 0, Position: 0, FullLengthSeconds: 0}
			audio.LoadAndPlay(playingBook, channelUpdateBookArt)
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
