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
}

func (p *Player) play() {
	if p.channel != 0 {
		err := p.channel.Play(false)
		if err == nil {
			UIUpdateTicker.Reset(PlayingBookTickerDuration)
			CurrentBookPositionUpdateTicker.Reset(CurrentBookPositionTickerDuration)
			p.playing = true
			storage.UpdateCurrentBook(p.currentBook.Path)
		}
		merror.ShowError("", err)
	}
}

func (p *Player) pause() {
	if player.channel != 0 {
		active, err := player.channel.IsActive()
		if err != nil {
			merror.ShowError("", err)
		} else {
			if active == bass.ACTIVE_PLAYING {
				err := p.channel.Pause()
				UIUpdateTicker.Stop()
				CurrentBookPositionUpdateTicker.Stop()
				merror.ShowError("", err)
				merror.PanicError(err)
			}
		}
	}
}

func (p *Player) stop() {
	if p.channel != 0 {
		err := p.channel.Stop()
		if err != nil {
			merror.ShowError("", err)
		} else {
			UIUpdateTicker.Stop()
			CurrentBookPositionUpdateTicker.Stop()
			p.playing = false
			p.currentBook.Position = time.Duration(0)
			p.currentBook.CurrentChapter = 0
			p.channel.SetPosition(0, bass.POS_BYTE)
			updateUIOnStop()
			saveCurrentPlayingBookPositionToDisk()
			// updateUICurrentlyPlayingInfo()
			// ChannelAudioState <- Stopped
		}
	}
}

const fastForwardRewindAmount = 30

func (p *Player) fastRewind() {
	if player.channel != 0 {
		active, err := player.channel.IsActive()
		merror.ShowError("", err)
		merror.PanicError(err)
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, err := p.channel.Seconds2Bytes(fastForwardRewindAmount)
			if err != nil {
				merror.ShowError("", err)
			} else {
				currentBytePosition, err := p.channel.GetPosition(bass.POS_BYTE)
				if err == nil {
					if currentBytePosition-bytePositionAmount < 0 {
						if skipToPreviousFile(p) {
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
		active, err := player.channel.IsActive()
		merror.ShowError("", err)
		merror.PanicError(err)
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, err := p.channel.Seconds2Bytes(fastForwardRewindAmount)
			if err != nil {
				merror.ShowError("", err)
			} else {
				currentBytePosition, err := p.channel.GetPosition(bass.POS_BYTE)
				if err != nil {
					merror.ShowError("", err)
				} else {
					byteLength, err := p.channel.GetLength(bass.POS_BYTE)
					if err == nil {
						if currentBytePosition+bytePositionAmount >= byteLength {
							if !skipToNextFile(p, false) {
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
	skipToNextFile(p, false)
}

func (p *Player) skipPrevious() {
	skipToPreviousFile(p)
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
