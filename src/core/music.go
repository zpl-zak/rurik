/*
   Copyright 2019 Dominik Madarász <zaklaus@madaraszd.net>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package core

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

var (
	// Shuffle marks whether the music should be shuffled
	Shuffle     bool
	trackNames  []string
	tracks      map[string]rl.Music
	track       *rl.Music
	trackIndex  int
	musicVolume float32
)

// LoadPlaylist loads a playlist from file
func LoadPlaylist(name string) {
	trackNames = []string{}
	tracks = make(map[string]rl.Music)

	txt := string(system.GetFile(fmt.Sprintf("music/%s", name), true))

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
		fln := fmt.Sprintf("music/%s", trackName)
		tr := rl.LoadMusicStreamFromMemory(string(system.GetFile(fln, true)))
		log.Printf("Loading track: %s!", trackName)
		tracks[trackName] = tr
		st = tr
	}

	track = &st

	rl.StopMusicStream(*track)
	rl.PlayMusicStream(*track)
	rl.SetMusicVolume(*track, musicVolume)
}

// SetMusicVolume sets the music volume
func SetMusicVolume(vol float32) {
	musicVolume = vol

	if track == nil {
		return
	}

	if rl.IsMusicPlaying(*track) {
		rl.SetMusicVolume(*track, musicVolume)
	}
}

// UpdateMusic checks for music playback and updates the tracklist
func UpdateMusic() {
	if track == nil {
		return
	}

	if rl.IsMusicPlaying(*track) {
		rl.UpdateMusicStream(*track)
		return
	}

	LoadNextTrack()
}

// PauseMusic pauses the music stream
func PauseMusic() {
	if track == nil {
		return
	}

	rl.PauseMusicStream(*track)
}
