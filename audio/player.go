package audio

import (
	"mamela/err"
	"mamela/types"

	bass "github.com/pteich/gobass"
)

type Player struct {
	updater     chan types.PlayingBook
	currentBook types.PlayingBook
	channel     bass.Channel
}

func (p *Player) play() {
	if p.channel != 0 {
		e := p.channel.Play(false)
		if e == nil {
			Ticker.Reset(TickerDuration)
			ChannelAudioState <- Playing
		}
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
				Ticker.Stop()
				ChannelAudioState <- Paused
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
			p.currentBook.CurrentChapter = 0
			p.channel.SetPosition(0, bass.POS_BYTE)
			updateUICurrentlyPlayingInfo()
			Ticker.Stop()
			ChannelAudioState <- Stopped
		}
	}
}

func (p *Player) fastRewind() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		err.PanicError(e)
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, e := p.channel.Seconds2Bytes(60)
			if e != nil {
				err.ShowError("", e)
			} else {
				currentBytePosition, e := p.channel.GetPosition(bass.POS_BYTE)
				if e == nil {
					if currentBytePosition-bytePositionAmount < 0 {
						if skipToPreviousFile(p) {
							completeFileByteLength, e := p.channel.GetLength(bass.POS_BYTE)
							if e == nil {
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
		active, e := player.channel.IsActive()
		err.PanicError(e)
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, e := p.channel.Seconds2Bytes(60)
			if e != nil {
				err.ShowError("", e)
			} else {
				currentBytePosition, e := p.channel.GetPosition(bass.POS_BYTE)
				if e != nil {
					err.ShowError("", e)
				} else {
					byteLength, e := p.channel.GetLength(bass.POS_BYTE)
					if e == nil {
						if currentBytePosition+bytePositionAmount >= byteLength {
							if !skipToNextFile(p) {
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
	skipToNextFile(p)
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
