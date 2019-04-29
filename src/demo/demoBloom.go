package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type bloomProg struct {
	TresholdTexture system.RenderTarget
	BlurTexture     []system.RenderTarget
	ExtractColors   system.Program
	BlurImage       system.Program
}

func newBloom() *bloomProg {
	b := &bloomProg{
		TresholdTexture: system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight),
		ExtractColors:   system.NewProgram("", "shaders/extractColors.fs"),
	}

	b.ExtractColors.SetShaderValue("size", []float32{float32(system.ScreenWidth), float32(system.ScreenHeight)}, 2)

	return b
}

func (b *bloomProg) Apply() {
	if core.WindowWasResized {
		rl.UnloadRenderTexture(b.TresholdTexture)
		b.TresholdTexture = system.CreateRenderTarget(system.ScreenWidth, system.ScreenHeight)
		b.ExtractColors.SetShaderValue("size", []float32{float32(system.ScreenWidth), float32(system.ScreenHeight)}, 2)
	}

	b.ExtractColors.RenderToTexture(core.WorldTexture, b.TresholdTexture)

	core.BlurRenderTarget(b.TresholdTexture, 10)
	core.PushRenderTarget(b.TresholdTexture, true, rl.BlendAdditive)
}
