package ui

import (
	"mamela/audio"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MyListItemWidget struct {
	widget.BaseWidget
	// Title *canvas.Text
	Title *widget.Button
}

func NewMyListItemWidget(b book) *MyListItemWidget {
	item := &MyListItemWidget{
		Title: widget.NewButton(b.title, func() {
			// go func() {
			audio.LoadAndPlay(b.fullPath)
			// }()
		}),
	}
	item.Title.Alignment = widget.ButtonAlignLeading
	// item.Title = 18
	item.ExtendBaseWidget(item)
	return item
}

func (item *MyListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	v := container.New(layout.NewPaddedLayout(), layout.NewSpacer(), item.Title)
	// content := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), item.Title, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer())
	// a := container.NewBo(nil, nil, nil, item.Title)
	// c := container.NewStack(canvas.NewRectangle(BgColourLight), v)
	return widget.NewSimpleRenderer(v)
}
