package main

import (
	"encoding/gob"
	"fmt"

	"github.com/zaklaus/rurik/src/core"
)

type demoClassData struct {
	Foo int
	Bar string
}

func (d *demoClassData) Serialize(enc *gob.Encoder) {
	enc.Encode(d)
}

func (d *demoClassData) Deserialize(dec *gob.Decoder) {
	dec.Decode(&d)
}

// NewTestClass is a custom type
func NewTestClass(o *core.Object) {
	fmt.Printf("Initializing custom type from demo!\n")

	o.AutoStart = true
	o.AnimTag = "Base"
	o.FileName = "ball"
	o.NewAnim()
	o.IsCollidable = false
	o.UserData = &demoClassData{
		Foo: 42,
		Bar: "Hi, sailor!",
	}

	demoScript := o.GetWorld().NewObjectPro("demo_loop_event", "script")
	demoScript.FileName = "eventDemo.js"
	demoScript.IsPersistent = false
	o.GetWorld().FinalizeObject(demoScript)

	demoScript.Trigger(demoScript, o)
}
