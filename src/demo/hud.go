package main

import (
	"log"
	"math/rand"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

const (
	frameWidth           = 52
	frameHeight          = 20
	frameSrcPosX         = 0
	frameSrcPosY         = 0
	frameDestPosX        = 5
	frameDestPosY        = 15
	frameScaling         = 2.75
	frameOverlayPosX     = 6
	frameOverlayPosY     = 64
	frameOverlayDestPosX = 5
	frameOverlayDestPosY = 9

	valueUpdatePauseTime = 0.3
	valueUpdateTime      = 0.4
)

type barStat struct {
	Value  float32
	SrcX   float32
	SrcY   float32
	DestX  float32
	DestY  float32
	Width  float32
	Height float32

	RemainingValueOffset float32
	ValueUpdateDelayTime float32
	ValueUpdateTime      float32
}

var (
	hudTexture *rl.Texture2D
	barStats   []barStat
)

const (
	barHealth = iota
	barMagicka
	barUltimate
)

func initHUD() {
	hudTexture = system.GetTexture("gfx/gamehud.png")

	if hudTexture == nil {
		log.Fatalln("HUD texture not found!")
		return
	}

	barStats = []barStat{
		barStat{
			Value:  50,
			SrcX:   0,
			SrcY:   48,
			DestX:  7,
			DestY:  10,
			Width:  44,
			Height: 3,
		},
		barStat{
			Value:  50,
			SrcX:   0,
			SrcY:   51,
			DestX:  10,
			DestY:  14,
			Width:  37,
			Height: 2,
		},
		barStat{
			Value:  50,
			SrcX:   0,
			SrcY:   53,
			DestX:  11,
			DestY:  16,
			Width:  23,
			Height: 1,
		},
	}
}

func drawHUD() {
	{
		frameRect := rl.NewRectangle(frameSrcPosX, frameSrcPosY, frameWidth, frameHeight)
		frameDestRect := rl.NewRectangle(frameDestPosX, frameDestPosY, frameWidth*frameScaling, frameHeight*frameScaling)
		rl.DrawTexturePro(
			*hudTexture,
			frameRect,
			frameDestRect,
			rl.Vector2{},
			0,
			rl.White,
		)
	}

	for _, v := range barStats {
		rect := rl.NewRectangle(v.SrcX, v.SrcY, v.Value*v.Width/100, v.Height)
		barWidth := float32(v.Value*v.Width/100) * frameScaling
		destRect := rl.NewRectangle(
			v.DestX*frameScaling+frameDestPosX,
			v.DestY*frameScaling+frameDestPosY,
			barWidth,
			v.Height*frameScaling,
		)

		offsetRect := rl.NewRectangle(
			v.DestX*frameScaling+frameDestPosX+barWidth,
			v.DestY*frameScaling+frameDestPosY,
			float32(v.RemainingValueOffset*v.Width/100)*frameScaling,
			v.Height*frameScaling,
		)

		if v.RemainingValueOffset > 0 {

			offset := v.Value - v.RemainingValueOffset
			barWidthS := float32(offset*v.Width/100) * frameScaling
			barWidthE := float32(v.Value*v.Width/100) * frameScaling
			barWidth = core.ScalarLerp(barWidthS, barWidthE, 1-v.ValueUpdateTime/valueUpdateTime)
			rect.Width = core.ScalarLerp(v.Value*v.Width/100, offset*v.Width/100, v.ValueUpdateTime/valueUpdateTime)
			offsetRect.X = v.DestX*frameScaling + frameDestPosX + barWidthS
			destRect.Width = barWidth
			rl.DrawRectangleRec(offsetRect, rl.RayWhite)
		} else if v.RemainingValueOffset < 0 {
			offset := -v.RemainingValueOffset
			offset = core.ScalarLerp(offset, 0, 1-v.ValueUpdateTime/valueUpdateTime)
			offsetRect.Width = float32(offset*v.Width/100) * frameScaling
			rl.DrawRectangleRec(offsetRect, rl.RayWhite)
		}

		rl.DrawTexturePro(
			*hudTexture,
			rect,
			destRect,
			rl.Vector2{},
			0,
			rl.White,
		)
	}

	if false {
		frameRect := rl.NewRectangle(frameOverlayPosX, frameOverlayPosY, frameWidth, frameHeight)
		frameDestRect := rl.NewRectangle(
			frameDestPosX+frameOverlayDestPosX*frameScaling,
			frameDestPosY+frameOverlayDestPosY*frameScaling,
			frameWidth*frameScaling,
			frameHeight*frameScaling,
		)
		rl.DrawTexturePro(
			*hudTexture,
			frameRect,
			frameDestRect,
			rl.Vector2{},
			0,
			rl.White,
		)
	}
}

func applyBarStatValueChange(stat int, value float32) {
	v := &barStats[stat]

	value = clamp(value, -v.Value, 100-v.Value)

	v.RemainingValueOffset = value
	v.ValueUpdateDelayTime = valueUpdatePauseTime
	v.ValueUpdateTime = valueUpdateTime
	v.Value += value
}

func updateHUD() {
	// NOTE: temp test
	{
		if rl.IsKeyPressed(rl.KeyF8) {
			for k := range barStats {
				applyBarStatValueChange(k, (rand.Float32()*2-1)*10)
			}
		}
	}

	for k := range barStats {
		v := &barStats[k]

		if v.RemainingValueOffset != 0 {
			if v.ValueUpdateDelayTime > 0 {
				v.ValueUpdateDelayTime -= system.FrameTime

				if v.ValueUpdateDelayTime < 0 {
					v.ValueUpdateDelayTime = 0
				}
			}

			if v.ValueUpdateDelayTime == 0 && v.ValueUpdateTime > 0 {
				v.ValueUpdateTime -= system.FrameTime

				if v.ValueUpdateTime < 0 {
					v.ValueUpdateTime = 0
				}
			}

			if v.ValueUpdateDelayTime == 0 && v.ValueUpdateTime == 0 {
				v.RemainingValueOffset = 0
			}
		}
	}
}

func clamp(x, a, b float32) float32 {
	if x < a {
		return a
	} else if x > b {
		return b
	}

	return x
}
