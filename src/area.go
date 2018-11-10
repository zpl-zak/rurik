package main

import (
	"strconv"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"
)

type area struct {
	Radius     int
	isInCircle bool
	Talk       *Object
}

// NewArea trigger zone using various styles of execution
func (o *Object) NewArea() {
	o.Finish = func(o *Object) {
		o.Proxy, _ = o.world.FindObject(o.Meta.Properties.GetString("proxy"))
		o.Radius, _ = strconv.Atoi(o.Meta.Properties.GetString("radius"))
		o.Radius *= 2
	}

	o.Update = func(o *Object, dt float32) {
		vd := raymath.Vector2Distance(getAreaOrigin(o), LocalPlayer.Position)

		if vd < float32(o.Radius) {
			o.isInCircle = true

			if IsKeyPressed("use") {
				o.Trigger(o, LocalPlayer)
			}
		} else {
			o.isInCircle = false
		}
	}

	o.Draw = func(o *Object) {
		if DebugMode {
			pos := getAreaOrigin(o)
			col := rl.NewColor(0, 255, 255, 32)
			if o.isInCircle {
				col = rl.NewColor(255, 0, 255, 32)
			}
			rl.DrawCircle(int32(pos.X), int32(pos.Y), float32(o.Radius), col)
		}
	}

	o.Trigger = func(o, inst *Object) {

		if LocalPlayer.Locked {
			return
		}

		talkFile := o.Meta.Properties.GetString("talk")

		if talkFile != "" {
			if o.Talk == nil {
				o.Talk = o.world.NewObject(nil)
				o.Talk.Name = o.Name + "_talk"

				o.Talk.FileName = talkFile
				o.Talk.NewTalk()
				o.world.Objects = append(o.world.Objects, o.Talk)
			}

			if !o.Talk.Started && (rl.GetTime()-o.Talk.LastTrigger) > 1 {
				o.Talk.Trigger(o.Talk, o)
			}
		}

		if o.Target != nil {
			o.Target.Trigger(o.Target, o)
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
