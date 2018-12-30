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

package core

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

var (
	// WorldTexture represents the render target used by the game world
	WorldTexture system.RenderTarget

	// UITexture represents the render target used by the interface
	UITexture system.RenderTarget

	// NullTexture should only be used for rendering with shaders
	// that don't accept any previous texture data
	NullTexture system.RenderTarget

	finalRenderTexture system.RenderTarget
	renderTextureQueue = []renderQueueEntry{}

	// RenderCamera is a read-only camera used only for rendering
	RenderCamera rl.Camera2D

	blurTextures []system.RenderTarget
	blurProgram  system.Program
)

type renderQueueEntry struct {
	Target    system.RenderTarget
	FlipY     bool
	blendMode rl.BlendMode
}

// PushRenderTarget appends the render target to the queue to be processed by the compositor pipeline
func PushRenderTarget(tex system.RenderTarget, flipY bool, blendMode rl.BlendMode) {
	renderTextureQueue = append(renderTextureQueue, renderQueueEntry{
		Target:    tex,
		FlipY:     flipY,
		blendMode: blendMode,
	})
}

func renderGame() {
	rl.BeginDrawing()
	{ // Render the game world
		rl.BeginTextureMode(WorldTexture)
		{
			drawProfiler.StartInvocation()
			{
				rl.ClearBackground(rl.Black)

				CurrentGameMode.Draw()
			}
			drawProfiler.StopInvocation()
		}
		rl.EndTextureMode()

		// Render all UI elements
		rl.BeginTextureMode(UITexture)
		{
			rl.ClearBackground(rl.Blank)
			CurrentGameMode.DrawUI()
			DrawEditor()
		}
		rl.EndTextureMode()

		// Render all post-fx elements
		CurrentGameMode.PostDraw()

		// Blend results into one final texture
		rl.BeginTextureMode(finalRenderTexture)
		{
			rl.ClearBackground(rl.Black)
			rl.DrawTexture(WorldTexture.Texture, 0, 0, rl.White)

			// process the render queue
			for _, r := range renderTextureQueue {
				rl.BeginBlendMode(r.blendMode)
				{
					v := r.Target
					height := float32(v.Texture.Height)
					if r.FlipY {
						height *= -1
					}
					rl.DrawTextureRec(
						v.Texture,
						rl.NewRectangle(0, 0, float32(v.Texture.Width), height),
						rl.Vector2{},
						rl.White,
					)
					rl.EndBlendMode()
				}
			}

			rl.BeginBlendMode(rl.BlendAlpha)
			{
				rl.DrawTexture(UITexture.Texture, 0, 0, rl.White)
			}
			rl.EndBlendMode()
		}
		rl.EndTextureMode()
	}
	rl.EndDrawing()

	// output final render texture onto the screen
	rl.DrawTexturePro(
		finalRenderTexture.Texture,
		rl.NewRectangle(0, 0, float32(system.ScreenWidth), float32(system.ScreenHeight)),
		rl.NewRectangle(0, 0, float32(system.WindowWidth), float32(system.WindowHeight)),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)

	renderTextureQueue = []renderQueueEntry{}
}

// BlurRenderTarget blurs the render target
func BlurRenderTarget(tex system.RenderTarget, maxIter int) {
	if blurProgram.Shader.ID == 0 {
		blurProgram = system.NewProgram("", "assets/shaders/blur.fs")
	}

	var hor int32 = 1
	srcTex := tex

	for i := 0; i < maxIter; i++ {
		blurProgram.SetShaderValuei("horizontal", []int32{hor}, 1)
		blurProgram.RenderToTexture(srcTex, blurTextures[hor])
		srcTex = blurTextures[hor]
		hor = 1 - hor
	}

	system.CopyToRenderTarget(srcTex, tex, hor == 1)
}

func updateSystemRenderTargets() {
	if blurTextures != nil {
		rl.UnloadRenderTexture(blurTextures[0])
		rl.UnloadRenderTexture(blurTextures[1])
	}

	blurTextures = []system.RenderTarget{
		system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight),
		system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight),
	}
}
