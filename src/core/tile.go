/*
 * @Author: V4 Games
 * @Date: 2018-12-10 03:31:58
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 14:51:42
 */

package core

import (
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
		//o.DebugVisible = false

		o.HorizontalFlip = rawGID&tileHorizontalFlipMask != 0
		o.VerticalFlip = rawGID&tileVerticalFlipMask != 0
		o.DiagonalFlip = rawGID&tileDiagonalFlipMask != 0

		o.IsCollidable = o.Meta.Properties.GetString("colType") == "" || (o.Meta.Properties.GetString("colType") != "" && o.Meta.Properties.GetString("colType") != "none")
	}

	o.GetAABB = func(o *Object) rl.RectangleInt32 {
		return rl.RectangleInt32{
			X:      int32(o.Position.X),
			Y:      int32(o.Position.Y) - int32(o.Height),
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
