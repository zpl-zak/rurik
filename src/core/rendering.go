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
		blurProgram = system.NewProgramFromCode("", blurProgramSrcCode)
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

/* Built-in shaders */

const blurProgramSrcCode = `
#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform float time;

// Output fragment color
out vec4 finalColor;

uniform bool horizontal;
uniform vec2 size = vec2(640, 480);

uniform float weight[3] = float[](0.2270270270, 0.3162162162, 0.0702702703);

void main()
{
    // Texel color fetching from texture sampler
    vec3 texelColor = texture(texture0, fragTexCoord).rgb*weight[0];
    vec2 texelOffset = 1.0 / textureSize(texture0, 0);

    if (horizontal) {
        for (int i = 1; i < 5; ++i) {
            texelColor += texture(texture0, fragTexCoord + vec2(texelOffset.x * i, 0.0)).rgb * weight[i];
            texelColor += texture(texture0, fragTexCoord - vec2(texelOffset.x * i, 0.0)).rgb * weight[i];
        }
    } else {
        for (int i = 1; i < 5; ++i) {
            texelColor += texture(texture0, fragTexCoord + vec2(0.0, texelOffset.x * i)).rgb * weight[i];
            texelColor += texture(texture0, fragTexCoord - vec2(0.0, texelOffset.x * i)).rgb * weight[i];
        }
    }

    finalColor = vec4(texelColor, 1.0);
}
`
