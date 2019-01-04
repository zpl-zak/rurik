package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"

	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

// This file showcases how easy it is to override class's behavior and implement your own logic for a core feature.
// Instead of relying on Styling structures, that on their own are quite limited, this approach
// offers full freedom to creativity at the expense of having duplicate code.

// NewTalkDemo class
func NewTalkDemo(o *core.Object) {
	o.NewTalk()

	o.DrawUI = func(o *core.Object) {
		if !o.Started {
			return
		}

		width := system.ScreenWidth
		start := system.ScreenHeight / 2
		ot := o.CurrentText

		rl.DrawText(ot.Text, 15, 30, 10, rl.RayWhite)

		// choices
		chsX := width / 2
		chsY := start + 40

		if len(ot.Choices) > 0 {
			for idx, ch := range ot.Choices {
				ypos := chsY + int32(idx)*15 - 2
				if idx == o.SelectedChoice {
					rl.DrawRectangle(chsX-100, ypos, 200, 15, rl.DarkPurple)
				}

				core.DrawTextCentered(
					ch.Text,
					chsX,
					chsY+int32(idx)*15,
					10,
					rl.White,
				)

				if core.IsMouseInRectangle(chsX, ypos, 200, 15) {
					if rl.IsMouseButtonDown(rl.MouseLeftButton) {
						rl.DrawRectangleLines(chsX-100, ypos, 200, 15, rl.Pink)
					} else {
						rl.DrawRectangleLines(chsX-100, ypos, 200, 15, rl.Purple)
					}

					if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
						o.SelectedChoice = idx

						o.MouseDoublePressTime = core.MouseDoublePress
					}
				}
			}
		} else {
			rl.DrawRectangle(chsX-100, chsY-2, 200, 15, rl.DarkPurple)
			core.DrawTextCentered(
				"Press E to continue...",
				chsX,
				chsY,
				10,
				rl.White,
			)
		}
	}
}