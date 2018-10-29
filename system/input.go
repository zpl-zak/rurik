package system

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	keybindings = make(map[string][]int32)
)

// InitInput initializes the input system
func InitInput() {
	keybindings["up"] = []int32{rl.KeyW, rl.KeyUp}
	keybindings["down"] = []int32{rl.KeyS, rl.KeyDown}
	keybindings["left"] = []int32{rl.KeyA, rl.KeyLeft}
	keybindings["right"] = []int32{rl.KeyD, rl.KeyRight}
	keybindings["use"] = []int32{rl.KeyE}
	keybindings["menu"] = []int32{rl.KeyEnter}
	keybindings["exit"] = []int32{rl.KeyEscape}
}

// IsKeyDown checks whether the key is down
func IsKeyDown(action string) bool {
	for _, v := range keybindings[action] {
		if rl.IsKeyDown(v) {
			return true
		}
	}

	return false
}

// IsKeyPressed checks whether the key is pressed
func IsKeyPressed(action string) bool {
	for _, v := range keybindings[action] {
		if rl.IsKeyPressed(v) {
			return true
		}
	}

	return false
}

// IsKeyReleased checks whether the key is released
func IsKeyReleased(action string) bool {
	for _, v := range keybindings[action] {
		if rl.IsKeyReleased(v) {
			return true
		}
	}

	return false
}
