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

// for music
func (g *Game) PlayOGGSound(filePath string) error {
	// Inicializar o contexto de áudio
	if g.audioContext == nil {
		g.audioContext = audio.NewContext(sampleRate)
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
	g.musicPlayer, err = g.audioContext.NewPlayer(audio.NewInfiniteLoop(d, d.Length()))
	if err != nil {
		return err
	}

	// Tocar a música em loop
	g.musicPlayer.SetVolume(0.1)
	g.musicPlayer.Play()

	return nil
}

// for sound effects
func (g *Game) PlayWAVSound(filePath string) error {
	// Inicializar o contexto de áudio
	if g.audioContext == nil {
		g.audioContext = audio.NewContext(sampleRate)
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
	g.musicPlayer, err = g.audioContext.NewPlayer(d)
	if err != nil {
		return err
	}

	// Tocar a música em loop
	g.musicPlayer.SetVolume(0.1)
	g.musicPlayer.Play()

	return nil
}
