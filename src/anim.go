package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type anim struct {
	AnimTag     string
	animStarted bool
}

// NewAnim animated sprite
func (o *Object) NewAnim() {
	o.Texture = GetTexture(o.FileName + ".png")
	o.IsCollidable = true

	if o.AnimTag == "" {
		o.AnimTag = o.Meta.Properties.GetString("tag")
	}

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
		if o.Proxy != nil {
			o.Ase = o.Proxy.Ase
		} else {
			aseData := GetAnimData(o.FileName)
			o.Ase = &aseData
		}

		if o.AutoStart {
			o.Trigger(o, nil)
		}
	}

	o.GetAABB = getSpriteAABB

	o.Draw = func(o *Object) {
		source := getSpriteRectangle(o)
		dest := getSpriteOrigin(o)

		if DebugMode {
			c := getSpriteAABB(o)
			rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
			drawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
		}

		rl.DrawTexturePro(o.Texture, source, dest, rl.Vector2{}, 0, SkyColor)
	}
}
