package main

import (
	"fmt"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/robertkrimen/otto"
)

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
		rl.CloseWindow()
		return otto.Value{}
	})

	vm.Set("CurrentWorld", o.World)
	vm.Set("CurrentMap", CurrentMap)
	vm.Set("Self", o)
	vm.Set("LocalPlayer", LocalPlayer)
	vm.Set("MainCamera", MainCamera)
}
