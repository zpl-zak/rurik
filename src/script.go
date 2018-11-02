package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/robertkrimen/otto"
)

type script struct {
	File   string
	Ctx    *otto.Otto
	Source string
}

// NewScript instance
func (o *Object) NewScript() {
	if o.File == "" {
		o.File = o.Meta.Properties.GetString("file")
	}

	src, err := ioutil.ReadFile(fmt.Sprintf("assets/map/%s/scripts/%s", mapName, o.File))

	if err != nil {
		log.Fatalf("Script object's %s file %s was not found!\n", o.Name, o.File)
		return
	}

	o.Source = string(src)

	o.Ctx = otto.New()
	initGameAPI(o.Ctx)

	o.Trigger = func(o, inst *Object) {
		_, err := o.Ctx.Eval(o.Source)

		if err != nil {
			log.Fatalf("Script error detected at '%s':%s: \n\t%s!\n", o.Name, o.File, err.Error())
			return
		}
	}
}
