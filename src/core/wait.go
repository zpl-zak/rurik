/*
   Copyright 2018 V4 Games

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
	"fmt"
	"log"
	"strconv"

	rl "github.com/zaklaus/raylib-go/raylib"
)

type wait struct {
	Duration  int
	Remaining int
	Script    *Object
}

// NewWait timer/unconditioned trigger instigator
func (o *Object) NewWait() {
	o.Duration, _ = strconv.Atoi(o.Meta.Properties.GetString("duration"))
	o.Trigger = triggerWait
	o.Update = updateWait

	o.Init = func(o *Object) {
		if o.AutoStart {
			o.Trigger(o, nil)
		}
	}

	o.Draw = func(o *Object) {
		if !DebugMode || !o.DebugVisible {
			return
		}

		rl.DrawCircle(int32(o.Position.X), int32(o.Position.Y), 5, rl.White)
		DrawTextCentered(fmt.Sprintf("%s", o.Name), int32(o.Position.X), int32(o.Position.Y)+5, 10, rl.White)
		DrawTextCentered(fmt.Sprintf("Remaining: %d ms", o.Remaining), int32(o.Position.X), int32(o.Position.Y)+15, 10, rl.White)
	}
}

func updateWait(o *Object, dt float32) {
	if !o.Started {
		return
	}

	if o.Remaining < 0 {
		o.Started = false
		o.Remaining = 0

		fireWait(o)

		return
	}

	o.Remaining -= int(1000.0 * dt)
}

func fireWait(o *Object) {
	if o.Target != nil {
		o.Target.Trigger(o.Target, o)
	}

	if o.Script != nil {
		o.Script.Trigger(o.Script, o)
	}

	if o.Target == nil && o.Script == nil {
		log.Printf("Timer %s has no target attached!\n", o.Name)
	}
}

func triggerWait(o, inst *Object) {
	scriptName := o.Meta.Properties.GetString("file")

	if scriptName != "" {
		o.Script = o.world.NewObjectPro(o.Name+"_script", "script")
		o.Script.FileName = scriptName
		o.Script.IsPersistent = false
		o.world.FinalizeObject(o.Script)
	}

	if o.Duration == 0 {
		fireWait(o)
		return
	}
	o.Remaining = o.Duration
	o.Started = true
}
