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
	blurTextures := []system.RenderTarget{
		system.CreateRenderTarget(screenW, screenH),
		system.CreateRenderTarget(screenW, screenH),
	}
	b := &bloomProg{
		TresholdTexture: system.CreateRenderTarget(screenW, screenH),
		BlurTexture:     blurTextures,
		ExtractColors:   system.NewProgram("", "assets/shaders/extractColors.fs"),
		BlurImage:       system.NewProgram("", "assets/shaders/blur.fs"),
	}

	b.ExtractColors.SetShaderValue("size", []float32{screenW, screenH}, 2)

	return b
}

func (b *bloomProg) Apply() {
	if core.WindowWasResized {
		rl.UnloadRenderTexture(*b.BlurTexture[0])
		rl.UnloadRenderTexture(*b.BlurTexture[1])
		rl.UnloadRenderTexture(*b.TresholdTexture)

		blurTextures := []system.RenderTarget{
			system.CreateRenderTarget(screenW, screenH),
			system.CreateRenderTarget(screenW, screenH),
		}

		b.TresholdTexture = system.CreateRenderTarget(screenW, screenH)
		b.BlurTexture = blurTextures

		b.ExtractColors.SetShaderValue("size", []float32{float32(system.ScreenWidth), float32(system.ScreenHeight)}, 2)
		b.BlurImage.SetShaderValue("size", []float32{float32(system.ScreenWidth), float32(system.ScreenHeight)}, 2)
	}

	b.ExtractColors.RenderToTexture(core.WorldTexture, b.TresholdTexture)

	system.BlurRenderTarget(b.TresholdTexture, 10)

	core.PushRenderTarget(b.TresholdTexture, false, rl.BlendAdditive)
}
