package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 44100
)

var (
	audioContext         *audio.Context
	MusicPlayer          *audio.Player
	HitSoundPlayer       *audio.Player
	KillSoundPlayer      *audio.Player
	GameoverSoundPlayer  *audio.Player
	PlayerHitSoundPlayer *audio.Player
)

// for music
func CreateMusicSound(filePath string) error {
	var err error
	MusicPlayer, err = CreateSound(filePath)
	if err != nil {
		return err
	}

	return nil
}

// for sound effects
func CreateHitSound(filePath string) error {
	var err error
	HitSoundPlayer, err = CreateSound(filePath)
	if err != nil {
		return err
	}

	return nil
}

// for sound effects
func CreateKillSound(filePath string) error {
	var err error
	KillSoundPlayer, err = CreateSound(filePath)
	if err != nil {
		return err
	}

	return nil
}

func CreatePlayerHitSound(filePath string) error {
	var err error
	PlayerHitSoundPlayer, err = CreateSound(filePath)
	if err != nil {
		return err
	}

	return nil
}

func CreateGameoverSound(filePath string) error {
	var err error
	GameoverSoundPlayer, err = CreateSound(filePath)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) SetVolume(increase bool) {
	if increase {
		g.globalVolume += 0.1
	} else {
		g.globalVolume -= 0.1
	}
	if g.globalVolume > 1 {
		g.globalVolume = 1
	} else if g.globalVolume < 0 {
		g.globalVolume = 0
	}

	MusicPlayer.SetVolume(g.globalVolume)
	HitSoundPlayer.SetVolume(g.globalVolume)
	KillSoundPlayer.SetVolume(g.globalVolume)
	GameoverSoundPlayer.SetVolume(g.globalVolume)
	PlayerHitSoundPlayer.SetVolume(g.globalVolume)
}

func SetVolumeValue(volume float64) {
	MusicPlayer.SetVolume(volume)
	HitSoundPlayer.SetVolume(volume)
	KillSoundPlayer.SetVolume(volume)
	GameoverSoundPlayer.SetVolume(volume)
	PlayerHitSoundPlayer.SetVolume(volume)

}

func MusicLoop() {
	if !MusicPlayer.IsPlaying() {
		MusicPlayer.Rewind()
	}
	if !HitSoundPlayer.IsPlaying() {
		HitSoundPlayer.Rewind()
	}
	if !KillSoundPlayer.IsPlaying() {
		KillSoundPlayer.Rewind()
	}
	if !PlayerHitSoundPlayer.IsPlaying() {
		PlayerHitSoundPlayer.Rewind()
	}
	if !GameoverSoundPlayer.IsPlaying() {
		GameoverSoundPlayer.Rewind()
	}
}

func CreateSound(filePath string) (*audio.Player, error) {
	// Inicializar o contexto de áudio
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}

	// Ler o arquivo de música
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	//get file extension
	ext := filepath.Ext(filePath)
	// fmt.Println("File extension: ", ext)

	if ext == ".ogg" {
		d, err := vorbis.DecodeWithSampleRate(sampleRate, f)
		if err != nil {
			return nil, err
		}

		player, err := audioContext.NewPlayer(audio.NewInfiniteLoop(d, d.Length()))
		if err != nil {
			return nil, err
		}

		return player, nil

	} else if ext == ".wav" {
		d, err := wav.DecodeWithSampleRate(sampleRate, f)
		if err != nil {
			return nil, err
		}
		// f.Close()

		// Criar um player de áudio
		player, err := audioContext.NewPlayer(d)
		if err != nil {
			return nil, err
		}

		return player, nil
	} else {
		log.Fatal("Audio File extension not supported")
	}
	// Decodificar o arquivo de música

	return nil, nil
}
