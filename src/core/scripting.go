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

	// ScriptingContext is a scripting ScriptingContext
	ScriptingContext *otto.Otto
)

func initDefaultEvents() {
	RegisterEvent("exitGame", func(in InvokeData) InvokeData {
		CloseGame()
		return nil
	})

	RegisterEvent("followPlayer", func(in InvokeData) InvokeData {
		var data struct{ Speed float64 }
		DecodeInvokeData(&data, in)

		if data.Speed != 0 {
			MainCamera.Speed = float32(data.Speed)
		}

		MainCamera.Mode = CameraModeFollow
		MainCamera.Follow = LocalPlayer
		return nil
	})

	RegisterEvent("cameraInterpolate", func(in InvokeData) InvokeData {
		var data struct {
			Speed   float64
			Start   string
			End     string
			Instant bool
		}
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
func DecodeInvokeData(data interface{}, in InvokeData) {
	inp := in.(map[string]interface{})
	ref := reflect.ValueOf(data).Elem()
	dataName := reflect.TypeOf(data).Elem().Name()

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
					dataName,
				)
			}
		} else {
			log.Printf("Property %s not found inside of %s while invoking an event!\n", k, dataName)
		}
	}
}

func initScriptingSystem() {
	initDefaultEvents()

	ScriptingContext = otto.New()

	ScriptingContext.Set("log", func(call otto.FunctionCall) otto.Value {
		obj := call.Argument(0)
		fmt.Println(obj)

		return otto.Value{}
	})

	ScriptingContext.Set("findObject", func(call otto.FunctionCall) otto.Value {
		arg, _ := call.Argument(0).ToString()
		wv, _ := ScriptingContext.Get("CurrentWorld")
		w, _ := wv.Export()
		obj, _ := w.(*World).FindObject(arg)
		ret, _ := ScriptingContext.ToValue(obj)
		return ret
	})

	ScriptingContext.Set("setProperty", func(call otto.FunctionCall) otto.Value {
		source, _ := call.Argument(0).Export()
		field, _ := call.Argument(1).ToString()
		value, _ := call.Argument(2).Export()
		v := reflect.ValueOf(source)
		vd := reflect.ValueOf(value)
		r := reflect.Indirect(v).FieldByName(field)
		r.Set(vd)

		return otto.Value{}
	})

	ScriptingContext.Set("exitGame", func(call otto.FunctionCall) otto.Value {
		CloseGame()
		return otto.Value{}
	})

	ScriptingContext.Set("invoke", func(call otto.FunctionCall) otto.Value {
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
		tryConv, done := otto.ToValue(res)
		if done == nil {
			return tryConv
		}

		retObj, err := jsoniter.MarshalToString(&res)
		if err != nil {
			log.Printf("Invalid invoke return value! %v\n", err)
		}

		ret, _ := ScriptingContext.Object(fmt.Sprintf("(%s)", retObj))

		return ret.Value()
	})

	ScriptingContext.Object("exports = {}")
}

// RegisterEvent registers a particular event
func RegisterEvent(name string, call func(data InvokeData) InvokeData) {
	EventHandlers[name] = call
}
