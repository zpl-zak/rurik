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
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	jsoniter "github.com/json-iterator/go"

	goaseprite "github.com/zaklaus/GoAseprite"
	rl "github.com/zaklaus/raylib-go/raylib"
)

var (
	textures   = make(map[string]rl.Texture2D)
	animData   = make(map[string]goaseprite.File)
	fileData   = make(map[string][]byte)
	isDBLoaded bool

	// MapName represents the currently loaded map name
	MapName string

	// AssetDatabase contains all the parsed data we use in game
	AssetDatabase AssetArchive
)

const (
	fileTypeGeneric = iota
	fileTypeSprite
	fileTypeAseprite
	fileTypeMap
	fileTypeMusic
)

type annotationChunk struct {
	IsHeader    bool   `json:"$head"`
	FileName    string `json:"file"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"desc"`
	Type        string `json:"type"`
}

type annotationData struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Description string `json:"desc"`

	Chunks []annotationChunk `json:"chunks"`
}

// AssetChunk describes an asset file
type AssetChunk struct {
	FileName    string
	Name        string
	Author      string
	Description string
	Type        uint16
	Data        []byte
}

// AssetArchive describes the game data
type AssetArchive struct {
	Name        string
	Author      string
	Version     string
	Description string

	Chunks []AssetChunk
}

// InitAssets initializes all asset info
func InitAssets(archiveName string, isDebugMode bool) {
	gob.Register(goaseprite.File{})
	gob.Register(AssetChunk{})
	gob.Register(AssetArchive{})
	gob.Register(rl.Texture2D{})

	if isDebugMode {
		tagFileName := fmt.Sprintf("tags/%s.rtag", strings.Split(path.Base(archiveName), ".")[0])

		if _, err := os.Stat(tagFileName); !os.IsNotExist(err) {
			tagData, _ := ioutil.ReadFile(tagFileName)
			tags := parseAnnotationFile(tagData)
			buildAssetStorage(archiveName, tags)
		}
	}

	archiveName = fmt.Sprintf("data/%s", archiveName)

	if _, err := os.Stat(archiveName); os.IsNotExist(err) {
		log.Fatalf("Could not load game data from %s!", archiveName)
	}

	dat, _ := ioutil.ReadFile(archiveName)
	ch := new(bytes.Buffer)
	ch.Write(dat)
	enc := gob.NewDecoder(ch)
	enc.Decode(&AssetDatabase)

	isDBLoaded = true
}

func parseAnnotationFile(data []byte) annotationData {
	var tags annotationData
	err := jsoniter.Unmarshal(data, &tags)

	if err != nil {
		log.Fatalf("Can't parse tag file! Syntax error detected.")
	}

	return tags
}

func mapFileTypeStringToID(class string) uint16 {
	switch class {
	case "gfx":
		fallthrough
	case "sprite":
		{
			return fileTypeSprite
		}

	case "anim":
		{
			return fileTypeAseprite
		}

	case "map":
		{
			return fileTypeMap
		}

	case "music":
		{
			return fileTypeMusic
		}

	default:
		{
			return fileTypeGeneric
		}
	}
}

func buildAssetStorage(filePath string, an annotationData) {
	var a AssetArchive
	a.Name = an.Name
	a.Author = an.Author
	a.Description = an.Description
	a.Version = an.Version
	a.Chunks = []AssetChunk{}

	lastName := a.Name
	lastAuthor := a.Author
	lastDescription := a.Description

	for _, v := range an.Chunks {
		if v.IsHeader {
			setPropertyIfSet(&lastName, v.Name, lastName)
			setPropertyIfSet(&lastAuthor, v.Author, lastAuthor)
			setPropertyIfSet(&lastDescription, v.Description, lastDescription)
		} else {
			var ac AssetChunk

			setPropertyIfSet(&ac.Name, v.Name, lastName)
			setPropertyIfSet(&ac.Author, v.Author, lastAuthor)
			setPropertyIfSet(&ac.Description, v.Description, lastDescription)
			ac.FileName = v.FileName

			ac.Type = mapFileTypeStringToID(v.Type)

			var err error
			ac.Data, err = ioutil.ReadFile(fmt.Sprintf("assets/%s", v.FileName))

			if err != nil {
				log.Fatalf("File %s could not be loaded!\n", v.FileName)
			}

			a.Chunks = append(a.Chunks, ac)
		}
	}

	ach := new(bytes.Buffer)
	enc := gob.NewEncoder(ach)
	err := enc.Encode(a)

	if err != nil {
		log.Fatalf("Error creating asset database! %v", err)
	}

	ioutil.WriteFile(fmt.Sprintf("data/%s", filePath), ach.Bytes(), 0644)
}

func setPropertyIfSet(src *string, value, fallback string) {
	if value == "" {
		*src = fallback
	} else {
		*src = value
	}
}

// FindAsset looks for asset by filename
func FindAsset(fileName string) *AssetChunk {
	for _, v := range AssetDatabase.Chunks {
		if fileName == v.FileName {
			return &v
		}
	}

	return nil
}

// GetTexture retrieves a cached texture from disk
func GetTexture(texturePath string) *rl.Texture2D {
	tx, ok := textures[texturePath]

	if ok {
		return &tx
	}

	a := FindAsset(texturePath)

	if a == nil {
		log.Fatalf("Could not open texture: %s!\n", texturePath)
		return nil
	}

	txImage := rl.LoadImageFromMemory(string(a.Data))
	tx = rl.LoadTextureFromImage(txImage)
	rl.UnloadImage(txImage)

	textures[texturePath] = tx

	return &tx
}

// GetAnimData retrieves a cached Aseprite anim data from a disk
func GetAnimData(animPath string) goaseprite.File {
	ani, ok := animData[animPath]

	if ok {
		return ani
	}

	a := FindAsset(animPath)

	if a == nil {
		log.Fatalf("Could not open aseprite file: %s!\n", animPath)
		return goaseprite.File{}
	}

	dat := goaseprite.Load(string(a.Data))

	if dat == nil {
		// TODO: err
	}

	ani = *dat
	animData[animPath] = ani
	return ani
}

// GetFile retrieves a file from a disk
func GetFile(path string, checkGlobalDir bool) []byte {
	if checkGlobalDir {
		if asset := FindAsset(path); asset != nil {
			return asset.Data
		}
	}

	return GetRootFile(fmt.Sprintf("map/%s/%s", MapName, path))
}

// GetRootFile retrieves a file inside of game root from a disk
func GetRootFile(path string) []byte {
	data, ok := fileData[path]

	if ok {
		return data
	}

	a := FindAsset(path)

	if a == nil {
		log.Fatalf("Could not open file: %s!\n", path)
		return nil
	}

	newData := a.Data
	fileData[path] = newData
	return newData
}
