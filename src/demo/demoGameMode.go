package main

import (
	"encoding/gob"
	"fmt"
	"math"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type demoGameMode struct {
	playState      int
	textWave       int32
	showHelpScreen bool
}

func (g *demoGameMode) Init() {
	core.LoadPlaylist("tracklist.txt")
	core.LoadNextTrack()

	// test class
	err := core.RegisterClass("demo_testclass", "TestClass", NewTestClass)

	// player class
	core.RegisterClass("player", "Player", NewPlayer)

	if err != nil {
		fmt.Printf("Custom type registration has failed: %s", err.Error())
	}

	initLevels()

	g.playState = stateMenu

	if fmapload {
		g.playState = statePlay
		core.LoadMap(playMapName)
		core.InitMap()
	}

	initShaders()
}

func (g *demoGameMode) Shutdown() {}

func (g *demoGameMode) Update() {
	switch g.playState {
	case statePaused:
		if rl.IsKeyPressed(rl.KeyEscape) {
			g.playState = statePlay
		}

		if system.IsKeyPressed("use") {
			core.FlushMaps()
			g.playState = stateLevelSelection
			levelSelection.selectedChoice = 0
			return
		}

	case stateMenu:
		g.textWave = int32(math.Round(math.Sin(float64(rl.GetTime()) * 10)))

		if system.IsKeyPressed("use") {
			g.playState = stateLevelSelection
		}

		if rl.IsKeyPressed(rl.KeyEscape) {
			core.CloseGame()
			return
		}

	case stateLevelSelection:
		g.updateLevelSelection()

		if rl.IsKeyPressed(rl.KeyEscape) {
			g.playState = stateMenu
		}

	case statePlay:
		core.UpdateMaps()
		updateDialogue()
		updateNotifications()

		if rl.IsKeyPressed(rl.KeyEscape) && core.CurrentMap.Name != "start" {
			g.playState = statePaused
		}

		if rl.IsKeyPressed(rl.KeyF1) {
			g.showHelpScreen = !g.showHelpScreen
		}
	}

	updateInternals(g)
}

func (g *demoGameMode) Serialize(enc *gob.Encoder) {
	data := demoGameSaveData{
		ObjectCounter: dynobjCounter,
	}

	enc.Encode(data)
}

func (g *demoGameMode) Deserialize(dec *gob.Decoder) {
	var saveData demoGameSaveData
	dec.Decode(&saveData)

	dynobjCounter = saveData.ObjectCounter
}

type demoGameSaveData struct {
	ObjectCounter int
}

func (g *demoGameMode) Draw() {
	drawBackground()

	rl.BeginMode2D(core.RenderCamera)
	{
		core.DrawMap(true)
	}
	rl.EndMode2D()
}

func (g *demoGameMode) DrawUI() {
	switch g.playState {
	case stateMenu:
		core.DrawTextCentered("Rurik Framework", system.ScreenWidth/2, system.ScreenHeight/2-20+g.textWave, 24, rl.RayWhite)
		core.DrawTextCentered("Press E/ENTER to continue", system.ScreenWidth/2, system.ScreenHeight/2+5+g.textWave, 14, rl.White)

	case statePaused:
		rl.DrawRectangle(0, 0, system.ScreenWidth, system.ScreenHeight, rl.Fade(rl.Black, 0.8))
		core.DrawTextCentered("Rurik Framework", system.ScreenWidth/2, system.ScreenHeight/2-20+g.textWave, 24, rl.RayWhite)
		core.DrawTextCentered("Press ESC to unpause or E/ENTER to return to the menu", system.ScreenWidth/2, system.ScreenHeight/2+5+g.textWave, 14, rl.White)

	case stateLevelSelection:
		core.DrawTextCentered("Rurik Framework", system.ScreenWidth/2, system.ScreenHeight/2-20+g.textWave, 24, rl.RayWhite)
		g.drawLevelSelection()

	case statePlay:
		core.DrawMapUI()
		drawDialogue()
		drawNotifications()

		if core.CurrentMap.Name != "start" {
			var xoffs int32 = 15
			yoffs := system.ScreenHeight - 120
			if g.showHelpScreen {
				rl.DrawText("Press F5 at any time to go back to the menu.", xoffs, yoffs-40, 12, rl.RayWhite)
				rl.DrawText("Press F2 to save and F3 to load a game state.", xoffs, yoffs-52, 12, rl.RayWhite)
				rl.DrawText("Press F9 to spawn a light object.", xoffs, yoffs-64, 12, rl.RayWhite)
				rl.DrawText("Use your mousewheel to zoom in/out.", xoffs, yoffs-76, 12, rl.RayWhite)
			} else {
				rl.DrawText("Press F1 for help.", xoffs-10, system.ScreenHeight-20, 12, rl.RayWhite)
			}
		} else {
			core.DrawTextCentered("Rurik Framework", system.ScreenWidth/2, system.ScreenHeight/3, 24, rl.RayWhite)
		}

		if core.CurrentMap.Name == "village" {
			// draw a minimap
			{
				rl.DrawRectangle(system.ScreenWidth-105, 5, 100, 100, rl.Blue)
				rl.DrawTexturePro(
					minimap.RenderTexture.Texture,
					rl.NewRectangle(0, 0,
						float32(minimap.RenderTexture.Texture.Width),
						float32(-minimap.RenderTexture.Texture.Height)),
					rl.NewRectangle(float32(system.ScreenWidth)-102, 8, 94, 94),
					rl.Vector2{},
					0,
					rl.White,
				)
			}

			// draw shadertoy example
			{
				rl.DrawRectangle(system.ScreenWidth-105, 110, 100, 100, rl.Fade(rl.Red, 0.6))
				rl.DrawTexturePro(
					shadertoy.RenderTexture.Texture,
					rl.NewRectangle(0, 0,
						float32(shadertoy.RenderTexture.Texture.Width),
						float32(shadertoy.RenderTexture.Texture.Height)),
					rl.NewRectangle(float32(system.ScreenWidth)-102, 113, 94, 94),
					rl.Vector2{},
					0,
					rl.White,
				)
			}
		}
	}
}

func (g *demoGameMode) PostDraw() {

	switch g.playState {
	case stateMenu:

	case statePaused:
		fallthrough

	case statePlay:
		// Generates and applies the lightmaps
		core.UpdateLightingSolution()

		if core.CurrentMap.Name == "village" || core.CurrentMap.Name == "sewer" {
			bloom.Apply()
			minimap.Apply()
			shadertoy.Apply()
		}
	}

}
