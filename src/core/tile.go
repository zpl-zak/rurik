/*
 * @Author: V4 Games
 * @Date: 2018-12-10 03:31:58
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 04:59:04
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

		o.HorizontalFlip = rawGID&tileHorizontalFlipMask != 0
		o.VerticalFlip = rawGID&tileVerticalFlipMask != 0
		o.DiagonalFlip = rawGID&tileDiagonalFlipMask != 0

		o.Position.Y -= float32(o.Height)
	}

	o.GetAABB = func(o *Object) rl.RectangleInt32 {
		return rl.RectangleInt32{
			X:      int32(o.Position.X),
			Y:      int32(o.Position.Y),
			Width:  32,
			Height: 32,
		}
	}

	o.Draw = func(o *Object) {
		source, tex := CurrentMap.GetTileDataFromID(o.TileID - 1)
		dest := rl.NewRectangle(o.Position.X, o.Position.Y, float32(o.Width), float32(o.Height))

		if DebugMode {
			c := o.GetAABB(o)
			rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
			drawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
		}

		rl.DrawTexturePro(*tex, source, dest, rl.Vector2{}, 0, SkyColor)
	}
}
