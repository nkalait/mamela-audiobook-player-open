package ui

import (
	"mamela/audio"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
		// Title: canvas.NewText(title, textColour),
		Title: widget.NewButton(b.title, func() {
			go func() {
				audio.Stop()
				// log.Println("play file " + b.fullPath)
				audio.LoadAndPlay(b.fullPath)
			}()
		}),
	}
	// item.Title.TextSize = 18
	item.ExtendBaseWidget(item)
	return item
}

// func (p *MyListItemWidget) BackgroundColor() color.Color {
// 	return color.RGBA{255, 20, 147, 255}
// }

func (item *MyListItemWidget) CreateRenderer() fyne.WidgetRenderer {
	v := container.New(layout.NewPaddedLayout(), layout.NewSpacer(), item.Title)
	// content := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), item.Title, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer())
	// a := container.NewBo(nil, nil, nil, item.Title)
	c := container.NewStack(canvas.NewRectangle(BgColourLight), v)
	return widget.NewSimpleRenderer(c)
}
