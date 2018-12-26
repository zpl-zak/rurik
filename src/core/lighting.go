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
	additiveLightTexture       system.RenderTarget
	multiplicativeLightTexture system.RenderTarget
	blurLightTextures          []system.RenderTarget

	blurProgram system.Program
)

func updateLightingSolution() {
	populateAdditiveLayer()
	blurLight(additiveLightTexture, 32)

	populateMultiplicativeLight()
	blurLight(multiplicativeLightTexture, 32)
	PushRenderTarget(multiplicativeLightTexture, false, rl.BlendMultiplied)

	PushRenderTarget(additiveLightTexture, false, rl.BlendAdditive)
	//system.CopyToRenderTarget(multiplicativeLightTexture, WorldTexture, true)
}

func populateAdditiveLayer() {
	objs := CurrentMap.World.GetObjectsOfType("light", false)

	if additiveLightTexture == nil || WindowWasResized {
		additiveLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		multiplicativeLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		blurLightTextures = []system.RenderTarget{
			system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight),
			system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight),
		}
	}

	rl.BeginTextureMode(*additiveLightTexture)
	{
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(RenderCamera)
		{
			for _, o := range objs {
				col := o.color
				col.A /= 2
				rl.DrawCircleGradient(
					int32(o.Position.X),
					int32(o.Position.Y),
					o.radius,
					col,
					rl.Blank,
				)
			}
		}
		rl.EndMode2D()
	}
	rl.EndTextureMode()
}

func updateShadersOnResize() {
	blurProgram.SetShaderValue("size", []float32{float32(system.ScreenWidth), float32(system.ScreenHeight)}, 2)
}

func blurLight(tex system.RenderTarget, maxIter int) {
	if blurProgram.Shader.ID == 0 {
		blurProgram = system.NewProgram("", "assets/shaders/blur.fs")
		updateShadersOnResize()
	}

	if WindowWasResized {
		updateShadersOnResize()
	}

	var hor int32 = 1
	srcTex := tex

	for i := 0; i < maxIter; i++ {
		blurProgram.SetShaderValuei("horizontal", []int32{hor}, 1)

		blurProgram.RenderToTexture(srcTex, blurLightTextures[hor])
		srcTex = blurLightTextures[hor]
		hor = 1 - hor
	}

	system.CopyToRenderTarget(srcTex, tex, hor == 1)
}

func populateMultiplicativeLight() {
	objs := CurrentMap.World.GetObjectsOfType("light", false)

	if additiveLightTexture == nil || WindowWasResized {
		additiveLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		multiplicativeLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		blurLightTextures = []system.RenderTarget{
			system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight),
			system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight),
		}
	}

	rl.BeginTextureMode(*multiplicativeLightTexture)
	{
		rl.ClearBackground(SkyColor)

		rl.BeginBlendMode(rl.BlendAdditive)
		{
			rl.BeginMode2D(RenderCamera)
			{
				for _, o := range objs {
					rl.DrawCircleGradient(
						int32(o.Position.X),
						int32(o.Position.Y),
						o.attenuation,
						o.color,
						rl.Blank,
					)
				}
			}
			rl.EndMode2D()
		}
		rl.EndBlendMode()
	}
	rl.EndTextureMode()
}
