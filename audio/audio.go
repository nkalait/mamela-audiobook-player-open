package audio

import (
	"fmt"
	"mamela/buildConstraints"
	"mamela/merror"
	"mamela/storage"
	"mamela/types"
	"os"
	"slices"
	"time"

	bass "github.com/pteich/gobass"
)

var LibDir = "lib" + buildConstraints.PathSeparator + "mac"

// Event listeners
var (
	ExitListener      = make(chan bool) // for stopping to listen to playing events
	exitAudio         = make(chan bool) // for unloading audio stuff
	BassInitiatedChan = make(chan bool)
)

// const (
// 	Stopped = iota
// 	Paused
// 	Playing
// )

// var ChannelAudioState = make(chan int)
var UpdateNowPlayingChannel = make(chan types.PlayingBook)

// Listens to events about changes to audio books root folder
var UpdateBookListChannel = make(chan bool)

// Holds data structures important to playing an audio book
var player Player

// half a second delay before updating UI
const PlayingBookTickerDuration = 500 * time.Millisecond

var UIUpdateTicker *time.Ticker = time.NewTicker(PlayingBookTickerDuration)

// 30 second delay before saving currently playing books
// play position to disk
const CurrentBookPositionTickerDuration = time.Second * 30

var CurrentBookPositionUpdateTicker *time.Ticker = time.NewTicker(CurrentBookPositionTickerDuration)

// Initiate Bass
func init() {
	UIUpdateTicker.Stop()
	CurrentBookPositionUpdateTicker.Stop()
	go func() {
		plugins := loadPlugins()
		initBass()
		<-exitAudio
		tearDown(plugins)
	}()

	go func() {
		for {
			select {
			case <-CurrentBookPositionUpdateTicker.C:
				saveCurrentPlayingBookPositionToDisk()
			}
		}
	}()

}

func saveCurrentPlayingBookPositionToDisk() {
	if len(storage.Data.BookList) > 0 {
		idx := slices.IndexFunc(storage.Data.BookList, func(b types.Book) bool {
			return b.Path == player.currentBook.Path
		})
		if idx > -1 {
			storage.Data.BookList[idx].Position = GetCurrentBookPlayingDuration(player.currentBook)
			storage.SaveBookListToStorageFile(storage.Data.BookList)
		}
	}
}

// Unload loaded Bass plugins and free all resources used by Bass
func tearDown(plugins []uint32) {
	for _, p := range plugins {
		bass.PluginFree(p)
	}
	bass.Free()
	os.Exit(0)
}

// Initialise Bass
func initBass() {
	err := bass.Init(-1, 44100, bass.DeviceStereo, 0, 0)
	merror.ShowError("Problem initiating bass", err)
	merror.PanicError(err)
	bass.SetVolume(100)
	BassInitiatedChan <- true
}

// Load plugins needed by Bass
func loadPlugins() []uint32 {
	pluginLibbassAac, err := bass.PluginLoad(LibDir+buildConstraints.PathSeparator+"libbass_aac.dylib", bass.StreamDecode)
	merror.ShowError("Problem loading plugin", err)
	merror.PanicError(err)
	pluginLibbassOpus, err := bass.PluginLoad(LibDir+buildConstraints.PathSeparator+"libbassopus.dylib", bass.StreamDecode)
	merror.ShowError("Problem loading plugin", err)
	merror.PanicError(err)

	plugins := make([]uint32, 2)
	plugins = append(plugins, pluginLibbassAac)
	plugins = append(plugins, pluginLibbassOpus)

	return plugins
}

