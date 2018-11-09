/*
 * @Author: V4 Games
 * @Date: 2018-11-09 02:23:41
 * @Last Modified by:   Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-11-09 02:23:41
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// Shuffle marks whether the music should be shuffled
	Shuffle     bool
	trackNames  []string
	tracks      map[string]rl.Music
	track       rl.Music
	trackIndex  int
	musicVolume float32
)

// LoadPlaylist loads a playlist from file
func LoadPlaylist(name string) {
	trackNames = []string{}
	tracks = make(map[string]rl.Music)

	data, err := ioutil.ReadFile(fmt.Sprintf("assets/music/%s", name))

	if err != nil {
		log.Fatalf("Cannot load playlist '%s': \n\t%s\n", name, err.Error())
		return
	}

	txt := string(data)

	trackNames = strings.Split(txt, "\n")

	trackIndex = 0
	LoadNextTrack()
}

// LoadNextTrack loads the next track in the playlist
func LoadNextTrack() {
	trackName := strings.TrimSpace(trackNames[trackIndex])

	if Shuffle {
		trackIndex = rand.Int() % len(trackNames)
	} else {
		trackIndex++

		if trackIndex >= len(trackNames) {
			trackIndex = 0
		}
	}

	if trackName == "" {
		return
	}

	st, ok := tracks[trackName]

	if !ok {
		tr := rl.LoadMusicStream(fmt.Sprintf("assets/music/%s", trackName))
		log.Printf("Loading track: %s!", trackName)
		tracks[trackName] = tr
		st = tr
	}

	track = st

	rl.StopMusicStream(track)
	rl.PlayMusicStream(track)
	rl.SetMusicVolume(track, musicVolume)
}

// SetMusicVolume sets the music volume
func SetMusicVolume(vol float32) {
	musicVolume = vol

	if rl.IsMusicPlaying(track) {
		rl.SetMusicVolume(track, musicVolume)
	}
}

// UpdateMusic checks for music playback and updates the tracklist
func UpdateMusic() {
	if rl.IsMusicPlaying(track) {
		rl.UpdateMusicStream(track)
		return
	}

	LoadNextTrack()
}

// PauseMusic pauses the music stream
func PauseMusic() {
	rl.PauseMusicStream(track)
}
