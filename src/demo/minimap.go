package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type minimapProg struct {
	RenderTexture system.RenderTarget
	WorldTexture  system.RenderTarget
	SobelEffect   system.Program
}

func newMinimap() *minimapProg {
	return &minimapProg{
		RenderTexture: system.CreateRenderTarget(screenW, screenH),
	}
}

func (m *minimapProg) Apply() {
	minimapCamera := rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(0, 0), 0, 1)
	minimapCamera.Offset = rl.Vector2{
		X: float32(int(-core.MainCamera.Position.X + screenW/2)),
		Y: float32(int(-core.MainCamera.Position.Y + screenH/2)),
	}

	rl.BeginTextureMode(*m.RenderTexture)
	{
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
}
