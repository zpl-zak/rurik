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

	o.Draw = func(o *Object) {
		if !DebugMode {
			return
		}

		color := rl.White

		if o.isColliding {
			color = rl.Red
			o.isColliding = false
		}

		rl.DrawRectangleLines(int32(o.Position.X), int32(o.Position.Y), int32(o.Meta.Width), int32(o.Meta.Height), color)

		c := o.GetAABB(o)
		drawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
	}

	o.GetAABB = func(o *Object) rl.RectangleInt32 {
		return rl.RectangleInt32{
			X:      int32(o.Position.X),
			Y:      int32(o.Position.Y),
			Width:  o.Size[0],
			Height: o.Size[1],
		}
	}
}

// CheckForCollision performs collision detection and resolution
func CheckForCollision(o *Object, deltaX, deltaY int32) (resolv.Collision, bool) {
	collisionProfiler.StartInvocation()

	if !o.IsCollidable {
		collisionProfiler.StopInvocation()
		return resolv.Collision{}, false
	}

	for _, c := range o.World.Objects {
		col, ok := resolveContact(o, c, deltaX, deltaY)

		if ok {
			collisionProfiler.StopInvocation()
			return col, true
		}
	}

	collisionProfiler.StopInvocation()
	return resolv.Collision{}, false
}

var (
	resolveFirst  resolv.Rectangle
	resolveSecond resolv.Rectangle
)

func resolveContact(a, b *Object, deltaX, deltaY int32) (resolv.Collision, bool) {

	if !b.IsCollidable || a == b {
		return resolv.Collision{}, false
	}

	rayRectangleInt32ToResolv(&resolveFirst, a.GetAABB(a))
	rayRectangleInt32ToResolv(&resolveSecond, b.GetAABB(b))

	try := resolv.Resolve(&resolveFirst, &resolveSecond, deltaX, deltaY)

	if try.Colliding() {
		if DebugMode {
			b.isColliding = true
		}

		return try, true
	}

	return resolv.Collision{}, false
}
