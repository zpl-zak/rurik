package main

import (
	"fmt"

	"github.com/zaklaus/rurik/src/core"
)

// NewTestClass is a custom type
func NewTestClass(o *core.Object) {
	fmt.Printf("Initializing custom type from demo!\n")

	o.AutoStart = true
	o.AnimTag = "Base"
	o.FileName = "ball"
	o.NewAnim()
	o.IsCollidable = false

	demoScript := o.GetWorld().NewObjectPro("demo_loop_event", "script")
	demoScript.FileName = "eventDemo.js"
	demoScript.IsPersistent = false
	o.GetWorld().FinalizeObject(demoScript)

	demoScript.Trigger(demoScript, o)
}
