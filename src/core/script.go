/*
 * @Author: V4 Games
 * @Date: 2018-11-09 02:36:51
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-08 21:01:10
 */

package core

import (
	"log"

	"github.com/json-iterator/go"
	"github.com/robertkrimen/otto"
)

type script struct {
	Ctx         *otto.Otto
	WasExecuted bool
	CanRepeat   bool
	Source      string
}

type scriptData struct {
	objectData
	WasExecuted bool `json:"done"`
	CanRepeat   bool `json:"rep"`
}

// NewScript sequence script
func (o *Object) NewScript() {
	o.Trigger = func(o, inst *Object) {
		if o.FileName == "" {
			o.FileName = o.Meta.Properties.GetString("file")
		}

		o.Source = GetScript(o.FileName)

		o.Ctx = otto.New()
		initGameAPI(o, o.Ctx)

		log.Printf("Loading script %s...\n", o.FileName)

		if !o.WasExecuted || o.CanRepeat {
			_, err := o.Ctx.Eval(o.Source)

			if err != nil {
				log.Fatalf("Script error detected at '%s':%s: \n\t%s!\n", o.Name, o.FileName, err.Error())
				return
			}
		}

		o.WasExecuted = true
	}

	o.Serialize = func(o *Object) string {
		val, _ := jsoniter.MarshalToString(&scriptData{
			WasExecuted: o.WasExecuted,
			CanRepeat:   o.CanRepeat,
		})

		return val
	}

	o.Deserialize = func(o *Object, v string) {
		var dat scriptData
		jsoniter.UnmarshalFromString(v, &dat)
		o.WasExecuted = dat.WasExecuted
		o.CanRepeat = dat.CanRepeat
	}
}
