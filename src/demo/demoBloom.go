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

	var hor int32 = 1
	maxIter := 10
	srcTex := b.TresholdTexture

	for i := 0; i < maxIter; i++ {
		b.BlurImage.SetShaderValuei("horizontal", []int32{hor}, 1)

		b.BlurImage.RenderToTexture(srcTex, b.BlurTexture[hor])
		srcTex = b.BlurTexture[hor]
		hor = 1 - hor
	}

	core.PushRenderTarget(b.BlurTexture[1], false, rl.BlendAdditive)
}
