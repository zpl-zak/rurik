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
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

type anim struct {
	AnimTag             string
	animStarted         bool
	PendingCurrentFrame int32
}

// NewAnim animated sprite
func (o *Object) NewAnim() {
	o.Trigger = func(o, inst *Object) {
		if o.Ase != nil {
			o.Ase.Play(o.AnimTag)

			if o.PendingCurrentFrame != -1 {
				o.Ase.CurrentFrame = o.PendingCurrentFrame
				o.PendingCurrentFrame = -1
			}
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

	o.GetAABB = GetSpriteAABB

	o.Draw = func(o *Object) {
		if o.Ase == nil {
			return
		}

		source := GetSpriteRectangle(o)
		dest := GetSpriteOrigin(o)

		if DebugMode && o.DebugVisible {
			c := GetSpriteAABB(o)
			rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
			DrawTextCentered(o.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
		}

		rl.DrawTexturePro(*o.Texture, source, dest, rl.Vector2{}, o.Rotation, SkyColor)
	}
}
