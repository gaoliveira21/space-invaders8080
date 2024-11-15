package io

import (
	"fmt"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	UFORepeatsSound byte = iota
	ShotSound
	ExplosionSound
	InvaderDieSound
	UFOHitSound
	FleetMovement1Sound
	FleetMovement2Sound
	FleetMovement3Sound
	FleetMovement4Sound
)

type SoundManager struct {
	sounds map[byte]*mix.Chunk
}

func NewSoundManager() *SoundManager {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		fmt.Printf("NewSoundManager: Could not init audio %s", err.Error())
		return nil
	}

	if err := mix.OpenAudio(44100, uint16(mix.DEFAULT_FORMAT), 2, 1024); err != nil {
		fmt.Printf("NewSoundManager: Could not open audio %s", err.Error())
		return nil
	}

	sm := &SoundManager{
		sounds: make(map[byte]*mix.Chunk),
	}

	soundFiles := map[byte]string{
		UFORepeatsSound:     "assets/ufo_lowpitch.wav",
		ShotSound:           "assets/shoot.wav",
		ExplosionSound:      "assets/explosion.wav",
		InvaderDieSound:     "assets/invaderkilled.wav",
		UFOHitSound:         "assets/ufo_highpitch.wav",
		FleetMovement1Sound: "assets/fastinvader1.wav",
		FleetMovement2Sound: "assets/fastinvader2.wav",
		FleetMovement3Sound: "assets/fastinvader3.wav",
		FleetMovement4Sound: "assets/fastinvader4.wav",
	}

	for id, file := range soundFiles {
		chunk, err := mix.LoadWAV(file)
		if err != nil {
			fmt.Printf("NewSoundManager: Could not load WAV file at %s %s", file, err.Error())
			return nil
		}
		sm.sounds[id] = chunk
	}

	return sm
}

func (sm *SoundManager) Play(soundId byte) {
	if chunk, ok := sm.sounds[soundId]; ok {
		if soundId >= FleetMovement1Sound {
			chunk.Volume(20)
		} else {
			chunk.Volume(10)
		}

		channel := -1
		switch soundId {
		case FleetMovement1Sound, FleetMovement2Sound, FleetMovement3Sound, FleetMovement4Sound:
			channel = 1
		}

		chunk.Play(channel, 0)
	}
}

func (sm *SoundManager) Cleanup() {
	for _, chunk := range sm.sounds {
		chunk.Free()
	}
	mix.CloseAudio()
}
