package main

import (
	"fmt"

	"madaraszd.net/zaklaus/rurik/src/core"
)

// NewTestClass is a custom type
func NewTestClass(o *core.Object) {
	fmt.Printf("Initializing custom type from demo!\n")

	o.AutoStart = true
	o.AnimTag = "Base"
	o.FileName = "ball"
	o.NewAnim()
	o.IsCollidable = false
}
