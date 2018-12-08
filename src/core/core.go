/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:26:53
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-08 21:39:44
 */

package core // madaraszd.net/zaklaus/rurik
import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// MainCamera is the primary camera used for the viewport
	MainCamera *Object

	// LocalPlayer is player's object
	LocalPlayer *Object

	// DebugMode switch
	DebugMode = true

	// TimeScale is game update time scale
	TimeScale = 1

	// CurrentGameMode is our main gameplay rules descriptor
	CurrentGameMode GameMode

	// ScreenTexture represents the render target
	ScreenTexture *rl.RenderTexture2D
)

const (
	// GameVersion describes itself
	GameVersion = "1.0.0"
)

// InitCore initializes the game engine
func InitCore(name string, windowW, windowH, screenW, screenH int32) {
	InitRenderer(name, windowW, windowH)
	ScreenTexture = CreateRenderTarget(screenW, screenH)
	InitInput()
	rl.InitAudioDevice()

	initObjectTypes()
	InitDatabase()
}

// Run executes the main game loop
func Run() {
	lastTime := float64(rl.GetTime())
	var unprocessedTime float64
	var frameCounter float64
	var frames int32

	InitGameProfilers()
	defer shutdown()

	if CurrentGameMode == nil {
		log.Fatalf("No GameMode has been set!\n")
		return
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

			gameModeProfiler.StartInvocation()
			CurrentGameMode.Update()
			gameModeProfiler.StopInvocation()

			shouldRender = true

			unprocessedTime -= float64(FrameTime)
		}

		if shouldRender {
			if DebugMode {
				drawProfiling()
			}

			rl.BeginTextureMode(*ScreenTexture)
			rl.BeginDrawing()
			drawProfiler.StartInvocation()
			{
				rl.ClearBackground(rl.Black)

				CurrentGameMode.Draw()

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

			rl.DrawTexturePro(screenTexture.Texture, rl.NewRectangle(0, 0, float32(ScreenWidth), -float32(ScreenHeight)),
				rl.NewRectangle(0, 0, float32(WindowWidth), float32(WindowHeight)), rl.NewVector2(0, 0), 0, rl.White)

			frames++
		} else {
			//time.Sleep(time.Millisecond)
		}
	}
}

func shutdown() {
	if rl.IsWindowReady() {
		CurrentGameMode.Shutdown()
		rl.CloseWindow()
	}
}

func setupDefaultCamera() {
	defCam := CurrentMap.World.NewObject(nil)

	defCam.Name = "main_camera"
	defCam.Class = "cam"
	defCam.Position = rl.Vector2{}

	defCam.NewCamera()
	defCam.Mode = CameraModeFollow
	defCam.Follow = LocalPlayer
	defCam.Visible = false

	CurrentMap.World.Objects = append(CurrentMap.World.Objects, defCam)
}
