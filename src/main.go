package main

import (
	"flag"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pkg/profile"
)

const (
	screenW = 640
	screenH = 480

	// Apply 2x upscaling
	windowW = screenW * 2
	windowH = screenH * 2
)

var (
	gameCamera rl.Camera2D
)

func main() {
	dbgMode := flag.Int("debug", 1, "Enable/disable debug mode. Works only in debug builds!")
	musicVol := flag.Int("musicvol", 10, "Music volume.")
	weatherTimeScale := flag.Float64("wtimescale", 1, "Weather time scale.")
	playMapName := flag.String("map", "demo", "Map name to play.")
	enableProfiler := flag.Bool("profile", false, "Enable profiling.")
	flag.Parse()

	if DebugMode {
		DebugMode = *dbgMode == 1
	}

	InitRenderer("Sample scene | Rurik Engine", windowW, windowH)
	CreateRenderTarget(screenW, screenH)
	InitInput()
	rl.InitAudioDevice()
	LoadPlaylist("tracklist.txt")
	LoadNextTrack()
	defer shutdown()

	SetMusicVolume(float32(*musicVol) / 100)

	InitCore()
	LoadMap(*playMapName)
	InitMap()

	screenTexture := GetRenderTarget()

	gameCamera = rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(0, 0), 0, 1)

	//bloom := rl.LoadShader("", "assets/shaders/bloom.fs")
	WeatherTimeScale = *weatherTimeScale

	lastTime := float64(rl.GetTime())
	var unprocessedTime float64
	var frameCounter float64
	var frames int32

	initGameProfilers()

	if *enableProfiler {
		defer profile.Start(profile.ProfilePath("build")).Stop()
	}

	for !rl.WindowShouldClose() {
		shouldRender := false
		startTime := float64(rl.GetTime())
		passedTime := startTime - lastTime
		lastTime = startTime
		unprocessedTime += passedTime
		frameCounter += passedTime

		if LocalPlayer == nil {
			log.Fatalln("Local player not defined!")
			return
		}

		if MainCamera == nil {
			setupDefaultCamera()
		}

		if frameCounter > 1 {
			updateProfiling(frameCounter, float64(frames))

			frames = 0
			frameCounter = 0
		}

		for unprocessedTime > float64(FrameTime) {

			UpdateEditor()

			musicProfiler.StartInvocation()
			UpdateMusic()
			musicProfiler.StopInvocation()

			updateProfiler.StartInvocation()
			UpdateMaps()
			updateProfiler.StopInvocation()

			UpdateMapUI()

			customProfiler.StartInvocation()
			updateEssentials()
			customProfiler.StopInvocation()

			shouldRender = true

			unprocessedTime -= float64(FrameTime)
		}

		if shouldRender {
			if DebugMode {
				drawProfiling()
			}

			rl.BeginTextureMode(*screenTexture)
			rl.BeginDrawing()
			drawProfiler.StartInvocation()
			{
				rl.ClearBackground(rl.Black)
				drawBackground()

				rl.BeginMode2D(gameCamera)
				{
					DrawMap()
				}
				rl.EndMode2D()

				DrawMapUI()
				DrawEditor()

				/* rl.BeginShaderMode(bloom)

				rl.DrawTextureRec(
					screenTexture.Texture,
					rl.NewRectangle(0, 0, float32(screenTexture.Texture.Width), float32(-screenTexture.Texture.Height)),
					rl.Vector2{},
					rl.White,
				)

				rl.EndShaderMode() */
			}
			drawProfiler.StopInvocation()
			rl.EndDrawing()

			rl.EndTextureMode()

			rl.DrawTexturePro(screenTexture.Texture, rl.NewRectangle(0, 0, screenW, -screenH),
				rl.NewRectangle(0, 0, float32(windowW), float32(windowH)), rl.NewVector2(0, 0), 0, rl.White)

			frames++
		} else {
			//time.Sleep(time.Millisecond)
		}
	}
}

func shutdown() {
	if rl.IsWindowReady() {
		rl.CloseWindow()
	}
}

func setupDefaultCamera() {
	defCam := CurrentMap.world.NewObject(nil)

	defCam.Name = "main_camera"
	defCam.Class = "cam"
	defCam.Position = rl.Vector2{}

	defCam.NewCamera()
	defCam.Mode = CameraModeFollow
	defCam.Follow = LocalPlayer
	defCam.Visible = false

	CurrentMap.world.Objects = append(CurrentMap.world.Objects, defCam)
}

func drawBackground() {
	bgImage := GetTexture("bg.png")

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

func updateEssentials() {
	if DebugMode && rl.IsKeyPressed(rl.KeyF5) {
		FlushMaps()
		LoadMap("demo")
		InitMap()
		return
	}

	if DebugMode && rl.IsKeyPressed(rl.KeyF7) {
		LoadNextTrack()
	}

	if DebugMode && rl.IsKeyPressed(rl.KeyF2) {
		saveSystem.SaveGame(0, "demo")
	}

	if DebugMode && rl.IsKeyPressed(rl.KeyF3) {
		saveSystem.LoadGame(0)
	}

	if DebugMode {
		wheel := rl.GetMouseWheelMove()
		if wheel != 0 {
			MainCamera.SetCameraZoom(MainCamera.Zoom + float32(wheel)*0.05)
		}
	}
	gameCamera.Zoom = MainCamera.Zoom

	gameCamera.Offset = rl.Vector2{
		X: float32(int(-MainCamera.Position.X*MainCamera.Zoom + screenW/2)),
		Y: float32(int(-MainCamera.Position.Y*MainCamera.Zoom + screenH/2)),
	}
}
