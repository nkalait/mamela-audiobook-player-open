package audio

import (
	"fmt"
	"log"
	"mamela/buildconstraints"
	"mamela/merror"
	"mamela/storage"
	"mamela/types"
	"slices"
	"time"

	bass "github.com/pteich/gobass"
)

const (
	PAUSED = iota
	PLAYING
	STOPPED
)

// Event listeners
var (
	ExitListener                    = make(chan bool) // for stopping to listen to playing events
	exitAudio                       = make(chan bool) // for unloading audio stuff
	BassInitiatedChan               = make(chan bool)
	NotifyInitReady                 = make(chan bool)
	UpdateNowPlayingChannel         = make(chan types.PlayingBook)
	UpdateBookListChannel           = make(chan bool)
	NotifyVolumeSliderDragged       = make(chan float64)
	NotifyBookPlayTime              = make(chan time.Duration)
	NotifyBookPlayTimeSliderDragged = make(chan float64)
	NotifyNewBookLoaded             = make(chan float64)
)

// Holds data structures important to playing an audio book
var player Player

// half a second delay before updating UI
const PlayingBookTickerDuration = 500 * time.Millisecond

var UIUpdateTicker *time.Ticker = time.NewTicker(PlayingBookTickerDuration)

// 15 second delay before saving currently playing books
// play position to disk
const CurrentBookPositionTickerDuration = time.Second * 15

var CurrentBookPositionUpdateTicker *time.Ticker = time.NewTicker(CurrentBookPositionTickerDuration)

// Initiate Bass
func init() {
	UIUpdateTicker.Stop()
	CurrentBookPositionUpdateTicker.Stop()
	go func() {
		<-NotifyInitReady
		plugins := loadPlugins()
		initBass()
		setVolumeSliderDragListener()
		setPlayTimeScrubberDragListener()
		<-exitAudio
		tearDown(plugins)
	}()

	go func() {
		for range CurrentBookPositionUpdateTicker.C {
			saveCurrentPlayingBookPositionToDisk()
		}
	}()

}

func setVolumeSliderDragListener() {
	go func() {
		for vol := range NotifyVolumeSliderDragged {
			player.setVolume(vol)
		}
	}()
}

func setPlayTimeScrubberDragListener() {
	go func() {
		for pos := range NotifyBookPlayTimeSliderDragged {
			if player.channel != 0 && pos != -1 {
				originalState := player.State
				if originalState == PLAYING {
					player.pause()
				}
				posDuration := time.Duration(pos * float64(time.Second))
				targetBytePos := 0
				if pos > 0 {
					targetBytePos = getTargetPositionInBytes(posDuration)
				} else {
					resetToBeginning()
				}
				err := player.channel.SetPosition(targetBytePos, bass.POS_BYTE)
				if err != nil {
					merror.ShowError("Could not set new position", err)
				}
				UpdateNowPlayingChannel <- player.currentBook
				if originalState == PLAYING {
					player.play()
				}
				saveDraggedPositionToDisk(posDuration)
			}
		}
	}()
}

func GetState() int {
	return player.State
}

func GetCurrentBookFullLength() float64 {
	return player.currentBook.FullLengthSeconds
}

func saveCurrentPlayingBookPositionToDisk() {
	if len(storage.Data.BookList) > 0 {
		idx := slices.IndexFunc(storage.Data.BookList, func(b types.Book) bool {
			return b.Path == player.currentBook.Path
		})
		if idx > -1 {
			storage.Data.BookList[idx].Position = GetCurrentBookPlayTime()
			storage.SaveBookListToStorageFile(storage.Data.BookList)
		}
	}
}

func saveDraggedPositionToDisk(pos time.Duration) {
	if len(storage.Data.BookList) == 0 {
		return
	}
	idx := slices.IndexFunc(storage.Data.BookList, func(b types.Book) bool {
		return b.Path == player.currentBook.Path
	})
	if idx == -1 {
		return
	}
	storage.Data.BookList[idx].Position = pos
	storage.SaveDataToStorageFile()
}

// Unload loaded Bass plugins and free all resources used by Bass
func tearDown(plugins []uint32) {
	for _, p := range plugins {
		bass.PluginFree(p)
	}
	bass.Free()
}

// Initialise Bass
func initBass() {
	err := bass.Init(-1, 44100, bass.DeviceStereo, 0, 0)
	merror.ShowError("Problem initiating bass", err)
	bass.SetVolume(100)
	BassInitiatedChan <- true
}

// Start listening to audio playing event and exit event
func StartChannelListener() {
	go func() {
	RoutineLoop:
		for {
			select {
			case <-UIUpdateTicker.C:
				updateUICurrentlyPlayingInfo()
			case <-ExitListener:
				break RoutineLoop
			}
		}
		exitAudio <- true
	}()
}

// Pad number below 10 with a zero
func Pad(i int) string {
	s := fmt.Sprint(i)
	if i < 10 {
		s = "0" + fmt.Sprint(i)
	}
	return s
}

