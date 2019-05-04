package main

import (
	"fmt"
	"math"

	rl "github.com/zaklaus/raylib-go/raylib"

	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type level struct {
	title   string
	mapName string
}

var levelSelection struct {
	selectedChoice       int
	levels               []level
	waveTime             int32
	banner               string
	mouseDoublePressTime int32
}

func initLevels() {
	levelSelection.levels = []level{
		level{
			title:   "Scripting & dialogues",
			mapName: "demo",
		},
		level{
			title:   "Lighting & Shaders",
			mapName: "village",
		},
		level{
			title:   "Lighting",
			mapName: "sewer",
		},
		level{
			title:   "Stress test",
			mapName: "stress",
		},
		level{
			title:   "Exit demo",
			mapName: "$exitGame",
		},
	}

	levelSelection.banner = "Welcome to Rurik Framework!\nThis demo showcases the framework's possibilities and features.\nMake a selection, please!"
}

func (g *demoGameMode) drawLevelSelection() {
	levelSelection.waveTime = int32(math.Round(math.Sin(float64(rl.GetTime()) * 40)))

	width := system.ScreenWidth
	start := system.ScreenHeight / 2

	rl.DrawText(levelSelection.banner, 15, 30, 10, rl.RayWhite)

	// choices
	chsX := width / 2
	chsY := start + 40

	rl.DrawRectangle(chsX-120+levelSelection.waveTime, chsY-20, 240+levelSelection.waveTime, int32(len(levelSelection.levels))*15+40, rl.Fade(rl.Black, 0.25))

	if levelSelection.mouseDoublePressTime > 0 {
		levelSelection.mouseDoublePressTime -= int32(1000 * (system.FrameTime * float32(core.TimeScale)))
	} else if levelSelection.mouseDoublePressTime < 0 {
		levelSelection.mouseDoublePressTime = 0
	}

	if len(levelSelection.levels) > 0 {
		for idx, ch := range levelSelection.levels {
			ypos := chsY + int32(idx)*15 - 2
			if idx == levelSelection.selectedChoice {
				rl.DrawRectangle(chsX-100, ypos, 200, 15, rl.DarkPurple)
			}

			core.DrawTextCentered(
				fmt.Sprintf("%s (%s)", ch.title, ch.mapName),
				chsX,
				chsY+int32(idx)*15,
				10,
				rl.White,
			)

			if core.IsMouseInRectangle(chsX-100, ypos, 200, 15) {
				if rl.IsMouseButtonDown(rl.MouseLeftButton) {
					rl.DrawRectangleLines(chsX-100, ypos, 200, 15, rl.Pink)
				} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
					levelSelection.selectedChoice = idx

					if levelSelection.mouseDoublePressTime > 0 {
						g.playLevelSelection()
					} else {
						levelSelection.mouseDoublePressTime = MouseDoublePress
					}
				} else {
					rl.DrawRectangleLines(chsX-100, ypos, 200, 15, rl.Purple)
				}
			}
		}
	}
}

func (g *demoGameMode) updateLevelSelection() {
	if system.IsKeyPressed("down") {
		levelSelection.selectedChoice++

		if levelSelection.selectedChoice >= len(levelSelection.levels) {
			levelSelection.selectedChoice = 0
		}
	}

	if system.IsKeyPressed("up") {
		levelSelection.selectedChoice--

		if levelSelection.selectedChoice < 0 {
			levelSelection.selectedChoice = len(levelSelection.levels) - 1
		}
	}

	if system.IsKeyPressed("use") {
		g.playLevelSelection()
	}
}

func (g *demoGameMode) loadLevel(mapName string) {
	core.FlushMaps()
	core.LoadMap(mapName)
	core.InitMap()
}

func (g *demoGameMode) playLevelSelection() {
	mapName := levelSelection.levels[levelSelection.selectedChoice].mapName

	if mapName == "$exitGame" {
		core.CloseGame()
		return
	}

	g.loadLevel(mapName)
	g.playState = statePlay
}
