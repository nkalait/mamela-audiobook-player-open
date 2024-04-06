package audio

import (
	"mamela/err"
	"mamela/types"
	"os"

	bass "github.com/pteich/gobass"
)

type Player struct {
	updater     chan types.PlayingBook
	currentBook types.PlayingBook
	channel     bass.Channel
}

func (p *Player) getCurrentFile() *os.File {
	path := p.currentBook.FullPath + "/" + p.currentBook.Chapters[p.currentBook.CurrentChapter]
	f, _ := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	return f
}

func (p *Player) play() {
	if p.channel != 0 {
		e := p.channel.Play(false)
		err.ShowError("", e)
	}
}

func (p *Player) pause() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		if e != nil {
			err.ShowError("", e)
		} else {
			if active == bass.ACTIVE_PLAYING {
				e := p.channel.Pause()
				err.PanicError(e)
			}
		}
	}
}

func (p *Player) stop() {
	if p.channel != 0 {
		e := p.channel.Stop()
		if e != nil {
			err.ShowError("", e)
		} else {
			p.channel.SetPosition(0, bass.POS_BYTE)
			updateUICurrentlyPlayingInfo()
		}
	}
}

func (p *Player) fastRewind() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		err.PanicError(e)
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, e := p.channel.Seconds2Bytes(10)
			if e != nil {
				err.ShowError("", e)
			} else {
				currentBytePosition, e := p.channel.GetPosition(bass.POS_BYTE)
				if e != nil {
					err.ShowError("", e)
				} else {
					if currentBytePosition-bytePositionAmount < 0 {
						e = p.channel.SetPosition(0, bass.POS_BYTE)
						if e != nil {
							err.ShowError("Error setting rewinding", e)
						}
					} else {
						e = p.channel.SetPosition(currentBytePosition-bytePositionAmount, bass.POS_BYTE)
						if e != nil {
							err.ShowError("Error setting rewinding", e)
						}
					}
				}
			}
		}
	}
}

func (p *Player) fastForward() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		err.PanicError(e)
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, e := p.channel.Seconds2Bytes(10)
			if e != nil {
				err.ShowError("", e)
			} else {
				currentBytePosition, e := p.channel.GetPosition(bass.POS_BYTE)
				if e != nil {
					err.ShowError("", e)
				} else {
					byteLength, e := p.channel.GetLength(bass.POS_BYTE)
					if e != nil {
						err.ShowError("Error setting fast forwarding", e)
					}
					if currentBytePosition+bytePositionAmount >= byteLength {
						e = p.channel.SetPosition(byteLength, bass.POS_BYTE)
						if e != nil {
							err.ShowError("Error setting fast forwarding", e)
						}
					} else {
						e = p.channel.SetPosition(currentBytePosition+bytePositionAmount, bass.POS_BYTE)
						if e != nil {
							err.ShowError("Error setting fast forwarding", e)
						}
					}
				}
			}
		}
	}
}

func (p *Player) skipNext() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		err.ShowError("Error skipping to next chapter", e)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_PAUSED {
			numChapters := len(player.currentBook.Chapters)
			if numChapters > 0 {
				if player.currentBook.CurrentChapter < numChapters-1 {
					player.currentBook.CurrentChapter = player.currentBook.CurrentChapter + 1
					LoadAndPlay(p.currentBook, nil)
				}
			}
		}
	}
}

func (p *Player) skipPrevious() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		err.ShowError("Error skipping to previous chapter", e)
		if active == bass.ACTIVE_PLAYING || active == bass.ACTIVE_PAUSED {
			numChapters := len(player.currentBook.Chapters)
			if numChapters > 0 {
				if player.currentBook.CurrentChapter > 0 {
					player.currentBook.CurrentChapter = player.currentBook.CurrentChapter - 1
					LoadAndPlay(p.currentBook, nil)
				} else {
					e = p.channel.SetPosition(0, bass.POS_BYTE)
					if e != nil {
						err.ShowError("Error to skipping to start", e)
					}
				}
			}
		}
	}
}

func Play() {
	player.play()
}
func Pause() {
	player.pause()
}
func Stop() {
	player.stop()
}
func FastRewind() {
	player.fastRewind()
}
func FastForward() {
	player.fastForward()
}
func SkipNext() {
	player.skipNext()
}
func SkipPrevious() {
	player.skipPrevious()
}
