package main

import (
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
	system.InitAssets([]string{"gfx.dta"}, true)

	rl.InitWindow(800, 450, "raylib [core] example - basic window")

	rl.SetTargetFPS(60)

	checked := rl.GenImageChecked(32, 32, 1, 1, rl.Black, rl.White)
	tex := rl.LoadTextureFromImage(checked)
	rl.UnloadImage(checked)

	cam := rl.NewCamera3D(
		rl.NewVector3(5, -5, 5),
		rl.NewVector3(0, 0, 0),
		rl.NewVector3(0, 1, 0),
		62,
		0,
	)

	modelsnap := system.CreateRenderTarget(screenW, screenH)

	rl.BeginTextureMode(modelsnap)
	{
		rl.BeginMode3D(cam)
		{
			rl.ClearBackground(rl.Blank)
			rl.DrawCubeTexture(tex, rl.Vector3{}, 4, 4, 4, rl.White)
			rl.DrawGrid(10, 1)
		}
		rl.EndMode3D()
	}
	rl.EndTextureMode()

	avatar := system.GetTexture("gfx/avatar.png")

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexturePro(modelsnap.Texture,
			rl.NewRectangle(0, 0, screenW, screenH),
			rl.NewRectangle(0, 0, screenW, screenH), rl.Vector2{}, 0,
			rl.White)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
		rl.DrawTexturePro(
			*avatar,
			rl.NewRectangle(0, 0, float32(avatar.Width), float32(avatar.Height)),
			rl.NewRectangle(0, 0, 200, 200),
			rl.Vector2{},
			0,
			rl.White,
		)

		rl.EndDrawing()
		rl.PollInputEvents()
	}

	rl.CloseWindow()
}
