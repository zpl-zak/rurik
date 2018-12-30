/*
   Copyright 2018 V4 Games

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package core // github.com/zaklaus/rurik
import (
	"log"
	"os"
	"runtime/pprof"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
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

	// IsRunning tells us whether the game is still running
	IsRunning = true

	// WindowWasResized states whether window was resized at this cycle
	WindowWasResized = false
)

const (
	// GameVersion describes itself
	GameVersion = "1.0.0"

	// DefaultDebugShowAll shows all visible debug elements
	DefaultDebugShowAll = false
)

// InitCore initializes the game engine
func InitCore(name string, windowW, windowH, screenW, screenH int32) {
	system.InitRenderer(name, windowW, windowH)
	system.ScreenWidth = screenW
	system.ScreenHeight = screenH
	system.ScaleRatio = system.WindowWidth / system.ScreenWidth
	updateSystemRenderTargets()

	WorldTexture = system.CreateRenderTarget(screenW, screenH)
	UITexture = system.CreateRenderTarget(screenW, screenH)
	NullTexture = system.CreateRenderTarget(screenW, screenH)
	finalRenderTexture = system.CreateRenderTarget(screenW, screenH)
	system.InitInput()
	rl.InitAudioDevice()

	InitGameProfilers()
	initScriptingSystem()
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

	if DebugMode {
		rl.SetTraceLog(rl.LogError | rl.LogWarning | rl.LogInfo)
	}

	RenderCamera = rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(0, 0), 0, 1)

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

		if MainCamera == nil || (MainCamera != nil && MainCamera.Name == "TempCamera__") {
			setupDefaultCamera()
		}

		if frameCounter > 1 {
			updateProfiling(frameCounter, float64(frames))

			frames = 0
			frameCounter = 0
		}
		updateWindow()

		for unprocessedTime > float64(system.FrameTime) {
			UpdateEditor()

			updateProfiler.StartInvocation()
			musicProfiler.StartInvocation()
			UpdateMusic()
			musicProfiler.StopInvocation()

			gameModeProfiler.StartInvocation()
			CurrentGameMode.Update()
			gameModeProfiler.StopInvocation()

			UpdateMapUI()

			FireEvent("onUpdate")
			updateProfiler.StopInvocation()

			shouldRender = true

			RenderCamera.Offset = rl.Vector2{
				X: float32(int(-MainCamera.Position.X*MainCamera.Zoom + float32(system.ScreenWidth)/2)),
				Y: float32(int(-MainCamera.Position.Y*MainCamera.Zoom + float32(system.ScreenHeight)/2)),
			}

			unprocessedTime -= float64(system.FrameTime)
		}

		if shouldRender {
			if DebugMode {
				drawProfiling()
			}

			renderGame()

			// WindowWasResized should be reset this cycle now
			if WindowWasResized {
				WindowWasResized = false
			}

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
	if CurrentMap == nil {
		MainCamera = &Object{Name: "TempCamera__"}
		return
	}

	defCam := CurrentMap.World.NewObjectPro("main_camera", "cam")
	defCam.Position = LocalPlayer.Position
	defCam.Mode = CameraModeFollow
	defCam.Follow = LocalPlayer
	defCam.Visible = false
	defCam.IsPersistent = false

	CurrentMap.World.AddObject(defCam)
}

func updateWindow() {
	width := int32(rl.GetScreenWidth())
	height := int32(rl.GetScreenHeight())

	if width != system.WindowWidth || height != system.WindowHeight {
		// Re-create all render textures and let user know about the change

		system.WindowWidth = width
		system.WindowHeight = height
		system.ScreenWidth = width / system.ScaleRatio
		system.ScreenHeight = height / system.ScaleRatio
		WindowWasResized = true

		rl.UnloadRenderTexture(WorldTexture)
		WorldTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)

		rl.UnloadRenderTexture(UITexture)
		UITexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)

		rl.UnloadRenderTexture(finalRenderTexture)
		finalRenderTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)

		rl.UnloadRenderTexture(NullTexture)
		NullTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)

		updateSystemRenderTargets()
	}
}
