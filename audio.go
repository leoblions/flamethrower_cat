package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const (
	AUD_SAMPLE_RATE        = 16000
	AUD_SAMPLE_RATE_A      = 24000
	AUD_FILE_EXTENSION     = ".ogg"
	AUD_FILE_SUBDIR        = "sounds"
	AUD_FILE_JUMP          = "jump"
	AUD_FILE_ATTACK        = "attack"
	AUD_FILE_DOOR          = "door"
	AUD_FILE_HIT           = "hit"
	AUD_FILE_LAVAHISS      = "lavahiss"
	AUD_FILE_DOOROPEN      = "dooropen"
	AUD_FILE_DOORCLOSE     = "doorclose"
	AUD_FILE_CANLID_REVERB = "canlid_reverb"
	AUD_DO_ASYNC           = true
)

type AudioPlayer struct {
	game             *Game
	waitgroup        sync.WaitGroup
	audioContext     *audio.Context
	soundAffectNames []string
	//jumpPlayer   *audio.Player
	//hitPlayer    *audio.Player
	soundEffectPlayers map[string]*audio.Player
}

func NewAudioPlayer(game *Game) *AudioPlayer {
	ap := &AudioPlayer{}
	ap.soundAffectNames = []string{
		"jump",
		"attack",
		"hit",
		"lavahiss",
		"dooropen",
		"doorclose",
		"canlid_reverb",
	}
	ap.initAudioPlayers()

	if AUD_DO_ASYNC {
		ap.waitgroup = sync.WaitGroup{}
	}

	return ap
}

func (ap *AudioPlayer) playSoundPlayer(plr *audio.Player) error {
	if err := plr.Rewind(); err != nil {
		return err
	}
	plr.Play()
	return nil
}

func (ap *AudioPlayer) playSoundByID(soundID string) error {
	soundEffect := ap.soundEffectPlayers[soundID]
	if AUD_DO_ASYNC {
		ap.waitgroup.Add(1)
		go func() {
			defer ap.waitgroup.Done()
			ap.playSoundPlayer(soundEffect)
		}()
	} else {
		ap.playSoundPlayer(soundEffect)
	}
	return nil
}

func (ap *AudioPlayer) playSoundByID_0(soundID string) error {
	soundEffect := ap.soundEffectPlayers[soundID]
	//fmt.Println("Play sound ,", soundID)
	if err := soundEffect.Rewind(); err != nil {
		return err
	}
	soundEffect.Play()
	return nil
}

func (ap *AudioPlayer) initAudioPlayerHelper(soundID string) error {
	if nil == ap.soundEffectPlayers {
		ap.soundEffectPlayers = map[string]*audio.Player{}
	}

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
	ap.soundEffectPlayers[soundID], err = ap.audioContext.NewPlayerF32(streamV)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Audio: added ogg ", filePath)
	//log.Println("Audio:  ", filePath)
	return nil
}

func (ap *AudioPlayer) initAudioPlayers() error {
	//get the canvas for playing sounds
	if ap.audioContext == nil {
		ap.audioContext = audio.NewContext(AUD_SAMPLE_RATE)
	}

	if nil == ap.soundEffectPlayers {
		ap.soundEffectPlayers = map[string]*audio.Player{}
	}

	for _, value := range ap.soundAffectNames {
		ap.initAudioPlayerHelper(value)
	}

	//ap.initAudioPlayerHelper(AUD_FILE_JUMP)

	//ap.initAudioPlayerHelper(AUD_FILE_ATTACK)

	return nil

}