// Convert seconds to time in hh : mm : ss
func SecondsToTimeText(seconds time.Duration) string {
	var h int = int(seconds.Seconds()) / 3600
	var m int = (int(seconds.Seconds()) / 60) % 60
	var s int = int(seconds.Seconds()) % 60

	sh := Pad(h)
	sm := Pad(m)
	ss := Pad(s)

	return sh + " : " + sm + " : " + ss
}

func GetCurrentBookPlayTime() time.Duration {
	pos := player.currentBook.Position
	if player.currentBook.CurrentChapter > 0 {
		for i := player.currentBook.CurrentChapter - 1; i >= 0; i-- {
			pos = pos + time.Duration(player.currentBook.Chapters[i].LengthInSeconds*1000000000)
		}
	}
	return pos
}

func updateUIOnStop() {
	UpdateNowPlayingChannel <- player.currentBook
	NotifyBookPlayTime <- 0
}

func ClearCurrentlyPlaying() {
	player.stop()
	CurrentBookPositionUpdateTicker.Stop()
	UIUpdateTicker.Stop()
	player.channel.Free()
	player.currentBook = types.PlayingBook{}
	UpdateNowPlayingChannel <- player.currentBook
}

// Update the currently playing audio book information on the UI
func updateUICurrentlyPlayingInfo() {
	if player.channel != 0 {
		active, err := player.channel.IsActive()
		merror.ShowError("Cannot check if channel is active", err)

		// We need active == bass.ACTIVE_STOPPED here in order to detect when
		// file has reached end
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_STOPPED {
			bytePosition, err := player.channel.GetPosition(bass.POS_BYTE)
			merror.ShowError("Cannot get byte position", err)
			pos, err := player.channel.Bytes2Seconds(bytePosition)
			merror.ShowError("Cannot convert bytes to seconds", err)

			// If audio book has multiple files; if a file in the book has reached the end then load the next file
			// and continue playing
			currentlyAt := player.currentBook.Position.Round(time.Millisecond * 500)
			skipAt := time.Duration(player.currentBook.Chapters[player.currentBook.CurrentChapter].LengthInSeconds * 1000000000).Round(time.Millisecond * 500)
			if currentlyAt == skipAt {
				skipToNextFile(&player, true, true, false)
			}

			// If we have reached the end of the audio book then stop playing
			posInWholeBook := GetCurrentBookPlayTime().Round(time.Millisecond * 500)
			wholeBookLength := time.Duration(player.currentBook.FullLengthSeconds * 1000000000).Round(time.Millisecond * 500)
			if posInWholeBook == wholeBookLength {
				player.currentBook.Finished = true
				UIUpdateTicker.Stop()
				CurrentBookPositionUpdateTicker.Stop()
				// ChannelAudioState <- Stopped
			}

			var d time.Duration = time.Duration(pos * 1000000000)
			player.currentBook.Position = time.Duration(d)
			NotifyBookPlayTime <- GetCurrentBookPlayTime()
		}
		UpdateNowPlayingChannel <- player.currentBook
	}
}

type UpdateFolderArtCallBack func(playingBook types.PlayingBook)

func LoadAndPlay(playingBook types.PlayingBook, continuePlaying bool, setPreviousPosition bool, updaterFolderArtCallback UpdateFolderArtCallBack) {
	stopPlayingIfPlaying()
	player.currentBook = playingBook

	chapter := player.currentBook.CurrentChapter
	loaded := loadAudioBookFile(storage.Data.Root + buildconstraints.PathSeparator + player.currentBook.Path + buildconstraints.PathSeparator + player.currentBook.Chapters[chapter].FileName)
	if loaded && setPreviousPosition {
		goToPreviousPosition()
	}

	if loaded {
		if continuePlaying {
			player.play()
		}

		if updaterFolderArtCallback != nil {
			updaterFolderArtCallback(player.currentBook)
		}
	}

	updateUICurrentlyPlayingInfo()
}

func stopPlayingIfPlaying() {
	if player.channel != 0 {
		a, err := player.channel.IsActive()
		if err == nil {
			if a == bass.ACTIVE_PLAYING || a == bass.ACTIVE_PAUSED {
				player.stop()
			}
		}
	}
}

func loadAudioBookFile(fullPath string) bool {
	var err error = nil
	// c, err := bass.StreamCreateURL("http://music.myradio.ua:8000/PopRock_news128.mp3", bass.DeviceStereo)
	player.channel, err = bass.StreamCreateFile(fullPath, 0, bass.AsyncFile)
	if err != nil {
		log.Println("Error loading media file: "+fullPath+"\n\n", err.Error())
		merror.ShowError("There seems to be a problem loading the the audio book file(s): "+fullPath, err)
		return false
	}
	return true
}

