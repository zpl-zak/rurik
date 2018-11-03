package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// NewTarget dummy target point
func (o *Object) NewTarget() {
	o.Draw = func(o *Object) {
		if !DebugMode {
			return
		}

		rl.DrawCircle(int32(o.Position.X), int32(o.Position.Y), 5, rl.White)
		drawTextCentered(o.Name, int32(o.Position.X), int32(o.Position.Y)+5, 10, rl.White)
	}
}
