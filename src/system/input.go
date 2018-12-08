/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:27:26
 * @Last Modified by:   Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-11-14 02:27:26
 */

package system

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	keybindings = make(map[string]inputAction)

	// GamepadDeadZone movement threshold
	GamepadDeadZone = 0.25

	// GamepadID represents the controller index
	GamepadID int32
)

type inputAction struct {
	positiveKeys []int32
	negativeKeys []int32
	allKeys      []int32
	joyAxis      int32
	joyButtons   []int32
}

// BindInputAction registers a new input action used by the game
func BindInputAction(name string, action inputAction) {
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

	BindInputAction("horizontal", inputAction{
		positiveKeys: []int32{rl.KeyD, rl.KeyRight},
		negativeKeys: []int32{rl.KeyA, rl.KeyLeft},
		joyAxis:      rl.GamepadXboxAxisLeftX,
	})

	BindInputAction("vertical", inputAction{
		positiveKeys: []int32{rl.KeyS, rl.KeyDown},
		negativeKeys: []int32{rl.KeyW, rl.KeyUp},
		joyAxis:      rl.GamepadXboxAxisLeftY,
	})

	BindInputAction("up", inputAction{
		allKeys:    []int32{rl.KeyW, rl.KeyUp},
		joyButtons: []int32{rl.GamepadXboxButtonUp},
	})

	BindInputAction("down", inputAction{
		allKeys:    []int32{rl.KeyS, rl.KeyDown},
		joyButtons: []int32{rl.GamepadXboxButtonDown},
	})

	BindInputAction("use", inputAction{
		allKeys:    []int32{rl.KeyE},
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
