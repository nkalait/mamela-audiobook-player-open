package ui

import (
	"image/color"
	"log"

	"mamela/audio"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func BuildUI(appLabel string) {
	mamelaApp := app.New()
	window := mamelaApp.NewWindow(appLabel)

	audioFile1 := "/Users/nada/Dev/mamela/gunshot.mp3"
	audioFile2 := "/Users/nada/Dev/mamela/song.mp3"

	button1 := widget.NewButton("click me", func() {
		// defer streamer.Close()
		go func() {
			audio.Stop()
			log.Println("play file 1")
			audio.LoadAndPlay(audioFile1)
		}()
	})

	button2 := widget.NewButton("click me", func() {
		// defer streamer.Close()
		go func() {
			audio.Stop()
			log.Println("play file 2")
			audio.LoadAndPlay(audioFile2)
		}()
	})

	text2 := canvas.NewText("2", color.White)
	text3 := canvas.NewText("3", color.White)
	grid := container.New(layout.NewGridLayout(2), button1, button2, text2, text3)

	window.SetContent(grid)
	window.ShowAndRun()
}
