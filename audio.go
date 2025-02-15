package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const (
	AUD_SAMPLE_RATE    = 16000
	AUD_SAMPLE_RATE_A  = 24000
	AUD_FILE_EXTENSION = ".ogg"
	AUD_FILE_SUBDIR    = "sounds"
	AUD_FILE_JUMP      = "jump"
	AUD_FILE_ATTACK    = "attack"
)

func (g *Game) initAudioPlayerHelper(soundID string) error {

	//check subdir exists
	if _, err := os.Stat(AUD_FILE_SUBDIR); errors.Is(err, os.ErrNotExist) {
		log.Fatal("Folder not found ", AUD_FILE_SUBDIR)
	}

	var err error

	filePath := path.Join(AUD_FILE_SUBDIR, soundID+AUD_FILE_EXTENSION)
	//check file exists
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("File not found ", filePath)
	} else {
		//log.Println("File  found ", jumpFilePath)
	}
	// get bytes array from file
	fileBytes, err := os.ReadFile(filePath)
	streamV, err := vorbis.DecodeF32(bytes.NewReader(fileBytes))
	if err != nil {
		log.Println("Audio: failed to decode ogg ", soundID)
		log.Fatal(err)
	}
	g.soundEffectPlayers[soundID], err = g.audioContext.NewPlayerF32(streamV)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (g *Game) initAudioPlayers() error {
	//get the canvas for playing sounds
	if g.audioContext == nil {
		g.audioContext = audio.NewContext(AUD_SAMPLE_RATE)
	}

	if nil == g.soundEffectPlayers {
		g.soundEffectPlayers = map[string]*audio.Player{}
	}

	g.initAudioPlayerHelper(AUD_FILE_JUMP)

	g.initAudioPlayerHelper(AUD_FILE_ATTACK)

	return nil

}
