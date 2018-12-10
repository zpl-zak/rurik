/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:28:08
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 14:23:42
 */

package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// NewTarget dummy target point
func (o *Object) NewTarget() {
	o.Draw = func(o *Object) {
		if !DebugMode || !o.DebugVisible {
			return
		}

		rl.DrawCircle(int32(o.Position.X), int32(o.Position.Y), 5, rl.White)
		drawTextCentered(o.Name, int32(o.Position.X), int32(o.Position.Y)+5, 10, rl.White)
	}
}
