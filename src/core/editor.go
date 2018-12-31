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

package core

import (
	rl "github.com/zaklaus/raylib-go/raylib"
)

var (
	rootIsCollapsed = true
	rootElement     = &editorElement{
		text:        "editor",
		isCollapsed: &rootIsCollapsed,
		children:    []*editorElement{},
	}
)

type editorElement struct {
	text        string
	isCollapsed *bool
	children    []*editorElement
}

func pushEditorElement(element *editorElement, text string, isCollapsed *bool) *editorElement {
	child := &editorElement{
		text:        text,
		isCollapsed: isCollapsed,
		children:    []*editorElement{},
	}
	element.children = append(element.children, child)

	return child
}

// DrawEditor draws debug UI
func DrawEditor() {
	if DebugMode {
		drawEditorElement(rootElement, 5, 5)
	}
}

// UpdateEditor updates editor debug UI
func UpdateEditor() {
	flushEditorElement()
}

func drawEditorElement(element *editorElement, offsetX, offsetY int32) int32 {
	color := rl.White
	var ext int32 = 10

	if element.isCollapsed != nil && IsMouseInRectangle(offsetX, offsetY, rl.MeasureText(element.text, 10), 10) {
		color = rl.Red

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			*element.isCollapsed = !*element.isCollapsed
		}
	}

	rl.DrawText(element.text, offsetX+1, offsetY+1, 10, rl.Black)
	rl.DrawText(element.text, offsetX, offsetY, 10, color)

	if element.isCollapsed != nil && *element.isCollapsed {
		return ext
	}

	for _, v := range element.children {
		ext += drawEditorElement(v, offsetX+5, offsetY+ext)
	}

	return ext
}

func flushEditorElement() {
	rootElement = &editorElement{
		text:        "editor",
		isCollapsed: &rootIsCollapsed,
		children:    []*editorElement{},
	}
}
