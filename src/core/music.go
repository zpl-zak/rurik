/*
   Copyright 2019 V4 Games

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
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"

	rl "github.com/zaklaus/raylib-go/raylib"
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
		fln := fmt.Sprintf("assets/music/%s", trackName)

		if _, err := os.Stat(fln); os.IsNotExist(err) {
			log.Printf("Music track %s not found!\n", trackName)
			return
		}
		tr := rl.LoadMusicStream(fmt.Sprintf("assets/music/%s", trackName))
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
