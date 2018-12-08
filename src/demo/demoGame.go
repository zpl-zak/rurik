package main

import (
	"flag"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pkg/profile"
	"madaraszd.net/zaklaus/rurik/src/core"
	"madaraszd.net/zaklaus/rurik/src/system"
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
	gameCamera    rl.Camera2D
)

type demoGameMode struct{}

func (g *demoGameMode) Init() {
	core.LoadPlaylist("tracklist.txt")
	core.LoadNextTrack()

	core.LoadMap(playMapName)
	core.InitMap()

	gameCamera = rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(0, 0), 0, 1)
}

func (g *demoGameMode) Draw() {
	drawBackground()

	rl.BeginMode2D(gameCamera)
	{
		core.DrawMap()
	}
	rl.EndMode2D()

	core.DrawMapUI()
}

func (g *demoGameMode) IgnoreUpdate() bool {
	return false
}

func (g *demoGameMode) Update() {
	if core.DebugMode && rl.IsKeyPressed(rl.KeyF5) {
		core.FlushMaps()
		core.LoadMap("demo")
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
			w.AddObject(newObject)

			dynobjCounter++
		}
	}

	if core.DebugMode {
		wheel := rl.GetMouseWheelMove()
		if wheel != 0 {
			core.MainCamera.SetCameraZoom(core.MainCamera.Zoom + float32(wheel)*0.05)
		}
	}

	gameCamera.Zoom = core.MainCamera.Zoom

	gameCamera.Offset = rl.Vector2{
		X: float32(int(-core.MainCamera.Position.X*core.MainCamera.Zoom + screenW/2)),
		Y: float32(int(-core.MainCamera.Position.Y*core.MainCamera.Zoom + screenH/2)),
	}
}

func drawBackground() {
	bgImage := system.GetTexture("bg.png")

	rows := int(screenW/bgImage.Width) + 1
	cols := int(screenH/bgImage.Height) + 1
	tileW := float32(bgImage.Width)
	tileH := float32(bgImage.Height)
	src := rl.NewRectangle(0, 0, tileW, tileH)

	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			rl.DrawTexturePro(
				bgImage,
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
	playMapName = *flag.String("map", "demo", "Map name to play.")
	enableProfiler := flag.Bool("profile", false, "Enable profiling.")
	flag.Parse()

	if core.DebugMode {
		core.DebugMode = *dbgMode == 1
	}

	core.InitCore("Demo game | Rurik Engine", windowW, windowH, screenW, screenH)

	demoGame := &demoGameMode{}

	//bloom := rl.LoadShader("", "assets/shaders/bloom.fs")
	core.SetMusicVolume(float32(*musicVol) / 100)
	core.WeatherTimeScale = *weatherTimeScale

	if *enableProfiler {
		defer profile.Start(profile.ProfilePath("build")).Stop()
	}

	core.Run(demoGame)
}

func (g *demoGameMode) Shutdown() {

}
