/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:26:11
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-08 21:02:24
 */

package core

import (
	"strconv"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"
	"madaraszd.net/zaklaus/rurik/src/system"
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

			if system.IsKeyPressed("use") {
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
				o.Talk = o.world.NewObjectPro(o.Name+"_talk", "talk")
				o.Talk.FileName = talkFile
				o.Talk.CanRepeat = true
				o.world.AddObject(o.Talk)
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
