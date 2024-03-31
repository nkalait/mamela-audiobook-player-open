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
		err.PanicError(e)
	}
}

func (p *Player) pause() {
	if player.channel != 0 {
		active, e := player.channel.IsActive()
		err.PanicError(e)
		if active == bass.ACTIVE_PLAYING {
			e := p.channel.Pause()
			err.PanicError(e)
		}
	}
}

func (p *Player) stop() {
	if p.channel != 0 {
		e := p.channel.Stop()
		err.PanicError(e)
		p.channel.SetPosition(0, bass.POS_BYTE)
		updateUIPlayingPosition(0)
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
