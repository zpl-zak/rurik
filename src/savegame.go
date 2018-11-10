/*
 * @Author: V4 Games
 * @Date: 2018-11-09 22:54:46
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-11-09 23:25:21
 */

package main

import (
	"io/ioutil"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/json-iterator/go"
)

var (
	saveSystem SaveSystem
)

// SaveSystem manages game save states
type SaveSystem struct {
	Version string      `json:"version"`
	States  []GameState `json:"gameStates"`
}

// GameState describes the serializable save state
type GameState struct {
	SaveName string
	Data     defaultSaveData `json:"data"`
}

// InitSaveSystem initializes the game state system
func (s *SaveSystem) InitSaveSystem() {
	dat, err := ioutil.ReadFile("gamesav.db")
	hasFailed := false

	if err == nil {
		var sav SaveSystem
		err = jsoniter.Unmarshal(dat, &sav)

		if err != nil {
			log.Printf("Gamesav.db is broken, ignoring...\n")
			hasFailed = true
		} else {
			*s = sav
		}
	} else {
		hasFailed = true
	}

	if hasFailed {
		s.States = make([]GameState, 10)
	}

	s.Version = GameVersion
}

// SaveGame saves the game state
func (s *SaveSystem) SaveGame(slotIndex int, stateName string) {
	state := GameState{
		SaveName: stateName,
	}

	state.Data = defaultSaveProvider(&state)

	s.States[slotIndex] = state

	data, _ := jsoniter.Marshal(s)

	ioutil.WriteFile("gamesav.db", data, 0600)
}

// LoadGame restores the game state
func (s *SaveSystem) LoadGame(slotIndex int) {
	state := &s.States[slotIndex]

	defaultLoadProvider(state)
}

// ShutdownSaveSystem closes down the connection
func (s *SaveSystem) ShutdownSaveSystem() {

}

type saveData interface{}

type defaultSaveData struct {
	saveData
	CurrentMap string           `json:"active"`
	Maps       []defaultMapData `json:"maps"`
}

type defaultMapData struct {
	MapName string              `json:"map"`
	Objects []defaultObjectData `json:"objects"`
}

type defaultObjectData struct {
	Name     string
	Position rl.Vector2
	Movement rl.Vector2
	Facing   rl.Vector2
	Data     interface{} `json:"data"`
}

func defaultSaveProvider(state *GameState) defaultSaveData {
	save := defaultSaveData{
		CurrentMap: CurrentMap.mapName,
		Maps:       []defaultMapData{},
	}

	for _, v := range Maps {
		mapData := defaultMapData{
			MapName: v.mapName,
			Objects: []defaultObjectData{},
		}

		for _, b := range v.world.Objects {
			obj := defaultObjectData{
				Name:     b.Name,
				Position: b.Position,
				Movement: b.Movement,
				Facing:   b.Facing,
				Data:     b.Serialize(b),
			}

			mapData.Objects = append(mapData.Objects, obj)
		}

		save.Maps = append(save.Maps, mapData)
	}

	return save
}

func defaultLoadProvider(state *GameState) {
	data := state.Data
	FlushMaps()
	LoadMap(data.CurrentMap)

	for _, mapData := range data.Maps {
		m := LoadMap(mapData.MapName)
		world := mapData.Objects

		for _, wo := range world {
			o, _ := m.world.FindObject(wo.Name)

			o.Position = wo.Position
			o.Movement = wo.Movement
			o.Facing = wo.Facing
		}
	}
}
