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
	"log"
	"strconv"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/raylib-go/raymath"
)

type area struct {
	isInCircle     bool
	Talk           *Object
	triggerOnHover bool
}

// NewArea trigger zone using various styles of execution
func (o *Object) NewArea() {
	o.Finish = func(o *Object) {
		o.Proxy, _ = o.world.FindObject(o.Meta.Properties.GetString("proxy"))
		r, _ := strconv.ParseFloat(o.Meta.Properties.GetString("radius"), 32)
		o.triggerOnHover = o.Meta.Properties.GetString("onhover") == "1"
		o.Radius = float32(r)
		o.Radius *= 2
	}

	o.Update = func(o *Object, dt float32) {
		hit := false
		for _, obj := range o.world.Objects {
			if obj.CanTrigger == false || obj == o.Proxy {
				continue
			}

			vd := raymath.Vector2Distance(getAreaOrigin(o), obj.Position)

			if vd < float32(o.Radius) {
				hit = true
				if obj.InsideArea(obj, o) || o.triggerOnHover {
					o.Trigger(o, obj)
				}
			}
		}

		o.isInCircle = hit
	}

	o.Draw = func(o *Object) {
		if DebugMode && o.DebugVisible {
			pos := getAreaOrigin(o)
			col := rl.NewColor(0, 255, 255, 32)
			if o.isInCircle {
				col = rl.NewColor(255, 0, 255, 32)
			}
			rl.DrawCircle(int32(pos.X), int32(pos.Y), float32(o.Radius), col)
		}
	}

	o.Trigger = func(o, inst *Object) {
		if o.EventName != "" {
			FireEvent(o.EventName, o.EventArgs)
		} else {
			log.Printf("Talk object '%s' has no event attached!\n", o.Name)
		}
	}
}

func getAreaOrigin(o *Object) rl.Vector2 {
	if o.Proxy != nil {
		p := o.Proxy.Position
		b := o.Proxy.GetAABB(o.Proxy)
		p.Y += float32(b.Height / 2)

		return p
	}

	return o.Position
}
