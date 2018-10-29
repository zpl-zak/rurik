package core

import (
	"fmt"

	"../system"
	rl "github.com/gen2brain/raylib-go/raylib"
	tiled "github.com/lafriks/go-tiled"
)

var (
	tilemap      *tiled.Map
	tilemapImage rl.Texture2D
	tilesets     int
	mapName      string
)

// LoadMap loads map data
func LoadMap(name string) {

	// TODO: Fix tileset path
	loadTilemap(fmt.Sprintf("assets/map/%s/%s.tmx", name, name), "assets/gfx/tileset.png")

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

// LoadTilemap loads the data of a map into the memory
func loadTilemap(tilemapPath, tilesetPath string) {
	tilemap, _ = tiled.LoadFromFile(tilemapPath)
	tilemapImage = system.GetTexture(tilesetPath)

	fmt.Println(tilemap.Tilesets)
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

	tileRow := int(tilemapImage.Width) / int(tileW)

	for _, layer := range tilemap.Layers {

		if !layer.Visible {
			continue
		}

		for tileIndex, tile := range layer.Tiles {
			id := int(tile.ID)

			if tile.IsNil() {
				continue
			}

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
