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

	dirtifyTexture system.RenderTarget
	dirtifyProgram system.Program
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
					rl.DrawTexturePro(
						v.Texture,
						rl.NewRectangle(0, 0, float32(v.Texture.Width), height),
						rl.NewRectangle(0, 0, float32(system.ScreenWidth), float32(system.ScreenHeight)),
						rl.Vector2{},
						0,
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

	blurProgram.SetShaderValuei("size", []int32{tex.Texture.Width, tex.Texture.Height}, 2)

	for i := 0; i < maxIter; i++ {
		blurProgram.SetShaderValuei("horizontal", []int32{hor}, 1)
		blurProgram.RenderToTexture(srcTex, blurTextures[hor])
		srcTex = blurTextures[hor]
		hor = 1 - hor
	}

	system.CopyToRenderTarget(srcTex, tex, hor == 1)
}

// DirtifyRenderTarget adds noise to the render target
func DirtifyRenderTarget(tex system.RenderTarget) {
	if dirtifyProgram.Shader.ID == 0 {
		dirtifyProgram = system.NewProgramFromCode("", dirtifyProgramSrcCode)
	}

	dirtifyProgram.SetShaderValuei("size", []int32{tex.Texture.Width, tex.Texture.Height}, 2)
	dirtifyProgram.SetShaderValue("time", []float32{rl.GetTime()}, 2)

	dirtifyProgram.RenderToTexture(tex, dirtifyTexture)
	system.CopyToRenderTarget(dirtifyTexture, tex, false)
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

	if dirtifyTexture.ID == 0 {
		rl.UnloadRenderTexture(dirtifyTexture)
	}

	dirtifyTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
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

const dirtifyProgramSrcCode = `
#version 330
#define M_PI 3.14159265358979323846

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform float time;
uniform vec2 size = vec2(640, 480);

// Output fragment color
out vec4 finalColor;

float rand(vec2 c){
	return fract(sin(dot(c.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

float mod289(float x){return x - floor(x * (1.0 / 289.0)) * 289.0;}
vec4 mod289(vec4 x){return x - floor(x * (1.0 / 289.0)) * 289.0;}
vec4 perm(vec4 x){return mod289(((x * 34.0) + 1.0) * x);}

float noise(vec3 p){
    vec3 a = floor(p);
    vec3 d = p - a;
    d = d * d * (3.0 - 2.0 * d);

    vec4 b = a.xxyy + vec4(0.0, 1.0, 0.0, 1.0);
    vec4 k1 = perm(b.xyxy);
    vec4 k2 = perm(k1.xyxy + b.zzww);

    vec4 c = k2 + a.zzzz;
    vec4 k3 = perm(c);
    vec4 k4 = perm(c + 1.0);

    vec4 o1 = fract(k3 * (1.0 / 41.0));
    vec4 o2 = fract(k4 * (1.0 / 41.0));

    vec4 o3 = o2 * d.z + o1 * (1.0 - d.z);
    vec2 o4 = o3.yw * d.x + o3.xz * (1.0 - d.x);

    return o4.y * d.y + o4.x * (1.0 - d.y);
}

void main()
{
	vec3 col = texture2D(texture0, fragTexCoord).rgb;
	vec3 mcol = mix(col, col+vec3(5,5,5), noise(col*0.111830397323219*rand(vec2(-time,2.0*time))));
	vec3 tcol = mix(col, mcol, rand(fragTexCoord));
	finalColor = vec4(tcol, 1.0);
}
`
