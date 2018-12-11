package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"madaraszd.net/zaklaus/rurik/src/core"
	"madaraszd.net/zaklaus/rurik/src/system"
)

type bloomProg struct {
	TresholdTexture  *rl.RenderTexture2D
	BlurTexture      []*rl.RenderTexture2D
	CompositeTexture *rl.RenderTexture2D
	ExtractColors    system.Program
	BlurImage        system.Program
}

func newBloom() *bloomProg {
	blurTextures := []*rl.RenderTexture2D{
		system.CreateRenderTarget(screenW, screenH),
		system.CreateRenderTarget(screenW, screenH),
	}
	b := &bloomProg{
		TresholdTexture:  system.CreateRenderTarget(screenW, screenH),
		BlurTexture:      blurTextures,
		CompositeTexture: system.CreateRenderTarget(screenW, screenH),
		ExtractColors:    system.NewProgram("", "assets/shaders/extractColors.fs"),
		BlurImage:        system.NewProgram("", "assets/shaders/blur.fs"),
	}

	b.ExtractColors.SetShaderValue("size", []float32{screenW, screenH}, 2)

	return b
}

func (b *bloomProg) Apply() {
	b.ExtractColors.BlitToRenderTarget(core.ScreenTexture, b.TresholdTexture)

	var hor int32 = 1
	first := true
	maxIter := 15

	for i := 0; i < maxIter; i++ {
		var srcTex *rl.RenderTexture2D
		b.BlurImage.SetShaderValuei("horizontal", []int32{hor}, 1)

		if first {
			srcTex = b.TresholdTexture
			first = false
		} else {
			srcTex = b.BlurTexture[1-hor]
		}

		b.BlurImage.BlitToRenderTarget(srcTex, b.BlurTexture[hor])
		hor = 1 - hor
	}

	rl.BeginTextureMode(*b.CompositeTexture)
	rl.BeginBlendMode(rl.BlendAdditive)
	{
		rl.DrawTexture(core.ScreenTexture.Texture, 0, 0, rl.White)
		rl.DrawTexture(b.BlurTexture[1].Texture, 0, 0, rl.White)
	}
	rl.EndBlendMode()
	rl.EndTextureMode()

	system.CopyToRenderTarget(b.CompositeTexture, core.ScreenTexture)
}
