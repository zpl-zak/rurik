/*
 * @Author: V4 Games
 * @Date: 2018-11-09 02:37:03
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-11 02:46:34
 */

package system

import (
	rl "github.com/zaklaus/raylib-go/raylib"
)

var (
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

	// FrameTime is the target update time
	FrameTime float32 = 1 / 60.0
)

// RenderTarget describes our render texture
type RenderTarget *rl.RenderTexture2D

// InitRenderer initializes the renderer and creates the window
func InitRenderer(title string, winW, winH int32) {
	rl.SetTraceLog(rl.LogError)
	//rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(winW, winH, title)
	rl.SetTargetFPS(0) // we use our own game loop timing.

	WindowWidth = winW
	WindowHeight = winH
}

// CreateRenderTarget generates a render target used by the game
func CreateRenderTarget(screenW, screenH int32) *rl.RenderTexture2D {
	screenTexture := rl.LoadRenderTexture(screenW, screenH)
	rl.SetTextureFilter(screenTexture.Texture, rl.FilterPoint)

	return &screenTexture
}

// CopyToRenderTarget copies one target to another
func CopyToRenderTarget(source, dest *rl.RenderTexture2D) {
	rl.BeginTextureMode(*dest)
	rl.DrawTexturePro(
		source.Texture,
		rl.NewRectangle(0, 0, float32(source.Texture.Width), float32(source.Texture.Height)),
		rl.NewRectangle(0, 0, float32(dest.Texture.Width), float32(dest.Texture.Height)),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)
	rl.EndTextureMode()
}
