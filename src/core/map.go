/*
 * @Author: V4 Games
 * @Date: 2018-11-08 16:05:27
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-08 20:47:00
 */

package core

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	tiled "github.com/lafriks/go-tiled"
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
	Image        rl.Texture2D
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

	if cmap.tilemap.Properties == nil {
		cmap.tilemap.Properties = &tiled.Properties{}
	}

	if err != nil {
		log.Fatalf("Map %s could not be loaded!\n", name)
		return nil
	}

	cmap.mapName = name

	world := &World{
		Objects: []*Object{},
	}

	if CurrentMap == nil {
		CurrentMap = cmap
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
	}
}

// FlushMaps disposes all data
func FlushMaps() {
	CurrentMap.World = nil
	CurrentMap = nil
	Maps = nil
	LocalPlayer = nil
	MainCamera = nil
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
	weatherProfiler.StartInvocation()
	CurrentMap.Weather.UpdateWeather()
	weatherProfiler.StopInvocation()

	for _, m := range Maps {
		m.World.UpdateObjects()
	}
}

// DrawMap draws the tilemap and all renderable objects
func DrawMap() {
	CurrentMap.DrawTilemap(false)
	CurrentMap.World.DrawObjects()
	CurrentMap.DrawTilemap(true) // render all overlays

	CurrentMap.Weather.DrawWeather()
}

// DrawMapUI draw current map's UI elements
func DrawMapUI() {
	CurrentMap.World.DrawObjectUI()
}

// UpdateMapUI draws debug UI
func UpdateMapUI() {
	if DebugMode {
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

func (m *Map) loadTilesetData(tilesetName string) *tilesetData {
	val, ok := m.tilesets[tilesetName]

	if ok {
		return val
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("assets/map/%s/%s", m.mapName, tilesetName))

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

	loadedTileset.Image = GetTexture(fmt.Sprintf("../map/%s/%s", m.mapName, loadedTileset.ImageInfo.Source))
	loadedTileset.IsCollapsed = true

	m.tilesets[tilesetName] = loadedTileset
	return loadedTileset
}

// CreateObjects iterates over all object definitions and spawns objects
func (m *Map) CreateObjects(w *World) {
	for _, objectGroup := range m.tilemap.ObjectGroups {
		for _, object := range objectGroup.Objects {
			obj := w.spawnObject(object)

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

			tilesetData := m.loadTilesetData(tile.Tileset.Source)

			if tilesetData == nil {
				log.Fatalf("Tileset data '%s' points to nil reference!\n", tile.Tileset.Source)
				return
			}

			tilemapImage := tilesetData.Image
			tileRow := int(tilemapImage.Width) / int(tileW)

			tileWorldX, tileWorldY := m.GetPositionFromID(uint32(tileIndex), tileW, tileH)

			tileX := float32(id%tileRow) * tileW
			tileY := float32(id/tileRow) * tileH

			sourceRect := rl.NewRectangle(tileX, tileY, tileW, tileH)
			var rot float32

			tilePos := rl.NewVector2(tileWorldX+tileW/2, tileWorldY+tileH/2)

			if !isPointWithinFrustum(tilePos) {
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

			rl.DrawTexturePro(tilemapImage,
				sourceRect,
				rl.NewRectangle(tilePos.X, tilePos.Y, tileW, tileH), rl.NewVector2(tileW/2, tileH/2), rot, SkyColor)
		}
	}
}

// GetPositionFromID returns XY world position based on tile ID
func (m *Map) GetPositionFromID(index uint32, tileW, tileH float32) (float32, float32) {
	tileWorldX := float32(index%uint32(m.tilemap.Width)) * tileW
	tileWorldY := float32(index/uint32(m.tilemap.Width)) * tileH

	return tileWorldX, tileWorldY
}