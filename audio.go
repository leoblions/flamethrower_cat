package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

/*

https://ebitengine.org/en/examples/audio.html

*/

const (
	AU_SAMPLE_RATE      = 48000
	AU_BYTES_PER_SAMPLE = 8
)

var (
// au_audio_stream audioStream
)

type soundFormat int

const (
	formatOgg soundFormat = iota
	formatMp3
)

type Audio struct {
	game         *Game
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
	soundFormat  soundFormat
}

func NewAudio(game *Game) *Audio {
	au := &Audio{}
	return au
}
