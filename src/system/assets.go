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

package system

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	goaseprite "github.com/solarlune/GoAseprite"
	rl "github.com/zaklaus/raylib-go/raylib"
)

var (
	textures = make(map[string]rl.Texture2D)
	animData = make(map[string]goaseprite.File)
	fileData = make(map[string][]byte)

	// MapName represents the currently loaded map name
	MapName string
)

// GetTexture retrieves a cached texture from disk
func GetTexture(texturePath string) *rl.Texture2D {
	texturePath = fmt.Sprintf("assets/gfx/%s", texturePath)

	tx, ok := textures[texturePath]

	if ok {
		return &tx
	}

	tx = rl.LoadTexture(texturePath)

	textures[texturePath] = tx

	return &tx
}

// GetAnimData retrieves a cached Aseprite anim data from a disk
func GetAnimData(animPath string) goaseprite.File {
	animPath = fmt.Sprintf("assets/gfx/%s.json", animPath)

	ani, ok := animData[animPath]

	if ok {
		return ani
	}

	ani = goaseprite.Load(animPath)
	animData[animPath] = ani
	return ani
}

// GetFile retrieves a file from a disk
func GetFile(path string, checkGlobalDir bool) []byte {
	if checkGlobalDir {
		if _, err := os.Stat(fmt.Sprintf("assets/%s", path)); !os.IsNotExist(err) {
			return GetRootFile(path)
		}
	}

	return GetRootFile(fmt.Sprintf("map/%s/%s", MapName, path))
}

// GetRootFile retrieves a file inside of game root from a disk
func GetRootFile(path string) []byte {
	path = fmt.Sprintf("assets/%s", path)

	data, ok := fileData[path]

	if ok {
		return data
	}

	newData, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("Could not load file %s\n", path)
		return nil
	}

	fileData[path] = newData
	return newData
}
