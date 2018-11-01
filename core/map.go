package core

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"

	"../system"
	rl "github.com/gen2brain/raylib-go/raylib"
	tiled "github.com/lafriks/go-tiled"
)

var (
	tilemap  *tiled.Map
	tilesets map[string]*tilesetData
	mapName  string
)

type tilesetImageData struct {
	Source string `xml:"source,attr"`
	Width  int32  `xml:"width,attr"`
	Height int32  `xml:"height,attr"`
}

type tilesetData struct {
	Version      string           `xml:"version,attr"`
	TiledVersion string           `xml:"tiledversion,attr"`
	Name         string           `xml:"name,attr"`
	TileWidth    int32            `xml:"tilewidth,attr"`
	TileHeight   int32            `xml:"tileheight,attr"`
	TileCount    int32            `xml:"tilecount,attr"`
	Columns      int32            `xml:"columns,attr"`
	ImageInfo    tilesetImageData `xml:"image"`
	Image        rl.Texture2D
}

// LoadMap loads map data
func LoadMap(name string) {
	var err error
	tilesets = make(map[string]*tilesetData)
	tilemap, err = tiled.LoadFromFile(fmt.Sprintf("assets/map/%s/%s.tmx", name, name))

	if err != nil {
		log.Fatalf("Map %s could not be loaded!\n", name)
		return
	}

	mapName = name

	flushObjects()
	CreateObjects()
	postProcessObjects()
}

// ReloadMap reloads map data
// Useful during development due to runtime asset hot-reloading capability.
func ReloadMap() {
	LoadMap(mapName)
}

func loadTilesetData(tilesetName string) *tilesetData {
	val, ok := tilesets[tilesetName]

	if ok {
		return val
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("assets/map/%s/%s", mapName, tilesetName))

	if err != nil {
		log.Fatalf("Tileset data %s could not be loaded!\n", tilesetName)
		return nil
	}

	loadedTileset := &tilesetData{}

	err = xml.Unmarshal(data, loadedTileset)

	if err != nil {
		log.Fatalf("Tileset data %s could not be parsed:\n\t %s!\n", tilesetName, err.Error())
		return nil
	}

	loadedTileset.Image = system.GetTexture(fmt.Sprintf("assets/map/%s/%s", mapName, loadedTileset.ImageInfo.Source))

	tilesets[tilesetName] = loadedTileset
	return loadedTileset
}

// CreateObjects iterates over all object definitions and spawns objects
func CreateObjects() {
	for _, objectGroup := range tilemap.ObjectGroups {
		for _, object := range objectGroup.Objects {
			spawnObject(object)
		}
	}
}

// DrawTilemap renders the loaded map
func DrawTilemap() {
	tileW := float32(tilemap.TileWidth)
	tileH := float32(tilemap.TileHeight)

	for _, layer := range tilemap.Layers {
		if !layer.Visible {
			continue
		}

		for tileIndex, tile := range layer.Tiles {
			id := int(tile.ID)

			if tile.IsNil() {
				continue
			}

			tilesetData := loadTilesetData(tile.Tileset.Source)

			if tilesetData == nil {
				log.Fatalf("Tileset data '%s' points to nil reference!\n", tile.Tileset.Source)
				return
			}

			tilemapImage := tilesetData.Image
			tileRow := int(tilemapImage.Width) / int(tileW)

			tileWorldX, tileWorldY := GetPositionFromID(uint32(tileIndex), tileW, tileH)

			tileX := float32(id%tileRow) * tileW
			tileY := float32(id/tileRow) * tileH

			sourceRect := rl.NewRectangle(tileX, tileY, tileW, tileH)
			var rot float32

			if tile.HorizontalFlip {
				sourceRect.Width *= -1
			}

			if tile.VerticalFlip {
				sourceRect.Height *= -1
			}

			if tile.DiagonalFlip {
				sourceRect.Width *= -1
				rot = 90
			}

			rl.DrawTexturePro(tilemapImage,
				sourceRect,
				rl.NewRectangle(tileWorldX+tileW/2, tileWorldY+tileH/2, tileW, tileH), rl.NewVector2(tileW/2, tileH/2), rot, rl.White)
		}
	}
}

// GetPositionFromID returns XY world position based on tile ID
func GetPositionFromID(index uint32, tileW, tileH float32) (float32, float32) {
	tileWorldX := float32(index%uint32(tilemap.Width)) * tileW
	tileWorldY := float32(index/uint32(tilemap.Width)) * tileH

	return tileWorldX, tileWorldY
}