// Start listening to audio playing event and exit event
func StartChannelListener(exitApp chan bool) {
	go func() {
	RoutineLoop:
		for {
			select {
			// case <-time.After(time.Second):
			case <-UIUpdateTicker.C:
				updateUICurrentlyPlayingInfo()
			case <-ExitListener:
				break RoutineLoop
			}
		}
		exitAudio <- true
		time.Sleep(time.Second * 2)
		exitApp <- true
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

func GetCurrentBookPlayingDuration(p types.PlayingBook) time.Duration {
	pos := p.Position
	if p.CurrentChapter > 0 {
		for i := p.CurrentChapter - 1; i >= 0; i-- {
			pos = pos + time.Duration(p.Chapters[i].LengthInSeconds*1000000000)
		}
	}
	return pos
}
func updateUIOnStop() {
	UpdateNowPlayingChannel <- player.currentBook
}

// Update the currently playing audio book information on the UI
func updateUICurrentlyPlayingInfo() {
	if player.channel != 0 {
		active, err := player.channel.IsActive()
		merror.ShowError("", err)
		merror.PanicError(err)

		// We need active == bass.ACTIVE_STOPPED here in order to detect when
		// file has reached end
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_STOPPED {
			bytePosition, err := player.channel.GetPosition(bass.POS_BYTE)
			merror.ShowError("", err)
			merror.PanicError(err)
			pos, err := player.channel.Bytes2Seconds(bytePosition)
			merror.ShowError("", err)
			merror.PanicError(err)

			currentlyAt := player.currentBook.Position.Round(time.Second)
			skipAt := time.Duration(player.currentBook.Chapters[player.currentBook.CurrentChapter].LengthInSeconds * 1000000000).Round(time.Second)
			if currentlyAt == skipAt {
				skipToNextFile(&player, true)
			}

			posInWholeBook := GetCurrentBookPlayingDuration(player.currentBook).Round(time.Second)
			wholeBookLength := time.Duration(player.currentBook.FullLengthSeconds * 1000000000).Round(time.Second)
			if posInWholeBook == wholeBookLength {
				player.currentBook.Finished = true
				UIUpdateTicker.Stop()
				CurrentBookPositionUpdateTicker.Stop()
				// ChannelAudioState <- Stopped
			}

			var d time.Duration = time.Duration(pos * 1000000000)
			player.currentBook.Position = time.Duration(d)
		}
		UpdateNowPlayingChannel <- player.currentBook
	}
}

type UpdateFolderArtCallBack func(playingBook types.PlayingBook)

func LoadAndPlay(playingBook types.PlayingBook, resumePlayback bool, updaterFolderArtCallback UpdateFolderArtCallBack) {
	// c, err := bass.StreamCreateURL("http://music.myradio.ua:8000/PopRock_news128.mp3", bass.DeviceStereo)
	stopPlayingIfPlaying(player.channel, player)
	player.currentBook = playingBook

	chapter := player.currentBook.CurrentChapter
	err := loadAudioBookFile(storage.Data.Root + buildConstraints.PathSeparator + player.currentBook.Path + buildConstraints.PathSeparator + player.currentBook.Chapters[chapter].FileName)
	if err == nil {
		startPlaying(resumePlayback)
	}

	if updaterFolderArtCallback != nil {
		updaterFolderArtCallback(player.currentBook)
	}
	updateUICurrentlyPlayingInfo()
}

func stopPlayingIfPlaying(c bass.Channel, p Player) {
	if c != 0 {
		a, err := c.IsActive()
		merror.ShowError("", err)
		merror.PanicError(err)
		if a == bass.ACTIVE_PLAYING || a == bass.ACTIVE_PAUSED {
			p.stop()
		}
	}
}

func loadAudioBookFile(fullPath string) error {
	var err error = nil
	player.channel, err = bass.StreamCreateFile(fullPath, 0, bass.AsyncFile)
	if err != nil {
		merror.ShowError("There seems to be a problem loading the the audio book file(s)", err)
	}

	return err
}

func startPlaying(resumePlayback bool) error {
	bytePos := 0
	if resumePlayback {
		bytePos = skipToLastPosition()
	}
	err := player.channel.SetPosition(bytePos, bass.POS_BYTE)
	if err != nil {
		merror.ShowError("There seems to be a problem playing the the audio book file(s)", err)
	} else {
		player.play()
	}
	return err
}

// We want to use the position saved to disk here so that we can resume playback
func skipToLastPosition() int {
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

	// Byte position of the last play position saved on disk
	bytePos, _ := player.channel.Seconds2Bytes(float64(savedPos.Seconds()))
	// Determine the last chapter that was playing while also decrementing bytePos
	// by the concatenated lengths of all the chapters that have played
	{
		concatLength := float64(0)
		for _, c := range player.currentBook.Chapters {
			concatLength += c.LengthInSeconds
			if savedPos.Seconds() > concatLength {
				b, err := player.channel.Seconds2Bytes(player.currentBook.Chapters[player.currentBook.CurrentChapter].LengthInSeconds)
				if err != nil {
					merror.ShowError("Could not start from last play position, will play from the beginning of the audio book", err)
					return 0
				}
				bytePos = bytePos - b
				player.currentBook.CurrentChapter++
			}
		}
	}

	// If the audio book has played pass at least one chapter, ie one file,
	// then load the appropriate file to play and set the appropriate position to start at
	if player.currentBook.CurrentChapter > 0 {
		newPos, err := player.channel.Bytes2Seconds(bytePos)
		if err != nil {
			merror.ShowError("Could not start from last play position, will play from the beginning of the audio book", err)
			return 0
		}
		player.currentBook.Position = time.Duration(newPos * 1000000000)
		loadAudioBookFile(storage.Data.Root + buildConstraints.PathSeparator + player.currentBook.Path + buildConstraints.PathSeparator + player.currentBook.Chapters[player.currentBook.CurrentChapter].FileName)
	}

	return bytePos
}

func skipToNextFile(p *Player, forceSkip bool) bool {
	skipped := false
	if p.channel != 0 {
		active, err := p.channel.IsActive()
		merror.ShowError("Error skipping to next chapter", err)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_PAUSED || forceSkip {
			numChapters := len(p.currentBook.Chapters)
			if numChapters > 0 {
				if p.currentBook.CurrentChapter < numChapters-1 {
					p.currentBook.CurrentChapter = p.currentBook.CurrentChapter + 1
					LoadAndPlay(p.currentBook, false, nil)
					skipped = true
				}
			}
		}
	}
	return skipped
}

func skipToPreviousFile(p *Player) bool {
	skipped := false
	if p.channel != 0 {
		active, err := p.channel.IsActive()
		merror.ShowError("Error skipping to previous chapter", err)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_PAUSED {
			numChapters := len(p.currentBook.Chapters)
			if numChapters > 0 {
				if p.currentBook.CurrentChapter > 0 {
					p.currentBook.CurrentChapter = p.currentBook.CurrentChapter - 1
					LoadAndPlay(p.currentBook, false, nil)
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
