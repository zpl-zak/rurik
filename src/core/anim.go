/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:25:59
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 22:13:24
 */

package core

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"madaraszd.net/zaklaus/rurik/src/system"
)

type anim struct {
	AnimTag     string
	animStarted bool
}

// NewAnim animated sprite
func (o *Object) NewAnim() {
	o.Trigger = func(o, inst *Object) {
		if o.Ase != nil {
			o.Ase.Play(o.AnimTag)
		}
		o.animStarted = true
	}

	o.Update = func(o *Object, dt float32) {
		if o.animStarted && o.Proxy == nil {
			o.Ase.Update(dt)
		}
	}

	o.Finish = func(o *Object) {
		o.Texture = system.GetTexture(o.FileName + ".png")

		if o.Proxy != nil {
			o.Ase = o.Proxy.Ase
		} else {
			if o.AnimTag == "" {
				o.AnimTag = o.Meta.Properties.GetString("tag")
			}

			aseData := system.GetAnimData(o.FileName)
			o.Ase = &aseData
		}

		if o.AutoStart {
			o.Trigger(o, nil)
		}
	}

	o.GetAABB = getSpriteAABB

	o.Draw = func(o *Object) {
		if o.Ase == nil {
			return
		}

		source := getSpriteRectangle(o)
		dest := getSpriteOrigin(o)

		if DebugMode && o.DebugVisible {
			c := getSpriteAABB(o)
			rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
			drawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
		}

		rl.DrawTexturePro(*o.Texture, source, dest, rl.Vector2{}, o.Rotation, SkyColor)
	}
}
