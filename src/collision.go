package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/resolv/resolv"
)

type collision struct {
	isColliding bool
}

// NewCollision static map collision
func (o *Object) NewCollision() {
	o.IsCollidable = true
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
	a := o.World.GetObjectsOfType("col", false)

	for _, c := range a {
		if c == o {
			continue
		}

		if o.Class == "col" {
			continue
		}

		col, ok := resolveContact(o, c, deltaX, deltaY)

		if ok {
			return col, true
		}
	}

	d := o.World.GetObjectsOfType("col", true)

	for _, c := range d {
		if !c.IsCollidable {
			continue
		}

		if c == o {
			continue
		}

		col, ok := resolveContact(o, c, deltaX, deltaY)

		if ok {
			return col, true
		}
	}

	return nil, false
}

func resolveContact(a, b *Object, deltaX, deltaY int32) (*resolv.Collision, bool) {
	first := rayRectangleInt32ToResolv(a.GetAABB(a))
	second := rayRectangleInt32ToResolv(b.GetAABB(b))

	try := resolv.Resolve(first, second, deltaX, deltaY)

	if try.Colliding() {
		if DebugMode {
			b.isColliding = true
		}

		return &try, true
	}

	return nil, false
}
