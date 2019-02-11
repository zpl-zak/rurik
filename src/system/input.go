/*
   Copyright 2019 V4 Games

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
)

// InputAction defines an action tagged by input
type InputAction struct {
	positiveKeys []int32
	negativeKeys []int32
	allKeys      []int32
	joyAxis      int32
	joyButtons   []int32
}

// BindInputAction registers a new input action used by the game
func BindInputAction(name string, action InputAction) {
	if action.allKeys == nil {
		if action.positiveKeys == nil {
			action.positiveKeys = []int32{}
		}

		if action.negativeKeys == nil {
			action.negativeKeys = []int32{}
		}

		action.allKeys = append(action.positiveKeys, action.negativeKeys...)
	}

	keybindings[name] = action
}

// InitInput initializes the input system
func InitInput() {
	GamepadID = 0

	BindInputAction("horizontal", InputAction{
		positiveKeys: []int32{rl.KeyD, rl.KeyRight},
		negativeKeys: []int32{rl.KeyA, rl.KeyLeft},
		joyAxis:      rl.GamepadXboxAxisLeftX,
	})

	BindInputAction("vertical", InputAction{
		positiveKeys: []int32{rl.KeyS, rl.KeyDown},
		negativeKeys: []int32{rl.KeyW, rl.KeyUp},
		joyAxis:      rl.GamepadXboxAxisLeftY,
	})

	BindInputAction("up", InputAction{
		allKeys:    []int32{rl.KeyW, rl.KeyUp},
		joyButtons: []int32{rl.GamepadXboxButtonUp},
	})

	BindInputAction("down", InputAction{
		allKeys:    []int32{rl.KeyS, rl.KeyDown},
		joyButtons: []int32{rl.GamepadXboxButtonDown},
	})

	BindInputAction("use", InputAction{
		allKeys:    []int32{rl.KeyE, rl.KeyEnter},
		joyButtons: []int32{rl.GamepadXboxButtonA},
	})
}

// IsKeyDown checks whether the key is down
func IsKeyDown(action string) bool {
	for _, v := range keybindings[action].allKeys {
		if rl.IsKeyDown(v) {
			return true
		}
	}

	for _, v := range keybindings[action].joyButtons {
		if rl.IsGamepadButtonDown(GamepadID, v) {
			return true
		}
	}

	return false
}

// IsKeyPressed checks whether the key is pressed
func IsKeyPressed(action string) bool {
	for _, v := range keybindings[action].allKeys {
		if rl.IsKeyPressed(v) {
			return true
		}
	}

	for _, v := range keybindings[action].joyButtons {
		if rl.IsGamepadButtonPressed(GamepadID, v) {
			return true
		}
	}

	return false
}

// IsKeyReleased checks whether the key is released
func IsKeyReleased(action string) bool {
	for _, v := range keybindings[action].allKeys {
		if rl.IsKeyReleased(v) {
			return true
		}
	}

	for _, v := range keybindings[action].joyButtons {
		if rl.IsGamepadButtonReleased(GamepadID, v) {
			return true
		}
	}

	return false
}

// GetAxis returns the axis value of an input
func GetAxis(action string) (rate float32) {
	a := keybindings[action]

	rate = rl.GetGamepadAxisMovement(GamepadID, a.joyAxis)

	if math.Abs(float64(rate)) < GamepadDeadZone {
		rate = 0
	}

	for _, v := range a.positiveKeys {
		if rl.IsKeyDown(v) {
			rate = 1
		}
	}

	for _, v := range a.negativeKeys {
		if rl.IsKeyDown(v) {
			rate = -1
		}
	}

	return
}
