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

	jsoniter "github.com/json-iterator/go"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

const (
	mouseDoublePress = 500
)

type talk struct {
	Texts                *Dialogue
	currentText          *Dialogue
	selectedChoice       int
	mouseDoublePressTime int32
}

type talkData struct {
	objectData
	WasExecuted bool `json:"done"`
	CanRepeat   bool `json:"rep"`
}

var (
	chatBanner *rl.Texture2D
)

// NewTalk dialogue handler
func (o *Object) NewTalk() {
	if chatBanner == nil {
		img := system.GetTexture("chat_banner.png")
		chatBanner = img
	}

	o.Serialize = func(o *Object) string {
		val, _ := jsoniter.MarshalToString(&scriptData{
			WasExecuted: o.WasExecuted,
			CanRepeat:   o.CanRepeat,
		})

		return val
	}

	o.Deserialize = func(o *Object, v string) {
		var dat talkData
		jsoniter.UnmarshalFromString(v, &dat)
		o.WasExecuted = dat.WasExecuted
		o.CanRepeat = dat.CanRepeat
	}

	o.Update = func(o *Object, dt float32) {
		if !o.Started {
			return
		}

		CanSave = bitsSet(CanSave, isInDialogue)

		if o.LastTrigger == 0 {
			o.LastTrigger = rl.GetTime()
			return
		}

		if o.mouseDoublePressTime > 0 {
			o.mouseDoublePressTime -= int32(1000 * dt)
		} else if o.mouseDoublePressTime < 0 {
			o.mouseDoublePressTime = 0
		}

		if len(o.currentText.Choices) > 0 {
			if system.IsKeyPressed("up") {
				o.selectedChoice--

				if o.selectedChoice < 0 {
					o.selectedChoice = len(o.currentText.Choices) - 1
				}
			}

			if system.IsKeyPressed("down") {
				o.selectedChoice++

				if o.selectedChoice >= len(o.currentText.Choices) {
					o.selectedChoice = 0
				}
			}
		}

		if system.IsKeyPressed("use") || (rl.IsMouseButtonReleased(rl.MouseLeftButton) && o.mouseDoublePressTime > 0) {
			if o.mouseDoublePressTime > 0 {
				o.mouseDoublePressTime = 0
			}

			tgt, _ := o.world.FindObject(o.currentText.Target)

			if len(o.currentText.Choices) > 0 {
				o.currentText = o.currentText.Choices[o.selectedChoice].Next
			} else {
				o.currentText = o.currentText.Next
			}

			if o.currentText == nil {
				o.Started = false
				CanSave = bitsClear(CanSave, isInDialogue)
				o.WasExecuted = true
			}

			if tgt != nil {
				tgt.Trigger(tgt, o)
			}
		}

		o.LastTrigger = rl.GetTime()
	}

	o.Trigger = func(o, inst *Object) {
		if o.WasExecuted && !o.CanRepeat {
			return
		}

		if o.Texts == nil {
			data := GetDialogue(o.FileName)
			o.Texts = data
		}

		o.currentText = o.Texts

		InitText(o.currentText)

		o.Started = true
		o.LastTrigger = 0
	}

	o.DrawUI = func(o *Object) {
		if !o.Started {
			return
		}

		var height int32 = 120
		var width int32 = 640
		start := system.ScreenHeight - height
		offsetX := (system.ScreenWidth - width) / 2

		rl.DrawRectangle(int32(offsetX), start, width, height, rl.Black)
		rl.DrawTexture(*chatBanner, int32(offsetX), start, rl.White)
		ot := o.currentText

		// Pos X: 5, Y: 5
		// Scale W: 34, 35
		if ot.AvatarFile != "" {
			rl.DrawTexturePro(
				*ot.Avatar,
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

	o.Init = func(o *Object) {
		if o.AutoStart {
			o.Trigger(o, nil)
		}
	}
}
