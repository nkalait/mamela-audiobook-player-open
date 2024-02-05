package ui

import (
	"mamela/audio"

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

func NewMyListItemWidget(b book) *MyListItemWidget {
	item := &MyListItemWidget{
		Title: widget.NewLabel(cases.Title(language.English).String(b.title)),
		Button: widget.NewButton("", func() {
			audio.LoadAndPlay(b.fullPath)
		}),
	}
	item.Title.Truncation = fyne.TextTruncateEllipsis
	// item.Title.Alignment = widget.ButtonAlignLeading
	item.ExtendBaseWidget(item)
	return item
}

func (item *MyListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	stack := container.NewStack(item.Button, item.Title)
	return widget.NewSimpleRenderer(stack)
}
