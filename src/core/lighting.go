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
)

func updateLightingSolution() {
	if additiveLightTexture == nil || WindowWasResized {
		if additiveLightTexture != nil {
			rl.UnloadRenderTexture(*additiveLightTexture)
			rl.UnloadRenderTexture(*multiplicativeLightTexture)
		}

		additiveLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		multiplicativeLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
	}

	populateAdditiveLayer()
	populateMultiplicativeLight()

	PushRenderTarget(multiplicativeLightTexture, false, rl.BlendMultiplied)
	PushRenderTarget(additiveLightTexture, false, rl.BlendAdditive)
}

func populateAdditiveLayer() {
	objs := CurrentMap.World.Objects

	rl.BeginTextureMode(*additiveLightTexture)
	{
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(RenderCamera)
		{
			for _, o := range objs {
				if !o.HasSpecularLight {
					continue
				}

				col := o.Color
				rl.DrawCircleGradient(
					int32(o.Position.X+16+o.Offset.X),
					int32(o.Position.Y-16+o.Offset.Y),
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
	objs := CurrentMap.World.Objects

	rl.BeginTextureMode(*multiplicativeLightTexture)
	{
		rl.ClearBackground(SkyColor)

		rl.BeginBlendMode(rl.BlendAdditive)
		{
			rl.BeginMode2D(RenderCamera)
			{
				for _, o := range objs {
					if !o.HasLight {
						continue
					}

					rl.DrawCircleGradient(
						int32(o.Position.X+16+o.Offset.X),
						int32(o.Position.Y-16+o.Offset.Y),
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
