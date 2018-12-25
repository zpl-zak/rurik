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
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/robertkrimen/otto"
)

// InvokeData is the incoming data from the DSL caller
type InvokeData interface{}

var (
	// EventHandlers consists of registered events you can invoke from the scripting side
	EventHandlers = make(map[string]func(data InvokeData) InvokeData)

	// ScriptingContext is a scripting vm
	ScriptingContext *otto.Otto
)

func initDefaultEvents() {
	RegisterEvent("exitGame", func(in InvokeData) InvokeData {
		CloseGame()
		return nil
	})

	RegisterEvent("followPlayer", func(in InvokeData) InvokeData {
		type followPlayerData struct {
			Speed float64
		}

		var data followPlayerData
		DecodeInvokeData(&data, in)

		if data.Speed != 0 {
			MainCamera.Speed = float32(data.Speed)
		}

		MainCamera.Mode = CameraModeFollow
		MainCamera.Follow = LocalPlayer
		return nil
	})

	RegisterEvent("cameraInterpolate", func(in InvokeData) InvokeData {
		type cameraInterpolateData struct {
			Speed   float64
			Start   string
			End     string
			Instant bool
		}

		var data cameraInterpolateData
		DecodeInvokeData(&data, in)

		if data.Speed != 0 {
			MainCamera.Speed = float32(data.Speed)
		}

		if data.Instant {
			MainCamera.First = true
		}

		MainCamera.Mode = CameraModeLerp
		MainCamera.Start, _ = CurrentMap.World.FindObject(data.Start)
		MainCamera.End, _ = CurrentMap.World.FindObject(data.End)

		return nil
	})

	RegisterEvent("testReturnValue", func(in InvokeData) InvokeData {
		return struct {
			Foo string
			Bar int32
		}{
			"hey",
			123,
		}
	})
}

// DecodeInvokeData decodes incoming data from the script
func DecodeInvokeData(data InvokeData, in InvokeData) {
	inp := in.(map[string]interface{})
	ref := reflect.ValueOf(data).Elem()

	for k, v := range inp {
		fieldSource := reflect.ValueOf(v)
		fieldDest := ref.FieldByName(k)

		if fieldDest.IsValid() && fieldDest.CanSet() {
			if fieldDest.Type() == fieldSource.Type() {
				fieldDest.Set(fieldSource)
			} else {
				log.Printf(
					"Property %s of type %s has a different type %s inside of %s !\n",
					k,
					fieldSource.Type().String(),
					fieldDest.Type().String(),
					reflect.TypeOf(data).Elem().Name(),
				)
			}
		} else {
			log.Printf("Property %s not found while invoking an event!\n", k)
		}
	}
}

func initScriptingSystem() {
	initDefaultEvents()

	ScriptingContext = otto.New()

	initGameAPI(ScriptingContext)
}

func initGameAPI(vm *otto.Otto) {
	vm.Set("log", func(call otto.FunctionCall) otto.Value {
		obj := call.Argument(0)
		fmt.Println(obj)

		return otto.Value{}
	})

	vm.Set("findObject", func(call otto.FunctionCall) otto.Value {
		arg, _ := call.Argument(0).ToString()
		wv, _ := vm.Get("CurrentWorld")
		w, _ := wv.Export()
		obj, _ := w.(*World).FindObject(arg)
		ret, _ := vm.ToValue(obj)
		return ret
	})

	vm.Set("setProperty", func(call otto.FunctionCall) otto.Value {
		source, _ := call.Argument(0).Export()
		field, _ := call.Argument(1).ToString()
		value, _ := call.Argument(2).Export()
		v := reflect.ValueOf(source)
		vd := reflect.ValueOf(value)
		r := reflect.Indirect(v).FieldByName(field)
		r.Set(vd)

		return otto.Value{}
	})

	vm.Set("exitGame", func(call otto.FunctionCall) otto.Value {
		CloseGame()
		return otto.Value{}
	})

	vm.Set("invoke", func(call otto.FunctionCall) otto.Value {
		eventName, _ := call.Argument(0).ToString()

		event, ok := EventHandlers[eventName]

		if !ok {
			log.Printf("Can't invoke event '%s'!\n", eventName)
			return otto.Value{}
		}

		var eventData InvokeData

		if len(call.ArgumentList) > 1 {
			eventData, _ = call.Argument(1).Export()
		}

		res := event(eventData)
		retObj, err := jsoniter.MarshalToString(&res)
		if err != nil {
			log.Printf("Invalid invoke return value! %v\n", err)
		}

		if retObj == "null" {
			return otto.Value{}
		}

		ret, _ := ScriptingContext.Object(fmt.Sprintf("(%s)", retObj))

		return ret.Value()
	})

	vm.Eval("var exports = {};")
}

// RegisterEvent registers a particular event
func RegisterEvent(name string, call func(data InvokeData) InvokeData) {
	EventHandlers[name] = call
}
