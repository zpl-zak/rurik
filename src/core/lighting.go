/*
   Copyright 2019 V4 Games

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

// UpdateLightingSolution generates all lightmaps and pushes them onto the blending stack
// This is an example implementation of 2D lighting inside of an engine, however due to the
// nature of this code, the whole lighting solution can be considered as an add-on, an optional
// feature you can easily enable/disable at your leisure. You can use this code directly, or create
// a different variation of it / develop your own solution. The principle behind this design is to
// show the modularity and versatility of the engine, giving you the power to tweak the game
// any means possible.
func UpdateLightingSolution() {
	if CurrentMap == nil {
		return
	}

	if additiveLightTexture.ID == 0 || WindowWasResized {
		if additiveLightTexture.ID != 0 {
			rl.UnloadRenderTexture(additiveLightTexture)
			rl.UnloadRenderTexture(multiplicativeLightTexture)
		}

		additiveLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		multiplicativeLightTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
	}

	lightingProfiler.StartInvocation()
	populateAdditiveLayer()
	populateMultiplicativeLight()
	lightingProfiler.StopInvocation()
	PushRenderTarget(multiplicativeLightTexture, false, rl.BlendMultiplied)
	PushRenderTarget(additiveLightTexture, false, rl.BlendAdditive)
}

func populateAdditiveLayer() {
	objs := CurrentMap.World.Objects

	rl.BeginTextureMode(additiveLightTexture)
	{
		rl.ClearBackground(rl.Blank)
		rl.BeginMode2D(RenderCamera)
		{
			for _, o := range objs {
				if !o.HasSpecularLight {
					continue
				}

				in := o.Color
				in.A = 90
				out := o.Color
				out.A = 0
				rl.DrawCircleGradient(
					int32(o.Position.X+16+o.Offset.X),
					int32(o.Position.Y-16+o.Offset.Y),
					float32(o.Radius),
					in,
					out,
				)
			}
		}
		rl.EndMode2D()
	}
	rl.EndTextureMode()

	BlurRenderTarget(additiveLightTexture, 64)
}

func populateMultiplicativeLight() {
	objs := CurrentMap.World.Objects

	rl.BeginTextureMode(multiplicativeLightTexture)
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
	BlurRenderTarget(multiplicativeLightTexture, 32)
}
