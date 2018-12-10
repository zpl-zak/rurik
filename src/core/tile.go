/*
 * @Author: V4 Games
 * @Date: 2018-12-10 03:31:58
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 17:45:39
 */

package core

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type tile struct {
	Width          int32
	Height         int32
	TileID         int
	HorizontalFlip bool
	VerticalFlip   bool
	DiagonalFlip   bool
}

// NewTile object tile
func (o *Object) NewTile() {
	o.Finish = func(o *Object) {
		rawGID := o.Meta.GID
		o.Meta.GID = rawGID &^ tileFlip
		o.TileID = int(o.Meta.GID)
		o.Width = int32(o.Meta.Width)
		o.Height = int32(o.Meta.Height)
		o.Ase = nil
		o.DebugVisible = false

		o.HorizontalFlip = rawGID&tileHorizontalFlipMask != 0
		o.VerticalFlip = rawGID&tileVerticalFlipMask != 0
		o.DiagonalFlip = rawGID&tileDiagonalFlipMask != 0

		o.IsCollidable = o.Meta.Properties.GetString("colType") == "" || (o.Meta.Properties.GetString("colType") != "" && o.Meta.Properties.GetString("colType") != "none")
	}

	o.GetAABB = func(o *Object) rl.RectangleInt32 {
		centerX := float64(o.Width) / 2
		centerY := -float64(o.Height) / 2
		cosR := math.Cos(float64(o.Rotation) / (180 / math.Pi))
		sinR := math.Sin(float64(o.Rotation) / (180 / math.Pi))
		rotCenterX := int32(centerX*cosR - centerY*sinR)
		rotCenterY := int32(centerX*sinR + centerY*cosR)

		return rl.RectangleInt32{
			X:      int32(o.Position.X) + rotCenterX - int32(centerX),
			Y:      int32(o.Position.Y) + rotCenterY + int32(centerY),
			Width:  o.Width,
			Height: o.Height,
		}
	}

	o.Update = func(o *Object, dt float32) {
		if rl.IsKeyDown(rl.KeyF) {
			o.Rotation++
		}
	}

	o.Draw = func(o *Object) {
		source, tex := CurrentMap.GetTileDataFromID(o.TileID - 1)
		dest := rl.NewRectangle(o.Position.X, o.Position.Y, float32(o.Width), float32(o.Height))

		if DebugMode && o.DebugVisible {
			c := o.GetAABB(o)
			rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
			{
				b := rl.Vector2{X: dest.X, Y: dest.Y}
				e := rl.Vector2{X: float32(c.X), Y: float32(c.Y)}
				rl.DrawCircle(int32(b.X), int32(b.Y), 2, rl.Green)
				rl.DrawCircle(int32(e.X), int32(e.Y), 2, rl.Red)
				rl.DrawLineEx(b, e, 1, rl.Yellow)
			}
			drawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
		}

		var rot float32

		if o.HorizontalFlip {
			source.Width *= -1
		}

		if o.VerticalFlip {
			source.Height *= -1
		}

		if o.DiagonalFlip {
			source.Width *= -1
			rot = 90
		}

		rl.DrawTexturePro(*tex, source, dest, rl.Vector2{X: 0, Y: float32(o.Height)}, rot+o.Rotation, SkyColor)
	}
}