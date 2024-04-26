package audio

import (
	"mamela/merror"
	"mamela/types"
	"time"

	bass "github.com/pteich/gobass"
)

type Player struct {
	updater     chan types.PlayingBook
	currentBook types.PlayingBook
	channel     bass.Channel
	playing     bool
}

func (p *Player) play() {
	if p.channel != 0 {
		e := p.channel.Play(false)
		if e == nil {
			// Ticker = time.NewTicker(TickerDuration)
			Ticker.Reset(TickerDuration)
			p.playing = true
			// ChannelAudioState <- Playing
		}
		merror.ShowError("", e)
	}
}

func (p *Player) pause() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		if e != nil {
			merror.ShowError("", e)
		} else {
			if active == bass.ACTIVE_PLAYING {
				e := p.channel.Pause()
				Ticker.Stop()
				// ChannelAudioState <- Paused
				merror.ShowError("", e)
				merror.PanicError(e)
			}
		}
	}
}

func (p *Player) stop() {
	if p.channel != 0 {
		e := p.channel.Stop()
		if e != nil {
			merror.ShowError("", e)
		} else {
			Ticker.Stop()
			p.playing = false
			p.currentBook.Position = time.Duration(0)
			p.currentBook.CurrentChapter = 0
			p.channel.SetPosition(0, bass.POS_BYTE)
			updateUIOnStop()
			// updateUICurrentlyPlayingInfo()
			// ChannelAudioState <- Stopped
		}
	}
}

const fastForwardRewindAmount = 30

func (p *Player) fastRewind() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		merror.ShowError("", e)
		merror.PanicError(e)
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, e := p.channel.Seconds2Bytes(fastForwardRewindAmount)
			if e != nil {
				merror.ShowError("", e)
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
		merror.ShowError("", e)
		merror.PanicError(e)
		if active == bass.ACTIVE_PLAYING {
			bytePositionAmount, e := p.channel.Seconds2Bytes(fastForwardRewindAmount)
			if e != nil {
				merror.ShowError("", e)
			} else {
				currentBytePosition, e := p.channel.GetPosition(bass.POS_BYTE)
				if e != nil {
					merror.ShowError("", e)
				} else {
					byteLength, e := p.channel.GetLength(bass.POS_BYTE)
					if e == nil {
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
