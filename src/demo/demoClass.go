package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/zaklaus/rurik/src/core"
)

type demoClassData struct {
	Foo int
	Bar string
}

func (d *demoClassData) Serialize() string {
	r, _ := jsoniter.MarshalToString(*d)
	return r
}

func (d *demoClassData) Deserialize(input string) {
	jsoniter.UnmarshalFromString(input, d)
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
