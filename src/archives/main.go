package main

import (
	"github.com/davecgh/go-spew/spew"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

const (
	screenW = 640
	screenH = 480

	// Apply 2x upscaling
	windowW = screenW * 2
	windowH = screenH * 2

	defaultCameraZoom = 1.79
)

func init() {
	rl.SetCallbackFunc(main)
}

func main() {
	system.InitAssets("default.dta", true)

	rl.InitWindow(800, 450, "raylib [core] example - basic window")

	rl.SetTargetFPS(60)

	avatar := system.FindAsset("gfx/avatar.png")
	spew.Dump(avatar)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
