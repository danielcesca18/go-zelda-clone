package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 44100
)

var (
	audioContext   *audio.Context
	MusicPlayer    *audio.Player
	HitSoundPlayer *audio.Player
)

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

}

func SetVolumeValue(volume float64) {
	MusicPlayer.SetVolume(volume)
	HitSoundPlayer.SetVolume(volume)

}

func MusicLoop() {
	if !MusicPlayer.IsPlaying() {
		MusicPlayer.Rewind()
	}
	if !HitSoundPlayer.IsPlaying() {
		HitSoundPlayer.Rewind()
	}
}

// for music
func CreateMusicSound(filePath string) error {
	// Inicializar o contexto de áudio
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}

	// Ler o arquivo de música
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	// Criar um ReadSeekCloser a partir do buffer
	// stream := audio.NewInfiniteLoop(bytes.NewReader(data), int64(len(data)))

	// Decodificar o arquivo de música
	d, err := vorbis.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		return err
	}
	// f.Close()

	// Criar um player de áudio
	MusicPlayer, err = audioContext.NewPlayer(audio.NewInfiniteLoop(d, d.Length()))
	if err != nil {
		return err
	}

	// // Tocar a música em loop
	// MusicPlayer.SetVolume(*volume)
	// g.musicPlayer.Play()

	return nil
}

// for sound effects
func CreateHitSound(filePath string) error {
	// Inicializar o contexto de áudio
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}

	// Ler o arquivo de música
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	// Decodificar o arquivo de música
	d, err := wav.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		return err
	}
	// f.Close()

	// Criar um player de áudio
	HitSoundPlayer, err = audioContext.NewPlayer(d)
	if err != nil {
		return err
	}

	// // Tocar a música em loop
	// musicPlayer.SetVolume(0.1)
	// musicPlayer.Play()

	return nil
}
