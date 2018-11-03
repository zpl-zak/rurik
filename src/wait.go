package main

import (
	"fmt"
	"log"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type wait struct {
	Duration  int
	Remaining int
	Script    *Object
}

// NewWait instance
func (o *Object) NewWait() {
	o.Duration, _ = strconv.Atoi(o.Meta.Properties.GetString("duration"))
	o.Trigger = triggerWait
	o.Update = updateWait

	scriptName := o.Meta.Properties.GetString("file")

	if scriptName != "" {
		o.Script = NewObject(nil)
		o.Script.Name = o.Name + "_script"

		o.Script.FileName = scriptName
		o.Script.NewScript()
		Objects = append(Objects, o.Script)
	}

	if o.AutoStart {
		o.Trigger(o, nil)
	}

	o.Draw = func(o *Object) {
		if !DebugMode {
			return
		}

		rl.DrawCircle(int32(o.Position.X), int32(o.Position.Y), 5, rl.White)
		drawTextCentered(fmt.Sprintf("%s", o.Name), int32(o.Position.X), int32(o.Position.Y)+5, 10, rl.White)
		drawTextCentered(fmt.Sprintf("Remaining: %d ms", o.Remaining), int32(o.Position.X), int32(o.Position.Y)+15, 10, rl.White)
	}
}

func updateWait(o *Object, dt float32) {
	if !o.Started {
		return
	}

	if o.Remaining < 0 {
		o.Started = false
		o.Remaining = 0

		if o.Target != nil {
			o.Target.Trigger(o.Target, o)
		}

		if o.Script != nil {
			o.Script.Trigger(o.Script, o)
		}

		if o.Target == nil && o.Script == nil {
			log.Printf("Timer %s has no target attached!\n", o.Name)
		}

		return
	}

	o.Remaining -= int(dt * 1000)
}

func triggerWait(o, inst *Object) {
	o.Remaining = o.Duration
	o.Started = true
}
