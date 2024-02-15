package ui

import (
	"mamela/audio"
	"mamela/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type MyListItemWidget struct {
	widget.BaseWidget
	Title  *widget.Label
	Button *widget.Button
}

func NewMyListItemWidget(b types.Book) *MyListItemWidget {
	item := &MyListItemWidget{
		Title: widget.NewLabel(cases.Title(language.English).String(b.Title)),
		Button: widget.NewButton("", func() {
			var playingBook types.PlayingBook = types.PlayingBook{b, 0}
			updateNowPlayingChannel <- playingBook
			audio.LoadAndPlay(playingBook)
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
