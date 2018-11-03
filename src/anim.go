package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	goaseprite "github.com/solarlune/GoAseprite"
)

type anim struct {
	animStarted bool
}

// NewAnim animated sprite
func (o *Object) NewAnim() {
	o.Ase = goaseprite.Load(fmt.Sprintf("assets/gfx/%s.json", o.FileName))
	o.Texture = GetTexture(fmt.Sprintf("assets/gfx/%s.png", o.FileName))
	o.IsCollidable = true

	o.Trigger = func(o, inst *Object) {
		o.Ase.Play(o.Meta.Properties.GetString("tag"))
		o.animStarted = true
	}

	o.Update = func(o *Object, dt float32) {
		if o.animStarted {
			o.Ase.Update(dt)
		}
	}

	o.GetAABB = getSpriteAABB

	o.Draw = func(o *Object) {
		source := getSpriteRectangle(o)
		dest := getSpriteOrigin(o)

		rl.DrawTexturePro(o.Texture, source, dest, rl.Vector2{}, 0, SkyColor)

		if DebugMode {
			c := getSpriteAABB(o)
			rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
			drawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
		}
	}

	if o.AutoStart {
		o.Trigger(o, nil)
	}
}
