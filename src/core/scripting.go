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

var (
	// EventHandlers consists of registered events you can invoke from the scripting side
	EventHandlers = make(map[string]func(data string) string)
)

func initDefaultEvents() {
	RegisterEvent("exitGame", func(in string) string {
		CloseGame()
		return "{}"
	})

	RegisterEvent("followPlayer", func(in string) string {
		type followPlayerData struct {
			Speed float32
		}

		var data followPlayerData
		jsoniter.UnmarshalFromString(in, &data)

		if data.Speed != 0 {
			MainCamera.Speed = data.Speed
		}

		MainCamera.Mode = CameraModeFollow
		MainCamera.Follow = LocalPlayer

		return "{}"
	})

	RegisterEvent("cameraInterpolate", func(in string) string {
		type cameraInterpolateData struct {
			Speed   float32
			Start   string
			End     string
			Instant bool
		}

		var data cameraInterpolateData
		jsoniter.UnmarshalFromString(in, &data)

		if data.Speed != 0 {
			MainCamera.Speed = data.Speed
		}

		if data.Instant {
			MainCamera.First = true
		}

		MainCamera.Mode = CameraModeLerp
		MainCamera.Start, _ = CurrentMap.World.FindObject(data.Start)
		MainCamera.End, _ = CurrentMap.World.FindObject(data.End)

		return "{}"
	})
}

func initGameAPI(o *Object, vm *otto.Otto) {
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

		eventData := ""

		if len(call.ArgumentList) > 1 {
			eventData, _ = call.Argument(1).ToString()
		}

		ret, _ := otto.ToValue(event(eventData))
		return ret
	})

	vm.Set("CurrentWorld", o.world)
	vm.Set("CurrentMap", CurrentMap)
	vm.Set("Self", o)
	vm.Set("LocalPlayer", LocalPlayer)
	vm.Set("MainCamera", MainCamera)
}

// RegisterEvent registers a particular event
func RegisterEvent(name string, call func(data string) string) {
	EventHandlers[name] = call
}
