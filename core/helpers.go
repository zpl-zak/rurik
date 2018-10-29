package core

import (
	"../system"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv/resolv"
)

func rayRectangleInt32ToResolv(i rl.RectangleInt32) *resolv.Rectangle {
	return &resolv.Rectangle{
		BasicShape: resolv.BasicShape{
			X:           i.X,
			Y:           i.Y,
			Collideable: true,
		},
		W: i.Width,
		H: i.Height,
	}
}

func drawTextCentered(text string, posX, posY, fontSize int32, color rl.Color) {
	if fontSize < 10 {
		fontSize = 10
	}

	rl.DrawText(text, posX-rl.MeasureText(text, fontSize)/2, posY, fontSize, color)
}

func vector2Lerp(v1, v2 rl.Vector2, amount float32) (result rl.Vector2) {
	result.X = v1.X + amount*(v2.X-v1.X)
	result.Y = v1.Y + amount*(v2.Y-v1.Y)

	return result
}

func scalarLerp(v1, v2 float32, amount float32) (result float32) {
	result = v1 + amount*(v2-v1)

	return result
}

func isMouseInRectangle(x, y, x2, y2 int32) bool {
	x2 = x + x2
	y2 = y + y2

	mo := rl.GetMousePosition()
	m := [2]int32{
		int32(mo.X) / system.ScaleRatio,
		int32(mo.Y) / system.ScaleRatio,
	}

	if m[0] > x && m[0] < x2 &&
		m[1] > y && m[1] < y2 {
		return true
	}

	return false
}
