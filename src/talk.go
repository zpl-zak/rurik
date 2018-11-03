package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	mouseDoublePress = 500
)

type talk struct {
	Texts                *Dialogue
	currentText          *Dialogue
	selectedChoice       int
	wasPrevLocked        bool
	mouseDoublePressTime int32
}

var (
	chatBanner *rl.Texture2D
)

// NewTalk instance
func (o *Object) NewTalk() {
	if chatBanner == nil {
		img := GetTexture("assets/gfx/chat_banner.png")
		chatBanner = &img
	}

	o.Update = func(o *Object, dt float32) {
		if !o.Started {
			return
		}

		if o.mouseDoublePressTime > 0 {
			o.mouseDoublePressTime -= int32(dt * 1000)
		} else if o.mouseDoublePressTime < 0 {
			o.mouseDoublePressTime = 0
		}

		if len(o.currentText.Choices) > 0 {
			if IsKeyPressed("up") {
				o.selectedChoice--

				if o.selectedChoice < 0 {
					o.selectedChoice = len(o.currentText.Choices) - 1
				}
			}

			if IsKeyPressed("down") {
				o.selectedChoice++

				if o.selectedChoice >= len(o.currentText.Choices) {
					o.selectedChoice = 0
				}
			}
		}

		if IsKeyPressed("use") || (rl.IsMouseButtonReleased(rl.MouseLeftButton) && o.mouseDoublePressTime > 0) {

			if o.mouseDoublePressTime > 0 {
				o.mouseDoublePressTime = 0
			}

			tgt, _ := FindObject(o.currentText.Target)

			if len(o.currentText.Choices) > 0 {
				o.currentText = o.currentText.Choices[o.selectedChoice].Next
			} else {
				o.currentText = o.currentText.Next
			}

			if o.currentText == nil {
				o.Started = false
				LocalPlayer.Locked = o.wasPrevLocked
			}

			if tgt != nil {
				tgt.Trigger(tgt, o)
			}
		}
	}

	o.Trigger = func(o, inst *Object) {
		data, err := ioutil.ReadFile(fmt.Sprintf("assets/map/%s/texts/%s", mapName, o.FileName))

		if err != nil {
			log.Fatalf("Could not load texts for %s [%s]: %s\n", o.Name, o.FileName, err.Error())
		}

		err = json.Unmarshal(data, &o.Texts)

		if err != nil {
			log.Fatalf("Error loading text %s for %s: %s\n", o.FileName, o.Name, err.Error())
			return
		}

		o.currentText = o.Texts

		initText(o.currentText)

		o.Started = true
		o.wasPrevLocked = LocalPlayer.Locked
		LocalPlayer.Locked = true
	}

	o.DrawUI = func(o *Object) {
		if !o.Started {
			return
		}

		var height int32 = 120
		var width int32 = 640
		start := ScreenHeight - height
		offsetX := (ScreenWidth - width) / 2

		rl.DrawRectangle(int32(offsetX), start, width, height, rl.Black)
		rl.DrawTexture(*chatBanner, int32(offsetX), start, rl.White)
		ot := o.currentText

		// Pos X: 5, Y: 5
		// Scale W: 34, 35
		if ot.AvatarFile != "" {
			rl.DrawTexturePro(
				ot.Avatar,
				rl.NewRectangle(0, 0, float32(ot.Avatar.Width), float32(ot.Avatar.Height)),
				rl.NewRectangle(float32(offsetX)+5, float32(start)+5, 32, 32),
				rl.Vector2{},
				0,
				rl.White,
			)
		}

		rl.DrawText(
			ot.Name,
			offsetX+45,
			start+16,
			10,
			rl.Orange,
		)

		rl.DrawText(
			ot.Text,
			offsetX+5,
			start+45,
			10,
			rl.White,
		)

		// choices
		chsX := width - 220
		chsY := start + 16

		if len(ot.Choices) > 0 {
			for idx, ch := range ot.Choices {
				ypos := chsY + int32(idx)*15 - 2
				if idx == o.selectedChoice {
					rl.DrawRectangle(chsX, ypos, 200, 15, rl.DarkPurple)
				}

				rl.DrawText(
					fmt.Sprintf("%d. %s", idx+1, ch.Text),
					chsX+5,
					chsY+int32(idx)*15,
					10,
					rl.White,
				)

				if isMouseInRectangle(chsX, ypos, 200, 15) {
					if rl.IsMouseButtonDown(rl.MouseLeftButton) {
						rl.DrawRectangleLines(chsX, ypos, 200, 15, rl.Pink)
					} else {
						rl.DrawRectangleLines(chsX, ypos, 200, 15, rl.Purple)
					}

					if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
						o.selectedChoice = idx

						o.mouseDoublePressTime = mouseDoublePress
					}
				}
			}
		} else {
			rl.DrawRectangle(chsX, chsY-2, 200, 15, rl.DarkPurple)
			rl.DrawText(
				"Press E to continue...",
				chsX+5,
				chsY,
				10,
				rl.White,
			)
		}
	}

	if o.AutoStart {
		o.Trigger(o, nil)
	}
}
