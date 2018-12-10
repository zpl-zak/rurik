/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:26:53
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 14:26:48
 */

package core // madaraszd.net/zaklaus/rurik
import (
	"log"
	"os"
	"runtime/pprof"

	rl "github.com/gen2brain/raylib-go/raylib"
	"madaraszd.net/zaklaus/rurik/src/system"
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

	// IsRunning tells us whether the game is still running
	IsRunning = true
)

const (
	// GameVersion describes itself
	GameVersion = "1.0.0"

	// DefaultDebugShowAll shows all visible debug elements
	DefaultDebugShowAll = true
)

// InitCore initializes the game engine
func InitCore(name string, windowW, windowH, screenW, screenH int32) {
	system.InitRenderer(name, windowW, windowH)
	ScreenTexture = system.CreateRenderTarget(screenW, screenH)
	system.InitInput()
	rl.InitAudioDevice()

	initDefaultEvents()
	initObjectTypes()
	InitDatabase()
}

// CloseGame exits the game gracefully
func CloseGame() {
	IsRunning = false
}

// Run executes the main game loop
func Run(newGameMode GameMode, enableProfiler bool) {
	CurrentGameMode = newGameMode
	CurrentGameMode.Init()

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

	var cpuProfiler *os.File

	if enableProfiler {
		cpuProfiler, _ = os.Create("build/cpu.pprof")

		pprof.StartCPUProfile(cpuProfiler)
	}

	for IsRunning {
		if rl.WindowShouldClose() {
			IsRunning = false
		}

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

		for unprocessedTime > float64(system.FrameTime) {

			UpdateEditor()

			musicProfiler.StartInvocation()
			UpdateMusic()
			musicProfiler.StopInvocation()

			if !CurrentGameMode.IgnoreUpdate() {
				updateProfiler.StartInvocation()
				UpdateMaps()
				updateProfiler.StopInvocation()
			}

			UpdateMapUI()

			gameModeProfiler.StartInvocation()
			CurrentGameMode.Update()
			gameModeProfiler.StopInvocation()

			shouldRender = true

			unprocessedTime -= float64(system.FrameTime)
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
			}
			drawProfiler.StopInvocation()
			rl.EndDrawing()

			rl.EndTextureMode()

			rl.DrawTexturePro(ScreenTexture.Texture, rl.NewRectangle(0, 0, float32(system.ScreenWidth), -float32(system.ScreenHeight)),
				rl.NewRectangle(0, 0, float32(system.WindowWidth), float32(system.WindowHeight)), rl.NewVector2(0, 0), 0, rl.White)

			frames++
		} else {
			//time.Sleep(time.Millisecond)
		}
	}

	if enableProfiler {
		pprof.StopCPUProfile()
		cpuProfiler.Close()
	}
}

func shutdown() {
	log.Println("Shutting down the engine...")
	CurrentGameMode.Shutdown()
	rl.CloseWindow()
	rl.CloseAudioDevice()
	os.Exit(0)
}

func setupDefaultCamera() {
	defCam := CurrentMap.World.NewObjectPro("main_camera", "cam")
	defCam.Position = LocalPlayer.Position
	defCam.Mode = CameraModeFollow
	defCam.Follow = LocalPlayer
	defCam.Visible = false

	CurrentMap.World.AddObject(defCam)
}
