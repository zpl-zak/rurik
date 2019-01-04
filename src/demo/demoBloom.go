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
		TresholdTexture: system.CreateRenderTarget(screenW, screenH),
		ExtractColors:   system.NewProgram("", "assets/shaders/extractColors.fs"),
	}

	b.ExtractColors.SetShaderValue("size", []float32{screenW, screenH}, 2)

	return b
}

func (b *bloomProg) Apply() {
	if core.WindowWasResized {
		b.TresholdTexture = system.CreateRenderTarget(screenW, screenH)
		b.ExtractColors.SetShaderValue("size", []float32{float32(system.ScreenWidth), float32(system.ScreenHeight)}, 2)
	}

	b.ExtractColors.RenderToTexture(core.WorldTexture, b.TresholdTexture)

	core.BlurRenderTarget(b.TresholdTexture, 10)
	core.PushRenderTarget(b.TresholdTexture, true, rl.BlendAdditive)
}
