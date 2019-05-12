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
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
	goaseprite "github.com/zaklaus/GoAseprite"
	rl "github.com/zaklaus/raylib-go/raylib"
	"gopkg.in/yaml.v2"
)

var (
	textures   = make(map[string]rl.Texture2D)
	animData   = make(map[string]goaseprite.File)
	fileData   = make(map[string][]byte)
	isDBLoaded bool

	// MapName represents the currently loaded map name
	MapName string

	// AssetDatabase contains all the parsed data we use in game
	AssetDatabase []AssetArchive

	// ArchiveEncryptionKey is the key we use to access our assets
	ArchiveEncryptionKey = "letmein123"

	// RegisteredTags contains a list of all tags usable for assets
	RegisteredTags = []string{"none"}
)

const (
	fileTypeGeneric = iota
	fileTypeSprite
	fileTypeAseprite
	fileTypeMap
	fileTypeMusic
)

type annotationTag struct {
	Name  string  `yaml:"name"`
	Value float32 `yaml:"value"`
}

type annotationChunk struct {
	IsHeader    bool            `yaml:"default"`
	FileName    string          `yaml:"file"`
	Name        string          `yaml:"name"`
	Author      string          `yaml:"author"`
	Description string          `yaml:"desc"`
	Type        string          `yaml:"type"`
	Tags        []annotationTag `yaml:"tags"`
	ExtraData   string          `yaml:"extra"`
}

type annotationData struct {
	Name        string `yaml:"name"`
	Author      string `yaml:"author"`
	Version     string `yaml:"version"`
	Description string `yaml:"desc"`

	Chunks []annotationChunk `yaml:"chunks"`
}

// AssetChunk describes an asset file
type AssetChunk struct {
	FileName    string
	Name        string
	Author      string
	Description string
	Type        uint16
	Data        []byte
	Tags        []uint64
	TagValues   []float64
	ExtraData   []byte
}

// AssetArchive describes the game data
type AssetArchive struct {
	Name        string
	Author      string
	Version     string
	Description string

	Chunks []AssetChunk
}

// MatchVector specifies matching tags and their weights for asset lookup
type MatchVector struct {
	Tags    []uint64
	Weights []float64
}

// InitAssets initializes all asset info
func InitAssets(archiveNames []string, isDebugMode bool) {
	for _, v := range archiveNames {
		if isDebugMode {
			tagFileName := fmt.Sprintf("tags/%s.rtag", strings.Split(path.Base(v), ".")[0])

			if _, err := os.Stat(tagFileName); !os.IsNotExist(err) {
				tagData, _ := ioutil.ReadFile(tagFileName)
				tags := parseAnnotationFile(tagData)

				log.Printf(
					"Archive '%s' has been loaded!\n-- Author: %s\n-- Version: %s\n-- Description: %s\n",
					tags.Name,
					tags.Version,
					tags.Author,
					tags.Description,
				)

				spew.Dump(tags)

				buildAssetStorage(v, tags)
			}
		}

		v = fmt.Sprintf("data/%s", v)

		if _, err := os.Stat(v); os.IsNotExist(err) {
			log.Fatalf("Could not load game data from %s!", v)
		}

		var db AssetArchive

		dat, _ := ioutil.ReadFile(v)
		dat = decrypt(dat, ArchiveEncryptionKey)
		ch := bytes.Buffer{}
		ch.Write(dat)
		var dch bytes.Buffer
		r, _ := zlib.NewReader(&ch)
		dch.ReadFrom(r)
		enc := gob.NewDecoder(&dch)
		enc.Decode(&db)

		AssetDatabase = append(AssetDatabase, db)
	}
	isDBLoaded = true
}

func parseAnnotationFile(data []byte) annotationData {
	var tags annotationData
	err := yaml.Unmarshal(data, &tags)

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

// PushTag registers a tag
func PushTag(tag string) uint64 {
	RegisteredTags = append(RegisteredTags, tag)
	return uint64(len(RegisteredTags))
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

			ac.Tags = []uint64{}
			ac.TagValues = []float64{}

			for _, vt := range v.Tags {
				found := false
				for x, vta := range RegisteredTags {
					if vta == vt.Name {
						ac.Tags = append(ac.Tags, uint64(x))
						ac.TagValues = append(ac.TagValues, float64(vt.Value))
						found = true
						break
					}
				}

				if !found {
					ac.Tags = append(ac.Tags, 0)
					ac.TagValues = append(ac.TagValues, 0.0)
				}
			}

			ac.ExtraData = []byte(v.ExtraData)

			a.Chunks = append(a.Chunks, ac)
		}
	}

	ach := bytes.Buffer{}
	enc := gob.NewEncoder(&ach)
	err := enc.Encode(a)

	var fb bytes.Buffer
	w := zlib.NewWriter(&fb)
	w.Write(ach.Bytes())
	w.Close()

	if err != nil {
		log.Fatalf("Error creating asset database! %v", err)
	}

	ioutil.WriteFile(fmt.Sprintf("data/%s", filePath), encrypt(fb.Bytes(), ArchiveEncryptionKey), 0644)
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
	for _, a := range AssetDatabase {
		for _, v := range a.Chunks {
			if fileName == v.FileName {
				return &v
			}
		}
	}

	return nil
}

// GetBestAsset retrieves an asset closest to the MatchVector's description
func GetBestAsset(vec MatchVector) *AssetChunk {
	var res *AssetChunk
	var bestMatch float64

	for _, db := range AssetDatabase {
		for _, v := range db.Chunks {
			var totalMatch float64
			for x, t := range v.Tags {
				a := float64(vec.Tags[t])
				var neg float64 = 1
				if a < 1 {
					neg = -1
				}
				b := v.TagValues[x]
				d0 := math.Abs(float64(a - b))
				d1 := math.Abs((a - 1000000*neg) - b)
				diff := 1 - math.Min(d0, d1)

				w := vec.Weights[t] * diff
				totalMatch += w
			}

			if bestMatch < totalMatch {
				bestMatch = totalMatch
				res = &v
			}
		}
	}

	return res
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

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}
