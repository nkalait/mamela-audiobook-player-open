package ui

import (
	"bytes"
	"fmt"
	"image"
	"mamela/audio"
	"mamela/types"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowden/tag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var bookArt *canvas.Image
var bookTitle *widget.Label
var bookFullLength *canvas.Text
var playingPosition *canvas.Text
var playerButtonPlay *widget.Button
var playerButtonPause *widget.Button
var playerButtonStop *widget.Button
var playerButtonFastRewind *widget.Button
var playerButtonFastForward *widget.Button
var playerButtonSkipNext *widget.Button
var playerButtonSkipPrevious *widget.Button

func createPlayingLayout(updateNowPlayingChannel chan types.PlayingBook) *fyne.Container {
	initUI()
	hideUIItems()

	containerTop := container.NewVBox(bookTitle, bookFullLength)
	containerPositionDetails := container.NewVBox(
		playingPosition,
		layoutPlayerButtons(),
	)
	playingVBox := container.NewPadded(container.NewBorder(
		containerTop,
		containerPositionDetails,
		nil,
		nil,
		bookArt,
	))

	go func() {
		for playingBook := range updateNowPlayingChannel {
			if bookTitle.Hidden {
				showUIItems()
			}
			updatePlaying(playingBook)
		}
	}()

	go func() {
		for playingState := range audio.ChannelAudioState {
			switch playingState {
			case audio.Stopped:
				playerButtonStop.Importance = widget.HighImportance
				playerButtonPlay.Importance = widget.MediumImportance
				playerButtonPause.Importance = widget.MediumImportance
			case audio.Paused:
				playerButtonStop.Importance = widget.MediumImportance
				playerButtonPlay.Importance = widget.MediumImportance
				playerButtonPause.Importance = widget.HighImportance
			case audio.Playing:
				playerButtonStop.Importance = widget.MediumImportance
				playerButtonPlay.Importance = widget.HighImportance
				playerButtonPause.Importance = widget.MediumImportance
			}
			playerButtonStop.Refresh()
			playerButtonPlay.Refresh()
			playerButtonPause.Refresh()
		}
	}()
	return playingVBox
}

func initUI() {
	initTitle()
	initBookArt()
	initPlayingPosition()
	initPlayerButtons()
	initFullBookLength()
}

func showUIItems() {
	bookArt.Show()
	bookTitle.Show()
	bookFullLength.Show()
	playingPosition.Show()

	playerButtonSkipPrevious.Show()
	playerButtonFastRewind.Show()
	playerButtonPause.Show()
	playerButtonStop.Show()
	playerButtonPlay.Show()
	playerButtonFastForward.Show()
	playerButtonSkipNext.Show()
}

func hideUIItems() {
	bookArt.Hide()
	bookTitle.Hide()
	bookFullLength.Hide()
	playingPosition.Hide()

	playerButtonSkipPrevious.Hide()
	playerButtonFastRewind.Hide()
	playerButtonPause.Hide()
	playerButtonStop.Hide()
	playerButtonPlay.Hide()
	playerButtonFastForward.Hide()
	playerButtonSkipNext.Hide()
}

func initBookArt() {
	bookArt = canvas.NewImageFromImage(nil)
	bookArt.FillMode = canvas.ImageFillContain
}

func initTitle() {
	// bookTitle = canvas.NewText("", textColour)
	bookTitle = widget.NewLabel("")
	// bookTitle.TextSize = 16
	bookTitle.Wrapping = fyne.TextWrapBreak
	bookTitle.TextStyle.Bold = true
	bookTitle.Alignment = fyne.TextAlignCenter
}

func initFullBookLength() {
	bookFullLength = canvas.NewText("", textColour)
	bookFullLength.TextSize = 18
	// bookFullLength.TextStyle.Bold = true
	bookFullLength.Alignment = fyne.TextAlignCenter
}

func initPlayingPosition() {
	playingPosition = canvas.NewText("", textColour)
	playingPosition.TextSize = 18
	playingPosition.Alignment = fyne.TextAlignCenter
}

func initPlayerButtons() {
	playerButtonPlay = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		audio.Play()
	})
	playerButtonPause = widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
		audio.Pause()
	})
	playerButtonStop = widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
		audio.Stop()
	})
	playerButtonFastRewind = widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		audio.FastRewind()
	})
	playerButtonFastForward = widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		audio.FastForward()
	})
	playerButtonSkipNext = widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
		audio.SkipNext()
	})
	playerButtonSkipPrevious = widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		audio.SkipPrevious()
	})
}

func layoutPlayerButtons() *fyne.Container {
	layout := container.NewHBox(
		playerButtonSkipPrevious,
		playerButtonFastRewind,
		playerButtonPause,
		playerButtonStop,
		playerButtonPlay,
		playerButtonFastForward,
		playerButtonSkipNext,
	)
	return container.NewCenter(layout)
}

func updatePlaying(p types.PlayingBook) {
	updateTitle(p.Title)
	updatePlayingPosition(p)
	d := time.Duration(math.Round(p.FullLengthSeconds * 1000000000))
	updateBookFullLength(audio.SecondsToTimeText(d))
}

func clearBookArt() {
	bookArt.Image = nil
	bookArt.Refresh()
}

func updateBookArt(pic *tag.Picture) {
	clearBookArt()
	if pic == nil {
		return
	}
	if len(pic.Data) == 0 {
		return
	}
	img, _, err := image.Decode(bytes.NewReader(pic.Data))
	if err != nil {
		fmt.Print(err)
		return
	}
	bookArt.Image = img
	bookArt.Refresh()
}

func updateTitle(title string) {
	bookTitle.Text = cases.Title(language.English).String(title)
	bookTitle.Refresh()
}

func updateBookFullLength(bookLength string) {
	bookFullLength.Text = bookLength
	bookFullLength.Refresh()
}

func updatePlayingPosition(p types.PlayingBook) {
	playingPosition.Text = audio.SecondsToTimeText(audio.GetCurrentBookPlayingDuration(p))
	playingPosition.Refresh()
}
