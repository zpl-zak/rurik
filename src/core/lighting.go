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
	aoLightTexture             system.RenderTarget
	origObjectsLightTexture    system.RenderTarget
)

func updateLightingSolution() {
	if additiveLightTexture == nil || WindowWasResized {
		additiveLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		multiplicativeLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		// aoLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		// origObjectsLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
	}

	populateAdditiveLayer()
	//blurLight(additiveLightTexture, 32)

	populateMultiplicativeLight()
	//blurLight(multiplicativeLightTexture, 32)

	//populateAoLight()
	//blurLight(aoLightTexture, 8)

	// Apply lighting layers
	// PushRenderTarget(aoLightTexture, false, rl.BlendMultiplied)
	// PushRenderTarget(origObjectsLightTexture, false, rl.BlendAlpha)
	PushRenderTarget(multiplicativeLightTexture, false, rl.BlendMultiplied)
	PushRenderTarget(additiveLightTexture, false, rl.BlendAdditive)
	//system.CopyToRenderTarget(multiplicativeLightTexture, WorldTexture, true)
}

func populateAdditiveLayer() {
	objs := CurrentMap.World.GetObjectsOfType("light", false)

	rl.BeginTextureMode(*additiveLightTexture)
	{
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(RenderCamera)
		{
			for _, o := range objs {
				col := o.Color
				col.A /= 2
				rl.DrawCircleGradient(
					int32(o.Position.X),
					int32(o.Position.Y),
					float32(o.Radius),
					col,
					rl.Blank,
				)
			}
		}
		rl.EndMode2D()
	}
	rl.EndTextureMode()
}

func populateMultiplicativeLight() {
	objs := CurrentMap.World.GetObjectsOfType("light", false)

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
						o.Attenuation,
						o.Color,
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

func populateAoLight() {
	sky := SkyColor
	SkyColor = rl.NewColor(0, 0, 0, 127)

	rl.BeginTextureMode(*aoLightTexture)
	{
		rl.ClearBackground(rl.White)
		rl.BeginMode2D(RenderCamera)
		CurrentMap.World.DrawObjects()
		rl.EndMode2D()
	}
	rl.EndTextureMode()

	SkyColor = sky

	rl.BeginTextureMode(*origObjectsLightTexture)
	{
		rl.ClearBackground(rl.Blank)
		rl.BeginMode2D(RenderCamera)
		CurrentMap.World.DrawObjects()
		rl.EndMode2D()
	}
	rl.EndTextureMode()
}
