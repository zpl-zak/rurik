/*
   Copyright 2018 V4 Games

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
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	tiled "github.com/zaklaus/go-tiled"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

var (
	// CurrentMap points to currently loaded map
	CurrentMap *Map

	// Maps holds all loaded maps
	Maps map[string]*Map

	mapNodeIsCollapsed      = true
	worldNodeIsCollapsed    = false
	tilesetsNodeIsCollapsed = true
	objectsNodeIsCollapsed  = true
	cullingEnabled          = true
)

const (
	tileHorizontalFlipMask = 0x80000000
	tileVerticalFlipMask   = 0x40000000
	tileDiagonalFlipMask   = 0x20000000
	tileFlip               = tileHorizontalFlipMask | tileVerticalFlipMask | tileDiagonalFlipMask
	tileGIDMask            = 0x0fffffff
)

// Map defines the environment and simulation region (world)
type Map struct {
	tilemap  *tiled.Map
	tilesets map[string]*tilesetData
	mapName  string
	World    *World
	Weather  Weather
}

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
	Image        *rl.Texture2D
	IsCollapsed  bool
}

// LoadMap loads map data
func LoadMap(name string) *Map {
	if Maps == nil {
		Maps = make(map[string]*Map)
	}

	m, ok := Maps[name]

	if ok {
		return m
	}

	cmap := &Map{}
	var err error
	cmap.tilesets = make(map[string]*tilesetData)
	cmap.tilemap, err = tiled.LoadFromFile(fmt.Sprintf("assets/map/%s/%s.tmx", name, name))

	if err != nil {
		log.Fatalf("Map %s could not be loaded: %s\n", name, err.Error())
		return nil
	}

	if cmap.tilemap.Properties == nil {
		cmap.tilemap.Properties = &tiled.Properties{}
	}

	cmap.mapName = name

	world := &World{
		Objects: []*Object{},
	}

	if CurrentMap == nil {
		CurrentMap = cmap
		system.MapName = name
	}

	cmap.CreateObjects(world)
	world.postProcessObjects()

	cmap.Weather = Weather{}
	cmap.Weather.WeatherInit(cmap)

	cmap.World = world

	Maps[name] = cmap

	return cmap
}

// SwitchMap selects the primarily rendered map
func SwitchMap(name string) {
	m, ok := Maps[name]

	if ok {
		CurrentMap = m
		system.MapName = name
	}

	if MainCamera == nil {
		setupDefaultCamera()
	}
}

// FlushMaps disposes all data
func FlushMaps() {
	if CurrentMap != nil {
		CurrentMap.World = nil
	}

	CurrentMap = nil
	Maps = nil
	LocalPlayer = nil
	MainCamera = nil
	initScriptingSystem()
}

// InitMap initializes current map (useful for new game/areas)
func InitMap() {
	if CurrentMap != nil {
		CurrentMap.World.InitObjects()
	} else {
		log.Fatalf("CurrentMap not set, can't initialize the map!\n")
		return
	}
}

// UpdateMaps updates all maps' simulation regions (worlds)
func UpdateMaps() {
	if CurrentMap == nil {
		return
	}

	weatherProfiler.StartInvocation()
	CurrentMap.Weather.UpdateWeather()
	weatherProfiler.StopInvocation()

	for _, m := range Maps {
		m.World.UpdateObjects()
	}
}

// DrawMap draws the tilemap and all renderable objects
func DrawMap(usesCulling bool) {
	cullingEnabled = usesCulling

	if CurrentMap == nil {
		return
	}

	CurrentMap.DrawTilemap(false)
	CurrentMap.World.DrawObjects()
	CurrentMap.DrawTilemap(true) // render all overlays

	CurrentMap.Weather.DrawWeather()
}

// DrawMapUI draw current map's UI elements
func DrawMapUI() {
	if CurrentMap == nil {
		return
	}

	CurrentMap.World.DrawObjectUI()
}

// UpdateMapUI draws debug UI
func UpdateMapUI() {
	if DebugMode {
		if CurrentMap == nil {
			return
		}

		mapNode := pushEditorElement(rootElement, "map", &mapNodeIsCollapsed)

		if !mapNodeIsCollapsed {
			pushEditorElement(mapNode, fmt.Sprintf("name: %s", CurrentMap.mapName), nil)
			pushEditorElement(mapNode, fmt.Sprintf("no. of tilesets: %d", len(CurrentMap.tilesets)), nil)

			tilesetsNode := pushEditorElement(mapNode, "tilesets", &tilesetsNodeIsCollapsed)

			if !tilesetsNodeIsCollapsed {
				i := 0
				for _, v := range CurrentMap.tilesets {
					tilesetNode := pushEditorElement(tilesetsNode, fmt.Sprintf("%d. %s", i, v.Name), &v.IsCollapsed)
					i++

					if !v.IsCollapsed {
						pushEditorElement(tilesetNode, fmt.Sprintf("name: %s", v.Name), nil)
						pushEditorElement(tilesetNode, fmt.Sprintf("image: %s", v.ImageInfo.Source), nil)
						pushEditorElement(tilesetNode, fmt.Sprintf("width: %d", v.ImageInfo.Width), nil)
						pushEditorElement(tilesetNode, fmt.Sprintf("height: %d", v.ImageInfo.Height), nil)
					}
				}
			}

			pushEditorElement(mapNode, fmt.Sprintf("map width: %d", CurrentMap.tilemap.Width), nil)
			pushEditorElement(mapNode, fmt.Sprintf("map height: %d", CurrentMap.tilemap.Height), nil)
			drawWorldUI(mapNode)
		}
	}
}

func drawWorldUI(mapNode *editorElement) {
	worldNode := pushEditorElement(mapNode, "world", &worldNodeIsCollapsed)

	if !worldNodeIsCollapsed {
		pushEditorElement(worldNode, fmt.Sprintf("object count: %d", len(CurrentMap.World.Objects)), nil)
		pushEditorElement(worldNode, fmt.Sprintf("global id cursor: %d", CurrentMap.World.GlobalIndex), nil)

		objsNode := pushEditorElement(worldNode, "objects", &objectsNodeIsCollapsed)

		if !objectsNodeIsCollapsed {
			for i, v := range CurrentMap.World.Objects {
				pushEditorElement(objsNode, fmt.Sprintf("%d. %s (%s)", i, v.Name, v.Class), nil)
			}
		}
	}
}

// ReloadMap reloads map data
// Useful during development due to runtime asset hot-reloading capability.
func ReloadMap(oldMap *Map) *Map {
	oldMap.World.flushObjects()
	return LoadMap(oldMap.mapName)
}

func (m *Map) loadMapTilesetData(tilesetName string) *tilesetData {
	tilesetName = path.Base(tilesetName)

	val, ok := m.tilesets[tilesetName]

	if ok {
		return val
	}

	loadedTileset := loadTilesetData(tilesetName)
	m.tilesets[tilesetName] = loadedTileset
	return loadedTileset
}

func loadTilesetData(tilesetName string) *tilesetData {
	data, err := ioutil.ReadFile(fmt.Sprintf("assets/tilesets/%s", tilesetName))

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

	loadedTileset.Image = system.GetTexture(fmt.Sprintf("../tilesets/%s", loadedTileset.ImageInfo.Source))
	loadedTileset.IsCollapsed = true

	return loadedTileset
}

type objectTemplate struct {
	Tileset tiled.Tileset `xml:"tileset"`
	Object  tiled.Object  `xml:"object"`
}

// CreateObjects iterates over all object definitions and spawns objects
func (m *Map) CreateObjects(w *World) {
	for _, objectGroup := range m.tilemap.ObjectGroups {
		isColGroup := objectGroup.Properties.GetString("col") == "1"
		for _, object := range objectGroup.Objects {
			if isColGroup {
				object.Type = "col"
			}

			var possibleTileset tiled.Tileset

			if object.Template != "" {
				object.Template = path.Join("assets", "templates", path.Base(object.Template))
				fmt.Printf("Creating object using template: %s...\n", object.Template)

				_, err := os.Stat(object.Template)
				if os.IsNotExist(err) {
					log.Fatalf("Could not load template %s!\n", object.Template)
					return
				}

				tplData, _ := ioutil.ReadFile(object.Template)

				var tpl objectTemplate
				xml.Unmarshal(tplData, &tpl)
				tplObject := tpl.Object

				for _, prop := range tplObject.Properties {
					isNew := true
					for _, newProp := range object.Properties {
						if prop.Name == newProp.Name {
							isNew = false
							break
						}
					}

					if isNew {
						object.Properties = append(object.Properties, prop)
					}
				}

				if tplObject.GID > 0 && object.GID == 0 {
					object.GID = tplObject.GID
					possibleTileset = tpl.Tileset
				}

				if object.Width == 0 {
					object.Width = tplObject.Width
				}

				if object.Height == 0 {
					object.Height = tplObject.Height
				}
			}

			obj := w.spawnObject(object)

			if possibleTileset.FirstGID > 0 {
				obj.LocalTileset = loadTilesetData(path.Base(possibleTileset.Source))
			}

			if obj != nil {
				w.AddObject(obj)
			}
		}
	}
}

// DrawTilemap renders the loaded map
func (m *Map) DrawTilemap(renderOverlays bool) {
	tileW := float32(m.tilemap.TileWidth)
	tileH := float32(m.tilemap.TileHeight)

	for _, layer := range m.tilemap.Layers {
		if !layer.Visible {
			continue
		}

		if (layer.Properties.GetString("isOverlay") == "1" && !renderOverlays) ||
			layer.Properties.GetString("isOverlay") != "1" && renderOverlays {
			continue
		}

		for tileIndex, tile := range layer.Tiles {
			id := int(tile.ID)

			if tile.IsNil() {
				continue
			}

			tilesetData := m.loadMapTilesetData(tile.Tileset.Source)

			if tilesetData == nil {
				log.Fatalf("Tileset data '%s' points to nil reference!\n", tile.Tileset.Source)
				return
			}

			tilemapImage := tilesetData.Image

			tileWorldX, tileWorldY := m.GetWorldPositionFromID(uint32(tileIndex), tileW, tileH)

			sourceRect, _ := m.GetTileDataFromID(id)
			var rot float32

			tilePos := rl.NewVector2(tileWorldX+tileW/2, tileWorldY+tileH/2)

			if !isPointWithinFrustum(tilePos) && cullingEnabled {
				continue
			}

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

			rl.DrawTexturePro(*tilemapImage,
				sourceRect,
				rl.NewRectangle(tilePos.X, tilePos.Y, tileW, tileH),
				rl.NewVector2(tileW/2, tileH/2),
				rot,
				SkyColor,
			)
		}
	}
}

// GetWorldPositionFromID returns XY world position based on tile ID
func (m *Map) GetWorldPositionFromID(index uint32, tileW, tileH float32) (float32, float32) {
	tileWorldX := float32(index%uint32(m.tilemap.Width)) * tileW
	tileWorldY := float32(index/uint32(m.tilemap.Width)) * tileH

	return tileWorldX, tileWorldY
}

// GetTileDataFromID retrieves tile source rectangle and source image based on the tile ID
func (m *Map) GetTileDataFromID(tileID int) (rl.Rectangle, *rl.Texture2D) {
	var tilesetID int

	for i, v := range m.tilemap.Tilesets {
		if uint32(tileID) >= v.FirstGID {
			tilesetID = i
		} else {
			break
		}
	}

	return GetFinalTileDataFromID(tileID, m.loadMapTilesetData(path.Base((m.tilemap.Tilesets[tilesetID].Source))))
}

// GetFinalTileDataFromID retrieves the final TileData from a specific tileset
func GetFinalTileDataFromID(tileID int, tileset *tilesetData) (rl.Rectangle, *rl.Texture2D) {
	tileRow := int(tileset.ImageInfo.Width) / int(tileset.TileWidth)

	tileX := float32(tileID%tileRow) * float32(tileset.TileWidth)
	tileY := float32(tileID/tileRow) * float32(tileset.TileHeight)

	return rl.NewRectangle(tileX, tileY, float32(tileset.TileWidth), float32(tileset.TileHeight)), tileset.Image
}
