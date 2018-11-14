/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:26:21
 * @Last Modified by:   Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-11-14 02:26:21
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/GoAseprite"
)

var (
	textures  = make(map[string]rl.Texture2D)
	scripts   = make(map[string]string)
	dialogues = make(map[string]Dialogue)
	animData  = make(map[string]goaseprite.File)
)

// GetTexture retrieves a cached texture from disk
func GetTexture(texturePath string) rl.Texture2D {
	texturePath = fmt.Sprintf("assets/gfx/%s", texturePath)

	tx, ok := textures[texturePath]

	if ok {
		return tx
	}

	tx = rl.LoadTexture(texturePath)

	textures[texturePath] = tx

	return tx
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

// GetScript retrieves a cached script from a disk
func GetScript(scriptPath string) string {
	scriptPath = fmt.Sprintf("assets/map/%s/scripts/%s", CurrentMap.mapName, scriptPath)

	scr, ok := scripts[scriptPath]

	if ok {
		return scr
	}

	data, err := ioutil.ReadFile(scriptPath)

	if err != nil {
		log.Fatalf("Could not load script: %s\n", scriptPath)
		return ""
	}

	scr = string(data)

	scripts[scriptPath] = scr
	return scr
}

// GetDialogue retrieves a cached dialogue from a disk
func GetDialogue(dialoguePath string) Dialogue {
	dialoguePath = fmt.Sprintf("assets/map/%s/texts/%s", CurrentMap.mapName, dialoguePath)

	dia, ok := dialogues[dialoguePath]

	if ok {
		return dia
	}

	data, err := ioutil.ReadFile(dialoguePath)

	if err != nil {
		log.Fatalf("Could not load texts for %s\n", dialoguePath)
	}

	err = json.Unmarshal(data, &dia)

	if err != nil {
		log.Fatalf("Error loading text %s\n", dialoguePath)
		return Dialogue{}
	}

	dialogues[dialoguePath] = dia
	return dia
}
