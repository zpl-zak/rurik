package core

import (
	ry "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawEditor draws debug UI
func DrawEditor() {
	ry.Button(rl.NewRectangle(5, 5, 10, 10), "...")
}
