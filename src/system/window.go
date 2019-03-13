/*
   Copyright 2019 Dominik Madarász <zaklaus@madaraszd.net>

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
type RenderTarget = rl.RenderTexture2D

// InitRenderer initializes the renderer and creates the window
func InitRenderer(title string, winW, winH int32) {
	rl.SetTraceLog(rl.LogError)
	rl.InitWindow(winW, winH, title)
	rl.SetTargetFPS(0) // we use our own game loop timing.

	WindowWidth = winW
	WindowHeight = winH
}

// CreateRenderTarget generates a render target used by the game
func CreateRenderTarget(screenW, screenH int32) rl.RenderTexture2D {
	screenTexture := rl.LoadRenderTexture(screenW, screenH)
	rl.SetTextureFilter(screenTexture.Texture, rl.FilterPoint)

	return screenTexture
}

// CopyToRenderTarget copies one target to another
func CopyToRenderTarget(source, dest RenderTarget, flipY bool) {
	var neg float32 = 1.0
	if flipY {
		neg = -1.0
	}

	rl.BeginTextureMode(dest)
	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(
		source.Texture,
		rl.NewRectangle(0, 0, float32(source.Texture.Width), float32(source.Texture.Height)*neg),
		rl.NewRectangle(0, 0, float32(dest.Texture.Width), float32(dest.Texture.Height)),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)
	rl.EndTextureMode()
}
