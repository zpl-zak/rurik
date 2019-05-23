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
	"github.com/solarlune/resolv/resolv"
	rl "github.com/zaklaus/raylib-go/raylib"
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
	o.DebugVisible = false

	o.Draw = func(o *Object) {
		if !DebugMode || !o.DebugVisible {
			return
		}

		color := rl.White

		if o.isColliding {
			color = rl.Red
			o.isColliding = false
		}

		rl.DrawRectangleLines(int32(o.Position.X), int32(o.Position.Y), int32(o.Meta.Width), int32(o.Meta.Height), color)

		c := o.GetAABB(o)
		DrawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
	}

	o.GetAABB = GetSolidAABB
}

// CheckForCollision performs collision detection and resolution
func CheckForCollision(o *Object, deltaX, deltaY int32) (resolv.Collision, bool) {
	collisionProfiler.StartInvocation()

	if !o.IsCollidable {
		collisionProfiler.StopInvocation()
		return resolv.Collision{}, false
	}

	for _, c := range o.world.Objects {
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
