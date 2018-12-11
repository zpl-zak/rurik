package main

import (
	"github.com/zaklaus/raylib-go/raylib"
	"madaraszd.net/zaklaus/rurik/src/core"
	"madaraszd.net/zaklaus/rurik/src/system"
)

type minimapProg struct {
	RenderTexture system.RenderTarget
	WorldTexture  system.RenderTarget
	SobelEffect   system.Program
}

func newMinimap() *minimapProg {
	return &minimapProg{
		RenderTexture: system.CreateRenderTarget(screenW, screenH),
		WorldTexture:  system.CreateRenderTarget(screenW, screenH),
		SobelEffect:   system.NewProgram("", "assets/shaders/sobel.fs"),
	}
}

func (m *minimapProg) Apply() {
	minimapCamera := rl.NewCamera2D(rl.NewVector2(0, 0), rl.NewVector2(0, 0), 0, 1)
	minimapCamera.Offset = rl.Vector2{
		X: float32(int(-core.MainCamera.Position.X + screenW/2)),
		Y: float32(int(-core.MainCamera.Position.Y + screenH/2)),
	}

	rl.BeginTextureMode(*m.WorldTexture)
	{
		rl.BeginMode2D(minimapCamera)
		{
			dbg := core.DebugMode
			core.DebugMode = false
			core.DrawMap(false)
			core.DebugMode = dbg
		}
		rl.EndMode2D()
	}
	rl.EndTextureMode()

	m.SobelEffect.RenderToTexture(m.WorldTexture, m.RenderTexture)
}
