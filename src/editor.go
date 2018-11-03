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
		flushEditorElement()
	}
}

func drawEditorElement(element *editorElement, offsetX, offsetY int32) {
	color := rl.White

	if element.isCollapsed != nil && isMouseInRectangle(offsetX, offsetY, rl.MeasureText(element.text, 10), 10) {
		color = rl.Red

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			*element.isCollapsed = !*element.isCollapsed
		}
	}

	rl.DrawText(element.text, offsetX, offsetY, 10, color)

	if element.isCollapsed != nil && *element.isCollapsed {
		return
	}

	for i, v := range element.children {
		drawEditorElement(v, offsetX+5, offsetY+int32(10*(i+1)))
	}
}

func flushEditorElement() {
	rootElement = &editorElement{
		text:        "editor",
		isCollapsed: &rootIsCollapsed,
		children:    []*editorElement{},
	}
}
