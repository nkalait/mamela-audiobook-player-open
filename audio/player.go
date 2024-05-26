package audio

import (
	"mamela/merror"
	"mamela/storage"
	"mamela/types"
	"time"

	bass "github.com/pteich/gobass"
)

type Player struct {
	currentBook types.PlayingBook
	channel     bass.Channel
	playing     bool
	State       int
}

func (p *Player) play() {
	if p.channel != 0 {
		p.setVolume(storage.GetVolumeLevel())
		err := p.channel.Play(false)
		if err == nil {
			UIUpdateTicker.Reset(PlayingBookTickerDuration)
			CurrentBookPositionUpdateTicker.Reset(CurrentBookPositionTickerDuration)
			p.playing = true
			p.State = PLAYING
			storage.UpdateCurrentBook(p.currentBook.Path)
			updateUICurrentlyPlayingInfo()
		}
		merror.ShowError("Cannot play", err)
	}
}

func (p *Player) pause() {
	if player.channel != 0 {
		active, err := player.channel.IsActive()
		if err == nil {
			if active == bass.ACTIVE_PLAYING {
				saveCurrentPlayingBookPositionToDisk()
				err := p.channel.Pause()
				merror.ShowError("Error pausing", err)
				p.playing = false
				p.State = PAUSED
				UIUpdateTicker.Stop()
				CurrentBookPositionUpdateTicker.Stop()
				updateUICurrentlyPlayingInfo()
			}
		}
	}
}

func (p *Player) stop() {
	if p.channel != 0 {
		err := p.channel.Stop()
		if err == nil {
			saveCurrentPlayingBookPositionToDisk()
			UIUpdateTicker.Stop()
			CurrentBookPositionUpdateTicker.Stop()
			p.playing = false
			p.State = STOPPED
			p.currentBook.Position = time.Duration(0)
			p.currentBook.CurrentChapter = 0
			p.channel.SetPosition(0, bass.POS_BYTE)
			updateUIOnStop()
			updateUICurrentlyPlayingInfo()
		}
	}
}

const fastForwardRewindAmount = 30

func (p *Player) fastRewind() {
	if player.channel != 0 {
		active, err := player.channel.IsActive()
		if err != nil {
			merror.ShowError("", err)
			return
		}
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, err := p.channel.Seconds2Bytes(fastForwardRewindAmount)
			if err != nil {
				merror.ShowError("", err)
			} else {
				currentBytePosition, err := p.channel.GetPosition(bass.POS_BYTE)
				if err == nil {
					if currentBytePosition-bytePositionAmount < 0 {
						if skipToPreviousFile(p, true, false) {
							completeFileByteLength, err := p.channel.GetLength(bass.POS_BYTE)
							if err == nil {
								newPos := completeFileByteLength - bytePositionAmount
								if currentBytePosition < bytePositionAmount {
									deductBy := bytePositionAmount - currentBytePosition
									newPos = completeFileByteLength - deductBy
								}
								p.channel.SetPosition(newPos, bass.POS_BYTE)
							}
						} else {
							p.channel.SetPosition(0, bass.POS_BYTE)
						}
					} else {
						p.channel.SetPosition(currentBytePosition-bytePositionAmount, bass.POS_BYTE)
					}
				}
			}
		}
		updateUICurrentlyPlayingInfo()
	}
}

func (p *Player) fastForward() {
	if player.channel != 0 {
		active, _ := player.channel.IsActive()
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, err := p.channel.Seconds2Bytes(fastForwardRewindAmount)
			if err != nil {
				merror.ShowError("Cannot convert seconds to bytes", err)
			} else {
				currentBytePosition, err := p.channel.GetPosition(bass.POS_BYTE)
				if err != nil {
					merror.ShowError("Cannot get byte position", err)
				} else {
					byteLength, err := p.channel.GetLength(bass.POS_BYTE)
					if err == nil {
						if currentBytePosition+bytePositionAmount >= byteLength {
							if !skipToNextFile(p, false, true, false) {
								p.channel.SetPosition(byteLength, bass.POS_BYTE)
							}
						} else {
							p.channel.SetPosition(currentBytePosition+bytePositionAmount, bass.POS_BYTE)
						}
					}

				}
			}
		}
		updateUICurrentlyPlayingInfo()
	}
}

func (p *Player) skipNext() {
	skipToNextFile(p, false, p.playing, false)
}

func (p *Player) skipPrevious() {
	if goToBeginningOfFile(p) {
		return
	}
	skipToPreviousFile(p, p.playing, false)
}

func (p *Player) setVolume(vol float64) {
	bass.SetConfig(bass.CONFIG_GVOL_STREAM, int(vol))
}

func GetVolume() int64 {
	vol, _ := bass.GetConfig(bass.CONFIG_GVOL_STREAM)
	return vol
}

func ClearPlayer() {
	if player.channel != 0 {
		player.channel.Free()
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
