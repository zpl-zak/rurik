package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type demoGameMode struct {
	playState int
}

func (g *demoGameMode) Init() {
	core.LoadPlaylist("tracklist.txt")
	core.LoadNextTrack()

	// test class
	err := core.RegisterClass("demo_testclass", "TestClass", NewTestClass)

	if err != nil {
		fmt.Printf("Custom type registration has failed: %s", err.Error())
	}

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
	case stateMenu:

	case statePlay:
		core.UpdateMaps()
	}

	updateInternals(g)
}

func (g *demoGameMode) Serialize() string {
	data := demoGameSaveData{
		ObjectCounter: dynobjCounter,
	}

	ret, _ := jsoniter.MarshalToString(data)
	return ret
}

func (g *demoGameMode) Deserialize(data string) {
	var saveData demoGameSaveData
	jsoniter.UnmarshalFromString(data, &saveData)

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
		rl.DrawText("Press F5 to load map", 185, 15, 16, rl.RayWhite)

	case statePlay:
		core.DrawMapUI()

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

func (g *demoGameMode) PostDraw() {
	// Generates and applies the lightmaps
	core.UpdateLightingSolution()

	bloom.Apply()
	minimap.Apply()
	shadertoy.Apply()
}
