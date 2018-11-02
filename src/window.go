package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	screenTexture rl.RenderTexture2D

	//ScreenWidth render target width
	ScreenWidth int32

	//ScreenHeight render target height
	ScreenHeight int32

	//WindowWidth window width
	WindowWidth int32

	//WindowHeight window height
	WindowHeight int32

	// ScaleRatio between window and screen resolution
	ScaleRatio int32
)

// InitRenderer initializes the renderer and creates the window
func InitRenderer(title string, winW, winH int32) {
	//rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(winW, winH, title)
	rl.SetTargetFPS(60)

	WindowWidth = winW
	WindowHeight = winH
}

// CreateRenderTarget generates a render target used by the game
func CreateRenderTarget(screenW, screenH int32) *rl.RenderTexture2D {
	screenTexture = rl.LoadRenderTexture(screenW, screenH)
	rl.SetTextureFilter(screenTexture.Texture, rl.FilterPoint)

	ScreenWidth = screenW
	ScreenHeight = screenH
	ScaleRatio = WindowWidth / ScreenWidth

	return &screenTexture
}

// GetRenderTarget retrieves the render target
func GetRenderTarget() *rl.RenderTexture2D {
	return &screenTexture
}
