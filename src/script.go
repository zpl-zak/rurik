package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/robertkrimen/otto"
)

type script struct {
	Ctx    *otto.Otto
	Source string
}

// NewScript sequence script
func (o *Object) NewScript() {
	if o.FileName == "" {
		o.FileName = o.Meta.Properties.GetString("file")
	}

	src, err := ioutil.ReadFile(fmt.Sprintf("assets/map/%s/scripts/%s", CurrentMap.mapName, o.FileName))

	if err != nil {
		log.Fatalf("Script object's %s file %s was not found!\n", o.Name, o.FileName)
		return
	}

	o.Source = string(src)

	o.Ctx = otto.New()
	initGameAPI(o, o.Ctx)

	log.Printf("Loading script %s...\n", o.FileName)

	o.Trigger = func(o, inst *Object) {
		_, err := o.Ctx.Eval(o.Source)

		if err != nil {
			log.Fatalf("Script error detected at '%s':%s: \n\t%s!\n", o.Name, o.FileName, err.Error())
			return
		}
	}
}
