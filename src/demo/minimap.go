package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type minimapProg struct {
	RenderTexture system.RenderTarget
}

func newMinimap() *minimapProg {
	v := &minimapProg{
		RenderTexture: system.CreateRenderTarget(320, 320),
	}

	rl.SetTextureFilter(v.RenderTexture.Texture, rl.FilterBilinear)

	return v
}

func (m *minimapProg) Apply() {
	if core.MainCamera == nil {
		return
	}

	minimapCamera := rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(0, 0), 0, 1)
	minimapCamera.Offset = rl.Vector2{
		X: float32(int(-core.MainCamera.Position.X + 320/2)),
		Y: float32(int(-core.MainCamera.Position.Y + 320/2)),
	}

	rl.BeginTextureMode(m.RenderTexture)
	{
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(minimapCamera)
		{
			dbg := core.DebugMode
			sky := core.SkyColor
			core.SkyColor = rl.White
			core.DebugMode = false
			core.DrawMap(false)
			core.DebugMode = dbg
			core.SkyColor = sky
		}
		rl.EndMode2D()
	}
	rl.EndTextureMode()

	core.BlurRenderTarget(m.RenderTexture, 32)
}
