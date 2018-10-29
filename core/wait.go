package core

import (
	"fmt"
	"log"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type wait struct {
	Duration  int
	Remaining int
}

// NewWait instance
func NewWait(o *Object) {
	o.Duration, _ = strconv.Atoi(o.Meta.Properties.GetString("duration"))
	autostart, _ := strconv.Atoi(o.Meta.Properties.GetString("autostart"))
	o.Trigger = triggerWait
	o.Update = updateWait

	scriptName := o.Meta.Properties.GetString("file")

	if scriptName != "" {
		scr := NewObject(nil)
		scr.Name = o.Name + "_script"

		scr.File = scriptName
		NewScript(scr)
		Objects = append(Objects, scr)
	}

	if autostart == 1 {
		o.Trigger(o, nil)
	}

	o.Draw = func(o *Object) {
		if !DebugMode {
			return
		}

		rl.DrawCircle(int32(o.Position.X), int32(o.Position.Y), 5, rl.White)
		drawTextCentered(fmt.Sprintf("%s\nRemaining: %d ms", o.Name, o.Remaining), int32(o.Position.X), int32(o.Position.Y)+5, 10, rl.White)
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
		} else {
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
