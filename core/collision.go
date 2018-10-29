package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv/resolv"
)

type collision struct {
	colType     string
	isColliding bool
}

// NewCollision instance
func (o *Object) NewCollision() {
	o.colType = o.Meta.Properties.GetString("colType")
	o.Size = []int32{int32(o.Meta.Width), int32(o.Meta.Height)}
	o.Draw = drawCollision
	o.GetAABB = getCollisionAABB
}

func getCollisionAABB(o *Object) rl.RectangleInt32 {
	return rl.RectangleInt32{
		X:      int32(o.Position.X),
		Y:      int32(o.Position.Y),
		Width:  o.Size[0],
		Height: o.Size[1],
	}
}

func drawCollision(o *Object) {
	if !DebugMode {
		return
	}

	color := rl.White

	if o.isColliding {
		color = rl.Red
		o.isColliding = false
	}

	rl.DrawRectangleLines(int32(o.Position.X), int32(o.Position.Y), int32(o.Meta.Width), int32(o.Meta.Height), color)

	c := getCollisionAABB(o)
	drawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
}

// CheckForCollision performs collision detection and resolution
func CheckForCollision(o *Object, deltaX, deltaY int32) (*resolv.Collision, bool) {
	a := GetObjectsOfType("col", false)

	for _, c := range a {
		if c == o {
			continue
		}

		if o.Class == "col" {
			continue
		}

		first := rayRectangleInt32ToResolv(o.GetAABB(o))
		second := rayRectangleInt32ToResolv(c.GetAABB(c))

		try := resolv.Resolve(first, second, deltaX, deltaY)

		if try.Colliding() {
			if DebugMode {
				c.isColliding = true
			}

			return &try, true
		}
	}

	return nil, false
}
