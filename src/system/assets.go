/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:26:21
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 03:38:10
 */

package system

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/GoAseprite"
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
