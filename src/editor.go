/*
 * @Author: V4 Games
 * @Date: 2018-11-09 02:36:36
 * @Last Modified by:   Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-11-09 02:36:36
 */

package main

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	rootIsCollapsed = false
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

	if element.isCollapsed != nil && isMouseInRectangle(offsetX, offsetY, rl.MeasureText(element.text, 10), 10) {
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