// We want to use the position saved to disk here so that we can resume playback
func goToPreviousPosition() error {
	savedPos := time.Duration(0)
	// Get the last play position that was saved to disk
	{
		idx := slices.IndexFunc(storage.Data.BookList, func(b types.Book) bool {
			return b.Path == player.currentBook.Path
		})
		if idx > -1 {
			savedPos = storage.Data.BookList[idx].Position
		}
	}
	targetBytePos := getTargetPositionInBytes(savedPos)
	err := player.channel.SetPosition(targetBytePos, bass.POS_BYTE)
	if err != nil {
		merror.ShowError("There seems to be a problem setting previous play position for this audio book", err)
	}
	return err
}

func getTargetPositionInBytes(targetPosition time.Duration) int {
	if targetPosition == 0 {
		return 0
	}

	chapter := 0

	bytePos, err := player.channel.Seconds2Bytes(targetPosition.Seconds())
	if err != nil {
		merror.ShowError("Invalid saved book position", err)
		return 0
	}

	// Determine the appropriate chapter to play, and byte position to start at in that chapter(file)
	{
		concatLength := float64(0)
		for _, c := range player.currentBook.Chapters {
			concatLength += c.LengthInSeconds
			if targetPosition.Seconds() > concatLength {
				b, err := player.channel.Seconds2Bytes(player.currentBook.Chapters[chapter].LengthInSeconds)
				if err != nil {
					merror.ShowError("Could not start from last play position, will play from the beginning of the audio book", err)
					return 0
				}
				bytePos = bytePos - b
				chapter++
			}
		}
	}

	// Load the appropriate file to play and set the appropriate position to start at
	player.currentBook.CurrentChapter = chapter
	loaded := loadAudioBookFile(storage.Data.Root + buildconstraints.PathSeparator + player.currentBook.Path + buildconstraints.PathSeparator + player.currentBook.Chapters[player.currentBook.CurrentChapter].FileName)
	newPos, err := player.channel.Bytes2Seconds(bytePos)
	if err != nil {
		merror.ShowError("Could not start from last play position, will play from the beginning of the audio book", err)
		return 0
	}
	player.currentBook.Position = time.Duration(newPos * 1000000000)

	if !loaded {
		bytePos = 0
	}

	return bytePos
}

func resetToBeginning() {
	if len(player.currentBook.Chapters) == 0 {
		return
	}
	player.currentBook.CurrentChapter = 0
	player.currentBook.Position = 0
	loadAudioBookFile(storage.Data.Root + buildconstraints.PathSeparator + player.currentBook.Path + buildconstraints.PathSeparator + player.currentBook.Chapters[0].FileName)
}

func skipToNextFile(p *Player, forceSkip bool, continuePlaying bool, setPreviousPosition bool) bool {
	skipped := false
	if p.channel != 0 {
		active, err := p.channel.IsActive()
		merror.ShowError("Error skipping to next chapter", err)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_PAUSED || forceSkip {
			numChapters := len(p.currentBook.Chapters)
			if numChapters > 0 {
				if p.currentBook.CurrentChapter < numChapters-1 {
					p.currentBook.CurrentChapter = p.currentBook.CurrentChapter + 1
					LoadAndPlay(p.currentBook, continuePlaying, setPreviousPosition, nil)
					skipped = true
				}
			}
		}
	}
	return skipped
}

//lint:ignore ST1011 honestly makes no sense
const leadingSeconds = 5 * time.Second

func goToBeginningOfFile(p *Player) bool {
	const errStr = "Error seeking to beginning of file"
	if p.channel != 0 {
		active, err := p.channel.IsActive()
		merror.ShowError(errStr, err)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_PAUSED {
			currentBytePosition, err := p.channel.GetPosition(bass.POS_BYTE)
			if err != nil {
				merror.ShowError(errStr, err)
				return false
			}

			currentSecondsPosition, err := p.channel.Bytes2Seconds(currentBytePosition)
			if err != nil {
				merror.ShowError(errStr, err)
				return false
			}
			if currentSecondsPosition >= float64(leadingSeconds.Seconds()) {
				err = p.channel.SetPosition(0, bass.POS_BYTE)
				if err != nil {
					merror.ShowError(errStr, err)
					return false
				}
				return true
			}
		}
	}
	return false
}

func skipToPreviousFile(p *Player, continuePlaying bool, setPreviousPosition bool) bool {
	skipped := false
	if p.channel != 0 {
		active, err := p.channel.IsActive()
		merror.ShowError("Error skipping to previous chapter", err)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_PAUSED {
			numChapters := len(p.currentBook.Chapters)
			if numChapters > 0 {
				if p.currentBook.CurrentChapter > 0 {
					p.currentBook.CurrentChapter = p.currentBook.CurrentChapter - 1
					LoadAndPlay(p.currentBook, continuePlaying, setPreviousPosition, nil)
					skipped = true
				} else {
					err = p.channel.SetPosition(0, bass.POS_BYTE)
					if err != nil {
						merror.ShowError("Error to skipping to start", err)
					}
				}
			}
		}
	}
	return skipped
}
