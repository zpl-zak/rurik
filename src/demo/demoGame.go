package main

import (
	"flag"
	"fmt"
	"path"
	"strings"

	jsoniter "github.com/json-iterator/go"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

const (
	screenW = 640
	screenH = 480

	// Apply 2x upscaling
	windowW = screenW * 2
	windowH = screenH * 2
)

var (
	dynobjCounter int
	playMapName   string
	bloom         *bloomProg
	shadertoy     *shadertoyProg
	minimap       *minimapProg
	pulseManager  *core.Object
)

type demoGameMode struct{}

func (g *demoGameMode) Init() {
	core.LoadPlaylist("tracklist.txt")
	core.LoadNextTrack()

	// test class
	err := core.RegisterClass("demo_testclass", "TestClass", NewTestClass)

	if err != nil {
		fmt.Printf("Custom type registration has failed: %s", err.Error())
	}

	core.LoadMap(playMapName)
	core.InitMap()

	initShaders()
}

func initShaders() {
	bloom = newBloom()
	shadertoy = newShadertoy()
	minimap = newMinimap()
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

func (g *demoGameMode) PostDraw() {
	// Generates and applies the lightmaps
	core.UpdateLightingSolution()

	bloom.Apply()
	minimap.Apply()
	shadertoy.Apply()
}

func (g *demoGameMode) IgnoreUpdate() bool {
	return false
}

func (g *demoGameMode) Update() {
	if pulseManager == nil {
		w := core.CurrentMap.World
		pulseManager = w.NewObjectPro("pulse_script", "script")
		pulseManager.FileName = "pulsatingLights.js"
		w.FinalizeObject(pulseManager)
		//pulseManager.Trigger(pulseManager, nil)
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF5) {
		core.FlushMaps()
		core.LoadMap(playMapName)
		core.InitMap()
		return
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF7) {
		core.LoadNextTrack()
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF2) {
		core.CurrentSaveSystem.SaveGame(0, "demo")
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF3) {
		core.CurrentSaveSystem.LoadGame(0)
	}

	if core.DebugMode && rl.IsKeyPressed(rl.KeyF9) {
		if core.CurrentMap != nil {
			w := core.CurrentMap.World
			name := fmt.Sprintf("dynobj:%d", dynobjCounter)
			newObject := w.NewObjectPro(name, "target")
			newObject.Position = core.LocalPlayer.Position
			newObject.Color = rl.Pink
			newObject.HasLight = true
			newObject.HasSpecularLight = true
			newObject.Attenuation = 256
			newObject.Radius = 8
			w.AddObject(newObject)

			dynobjCounter++
		}
	}

	if core.DebugMode {
		wheel := rl.GetMouseWheelMove()
		if wheel != 0 {
			core.MainCamera.SetCameraZoom(core.MainCamera.Zoom + float32(wheel)*0.12)
		}
	}

	core.RenderCamera.Zoom = core.MainCamera.Zoom

	if core.WindowWasResized {
		initShaders()
	}
}

func drawBackground() {
	bgImage := system.GetTexture("bg.png")

	rows := int(system.ScreenWidth/bgImage.Width) + 1
	cols := int(system.ScreenHeight/bgImage.Height) + 1
	tileW := float32(bgImage.Width)
	tileH := float32(bgImage.Height)
	src := rl.NewRectangle(0, 0, tileW, tileH)

	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			rl.DrawTexturePro(
				*bgImage,
				src,
				rl.NewRectangle(float32(j)*tileW, float32(i)*tileH, tileW, tileH),
				rl.Vector2{},
				0,
				rl.White,
			)
		}
	}
}

func main() {
	dbgMode := flag.Int("debug", 1, "Enable/disable debug mode. Works only in debug builds!")
	musicVol := flag.Int("musicvol", 10, "Music volume.")
	weatherTimeScale := flag.Float64("wtimescale", 1, "Weather time scale.")
	mapName := flag.String("map", "village", "Map name to play.")
	enableProfiler := flag.Bool("profile", false, "Enable profiling.")
	flag.Parse()

	playMapName = *mapName
	playMapName = path.Base(playMapName)
	playMapName = strings.Split(playMapName, ".")[0]

	if core.DebugMode {
		core.DebugMode = *dbgMode == 1
	}

	core.InitCore("Demo game | Rurik Engine", windowW, windowH, screenW, screenH)

	demoGame := &demoGameMode{}

	core.SetMusicVolume(float32(*musicVol) / 100)
	core.WeatherTimeScale = *weatherTimeScale

	core.Run(demoGame, *enableProfiler)
}

func (g *demoGameMode) Shutdown() {

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
