/*
   Copyright 2019 Dominik Madarász <zaklaus@madaraszd.net>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package core

import (
	"strings"

	"github.com/solarlune/resolv/resolv"
	rl "github.com/zaklaus/raylib-go/raylib"
)

var (
	colDbgX int32
	colDbgY int32
)

type collision struct {
	isColliding      bool
	ContainedObjects []TriggerContact
}

// TriggerContact specifies a trigger point
type TriggerContact struct {
	Object     *Object
	Res        resolv.Collision
	wasUpdated bool
}

// NewCollision static map collision
func (o *Object) NewCollision() {
	o.IsCollidable = true
	o.Size = []int32{int32(o.Meta.Width), int32(o.Meta.Height)}

	if o.PolyLines != nil {
		o.CollisionType = "slope"
	} else if o.CollisionType == "" {
		o.CollisionType = "solid"
	}

	o.Draw = func(o *Object) {
		if !DebugMode || !o.DebugVisible {
			return
		}

		color := rl.White

		if o.isColliding {
			color = rl.Red
			o.isColliding = false
		}

		rl.DrawCircle(colDbgX, colDbgY, 5, rl.Blue)

		if o.PolyLines != nil {
			for _, pl := range o.PolyLines {
				for idx := 0; idx < len(*pl.Points)-1; idx++ {
					pts := *pl.Points
					p0 := pts[idx+0]
					p1 := pts[idx+1]

					rl.DrawLine(
						int32(o.Position.X)+int32(p0.X),
						int32(o.Position.Y)+int32(p0.Y),
						int32(o.Position.X)+int32(p1.X),
						int32(o.Position.Y)+int32(p1.Y),
						color,
					)
				}
			}
		} else {
			rl.DrawRectangleLines(int32(o.Position.X), int32(o.Position.Y), int32(o.Meta.Width), int32(o.Meta.Height), color)
		}

		c := o.GetAABB(o)
		DrawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
	}

	o.GetAABB = GetSolidAABB
}

// CheckForCollision performs collision detection and resolution
func CheckForCollision(o *Object, deltaX, deltaY int32) (resolv.Collision, bool) {
	return CheckForCollisionEx("solid+trigger", o, deltaX, deltaY)
}

// CheckForCollisionEx performs collision detection and resolution
func CheckForCollisionEx(collisionType string, o *Object, deltaX, deltaY int32) (resolv.Collision, bool) {
	collisionProfiler.StartInvocation()

	if !o.IsCollidable {
		collisionProfiler.StopInvocation()
		return resolv.Collision{}, false
	}

	colTypes := strings.Split(collisionType, "+")

	for _, c := range o.world.Objects {
		sig := false
		for _, ct := range colTypes {
			if c.CollisionType == ct || ct == "*" {
				sig = true
			}
		}

		if !sig {
			continue
		}

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

	var try resolv.Collision

	// NOTE: Slope handling
	if b.PolyLines != nil && b.CollisionType == "slope" {
		for _, pl := range b.PolyLines {
			done := false
			for idx := 0; idx < len(*pl.Points)-1; idx++ {
				pts := *pl.Points
				p0 := pts[idx+0]
				p1 := pts[idx+1]
				line := resolv.NewLine(
					int32(b.Position.X)+int32(p0.X),
					int32(b.Position.Y)+int32(p0.Y),
					int32(b.Position.X)+int32(p1.X),
					int32(b.Position.Y)+int32(p1.Y),
				)

				try = resolv.Resolve(&resolveFirst, line, 0, deltaY)

				if try.Colliding() {
					xpos := a.Position.X - b.Position.X
					m := float32(p1.Y-p0.Y) / float32(p1.X-p0.X)
					bc := float32(p0.Y) - (m * float32(p0.X))
					ypos := m*(xpos) + bc

					if DebugMode {
						colDbgX = int32(b.Position.X) + int32(xpos)
						colDbgY = int32(b.Position.Y) + int32(ypos)
					}

					a.Position.Y = float32(b.Position.Y) + float32(ypos) - 20

					if try.Teleporting {
						try.ResolveX = deltaX
						try.Teleporting = false
					}

					done = true
					break
				}
			}

			if done {
				break
			}
		}
	}

	if !try.Colliding() && b.CollisionType != "slope" {
		rayRectangleInt32ToResolv(&resolveSecond, b.GetAABB(b))
		try = resolv.Resolve(&resolveFirst, &resolveSecond, deltaX, deltaY)
	}

	if try.Colliding() {
		if DebugMode {
			b.isColliding = true
		}

		a.HandleCollision(&try, a, b)
		b.HandleCollision(&try, b, a)

		if b.CollisionType == "trigger" {
			ct := findExistingContainedObject(b, a, try)

			if ct == nil {
				a.HandleCollisionEnter(&try, a, b)
				b.HandleCollisionEnter(&try, b, a)

				ctx := TriggerContact{
					Object:     a,
					Res:        try,
					wasUpdated: true,
				}

				b.ContainedObjects = append(b.ContainedObjects, ctx)
			} else {
				ct.wasUpdated = true
				ct.Res = try
			}

			return resolv.Collision{}, false
		}

		return try, true
	}

	return resolv.Collision{}, false
}

func findExistingContainedObject(o, other *Object, res resolv.Collision) *TriggerContact {
	for k := range o.ContainedObjects {
		v := &o.ContainedObjects[k]

		if v.Object == other {
			return v
		}
	}

	return nil
}

func (o *Object) updateTriggerArea() {
	newObjects := []TriggerContact{}

	for k := range o.ContainedObjects {
		v := &o.ContainedObjects[k]

		if v.wasUpdated {
			v.wasUpdated = false
			newObjects = append(newObjects, *v)
		} else {
			v.Object.HandleCollisionLeave(&v.Res, v.Object, o)
			o.HandleCollisionLeave(&v.Res, o, v.Object)
		}
	}

	o.ContainedObjects = newObjects
}
