/*
   Copyright 2019 Dominik Madarász <zaklaus@madaraszd.net>

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

package system

import (
	"math"

	rl "github.com/zaklaus/raylib-go/raylib"
)

var (
	keybindings = make(map[string]InputAction)

	// GamepadDeadZone movement threshold
	GamepadDeadZone = 0.25

	// GamepadID represents the controller index
	GamepadID int32

	// MouseDelta contains last mouse movement
	MouseDelta [2]int32

	lastMousePosition [2]int32
)

// InputAction defines an action tagged by input
type InputAction struct {
	PositiveKeys []int32
	NegativeKeys []int32
	AllKeys      []int32
	JoyAxis      int32
	JoyButtons   []int32
}

// UpdateInput updates the user input
func UpdateInput() {
	rl.PollInputEvents()

	// Update MouseDelta
	newMousePosition := GetMousePosition()
	MouseDelta = [2]int32{
		newMousePosition[0] - lastMousePosition[0],
		newMousePosition[1] - lastMousePosition[1],
	}
	lastMousePosition = newMousePosition
}

// BindInputAction registers a new input action used by the game
func BindInputAction(name string, action InputAction) {
	if action.AllKeys == nil {
		if action.PositiveKeys == nil {
			action.PositiveKeys = []int32{}
		}

		if action.NegativeKeys == nil {
			action.NegativeKeys = []int32{}
		}

		action.AllKeys = append(action.PositiveKeys, action.NegativeKeys...)
	}

	keybindings[name] = action
}

// InitInput initializes the input system
func InitInput() {
	GamepadID = 0

	BindInputAction("horizontal", InputAction{
		PositiveKeys: []int32{rl.KeyD, rl.KeyRight},
		NegativeKeys: []int32{rl.KeyA, rl.KeyLeft},
		JoyAxis:      rl.GamepadXboxAxisLeftX,
	})

	BindInputAction("vertical", InputAction{
		PositiveKeys: []int32{rl.KeyS, rl.KeyDown},
		NegativeKeys: []int32{rl.KeyW, rl.KeyUp},
		JoyAxis:      rl.GamepadXboxAxisLeftY,
	})

	BindInputAction("up", InputAction{
		AllKeys:    []int32{rl.KeyW, rl.KeyUp},
		JoyButtons: []int32{rl.GamepadXboxButtonUp},
	})

	BindInputAction("down", InputAction{
		AllKeys:    []int32{rl.KeyS, rl.KeyDown},
		JoyButtons: []int32{rl.GamepadXboxButtonDown},
	})

	BindInputAction("use", InputAction{
		AllKeys:    []int32{rl.KeyE, rl.KeyEnter},
		JoyButtons: []int32{rl.GamepadXboxButtonA},
	})
}

// IsKeyDown checks whether the key is down
func IsKeyDown(action string) bool {
	for _, v := range keybindings[action].AllKeys {
		if rl.IsKeyDown(v) {
			return true
		}
	}

	for _, v := range keybindings[action].JoyButtons {
		if rl.IsGamepadButtonDown(GamepadID, v) {
			return true
		}
	}

	return false
}

// IsKeyPressed checks whether the key is pressed
func IsKeyPressed(action string) bool {
	for _, v := range keybindings[action].AllKeys {
		if rl.IsKeyPressed(v) {
			return true
		}
	}

	for _, v := range keybindings[action].JoyButtons {
		if rl.IsGamepadButtonPressed(GamepadID, v) {
			return true
		}
	}

	return false
}

// IsKeyReleased checks whether the key is released
func IsKeyReleased(action string) bool {
	for _, v := range keybindings[action].AllKeys {
		if rl.IsKeyReleased(v) {
			return true
		}
	}

	for _, v := range keybindings[action].JoyButtons {
		if rl.IsGamepadButtonReleased(GamepadID, v) {
			return true
		}
	}

	return false
}

// GetAxis returns the axis value of an input
func GetAxis(action string) (rate float32) {
	a := keybindings[action]

	rate = rl.GetGamepadAxisMovement(GamepadID, a.JoyAxis)

	if math.Abs(float64(rate)) < GamepadDeadZone {
		rate = 0
	}

	for _, v := range a.PositiveKeys {
		if rl.IsKeyDown(v) {
			rate = 1
		}
	}

	for _, v := range a.NegativeKeys {
		if rl.IsKeyDown(v) {
			rate = -1
		}
	}

	return
}

// GetMousePosition returns a fixed mouse position
func GetMousePosition() [2]int32 {
	mo := rl.GetMousePosition()
	m := [2]int32{
		int32(mo.X / ScaleRatio),
		int32(mo.Y / ScaleRatio),
	}

	return m
}
