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
	additiveLightTexture       system.RenderTarget
	multiplicativeLightTexture system.RenderTarget

	showLightmap bool
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
	lmState := rl.BlendMultiplied

	if showLightmap {
		lmState = rl.BlendAlpha
	}

	PushRenderTarget(multiplicativeLightTexture, false, lmState)
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
				if !o.HasSpecularLight || !o.Visible {
					continue
				}

				in := o.Color
				in.A = 90
				out := o.Color
				out.A = 0
				rl.DrawCircleGradient(
					int32(o.Position.X+o.Offset.X),
					int32(o.Position.Y+o.Offset.Y),
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
		rl.ClearBackground(SkyColor) //Vec3ToColor(LerpColor(ColorToVec3(SkyColor), ColorToVec3(rl.Black), 0.8)))

		rl.BeginBlendMode(rl.BlendAdditive)
		{
			rl.BeginMode2D(RenderCamera)
			{
				for _, o := range objs {
					if !o.HasLight || !o.Visible {
						continue
					}

					rl.DrawCircleGradient(
						int32(o.Position.X+o.Offset.X),
						int32(o.Position.Y+o.Offset.Y),
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
