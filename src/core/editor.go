/*
   Copyright 2019 Dominik Madarász <zaklaus@madaraszd.net>

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
	"math"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

var (
	rootIsCollapsed = true
	rootElement     = &editorElement{
		text:        "editor",
		isCollapsed: &rootIsCollapsed,
		children:    []*editorElement{},
	}
)

const (
	defaultGraphHeight     = 50
	defaultGraphWidth      = 200
	defaultGraphDataMargin = 5
)

type editorElement struct {
	text        string
	isCollapsed *bool
	children    []*editorElement
	callback    func()

	// Graphs
	graphEnabled bool
	lineColor    rl.Color
	dataMargin   int32
	graphHeight  int32
	graphWidth   int32
	pointData    []float64
	useCurves    bool
	ValueSuffix  string

	// Buttons
	isButton         bool
	buttonColor      rl.Color
	buttonHoverColor rl.Color
	buttonPressColor rl.Color
}

func pushEditorElement(element *editorElement, text string, isCollapsed *bool) *editorElement {
	return pushEditorElementEx(element, text, isCollapsed, func() {})
}

func pushEditorElementEx(element *editorElement, text string, isCollapsed *bool, callback func()) *editorElement {
	child := &editorElement{
		text:        text,
		isCollapsed: isCollapsed,
		children:    []*editorElement{},
		callback:    callback,
	}
	element.children = append(element.children, child)

	return child
}

func setUpButton(element *editorElement, callback func()) {
	element.isButton = true
	element.buttonColor = rl.DarkPurple
	element.buttonHoverColor = rl.Purple
	element.buttonPressColor = rl.Pink
	element.callback = callback
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

	DebugShowAll = !rootIsCollapsed
}

func drawGraph(element *editorElement, offsetX, offsetY int32) int32 {
	var height, width int32

	height = element.graphHeight
	width = element.graphWidth

	rl.DrawRectangle(offsetX, offsetY, width, height, rl.NewColor(40, 40, 40, 140))

	if len(element.pointData) < 1 {
		return height + 5
	}

	// Value mapping
	var biggestValue float64
	var smallestValue = math.MaxFloat64
	var sum float64
	var avgValue float64
	var nodeCount int32

	var graphXTreshold int32
	actualGraphWidth := int32(len(element.pointData)) * element.dataMargin

	if actualGraphWidth > width {
		graphXTreshold = actualGraphWidth - width
	}

	for x, v := range element.pointData {
		if (int32(x) * element.dataMargin) < graphXTreshold {
			continue
		}

		if v > biggestValue {
			biggestValue = v
		}

		if v < smallestValue {
			smallestValue = v
		}

		sum += v
		nodeCount++
	}

	avgValue = (sum / float64(nodeCount))

	if smallestValue == biggestValue {
		biggestValue++
	}

	scaleY := float32(height) / float32(biggestValue-smallestValue)

	// Plotting
	var oldValue int32 = -1
	for x, v := range element.pointData {
		if (int32(x) * element.dataMargin) < graphXTreshold {
			continue
		}

		scaledValue := int32(float32(v-smallestValue) * float32(scaleY))
		rl.DrawCircle(offsetX-graphXTreshold+(int32(x)*element.dataMargin), offsetY+height-scaledValue, 1, element.lineColor)

		if oldValue != -1 {
			if element.useCurves {
				rl.DrawLineBezier(
					rl.NewVector2(
						float32(offsetX-graphXTreshold+(int32(x-1)*element.dataMargin)),
						float32(offsetY+height-oldValue),
					),
					rl.NewVector2(
						float32(offsetX-graphXTreshold+(int32(x)*element.dataMargin)),
						float32(offsetY+height-scaledValue),
					),
					1,
					element.lineColor,
				)
			} else {
				rl.DrawLine(
					offsetX-graphXTreshold+(int32(x-1)*element.dataMargin),
					offsetY+height-oldValue,
					offsetX-graphXTreshold+(int32(x)*element.dataMargin),
					offsetY+height-scaledValue,
					element.lineColor,
				)
			}
		}

		oldValue = scaledValue
	}

	isInRectangle := IsMouseInRectangle(
		offsetX,
		offsetY,
		width,
		height,
	)

	// shows specific approximation of a value on a graph
	if isInRectangle {
		mo := rl.GetMousePosition()
		m := [2]int32{
			int32(mo.X) / system.ScaleRatio,
			int32(mo.Y) / system.ScaleRatio,
		}

		// horizontal line
		rl.DrawLine(
			m[0],
			offsetY,
			m[0],
			offsetY+height,
			rl.Red,
		)

		var closestPointPastX int
		var skippedNodes int

		for x := range element.pointData {
			if (x*int(element.dataMargin) - int(graphXTreshold)) < int(m[0]-offsetX) {
				closestPointPastX = x
			} else {
				break
			}
		}

		adjustment := 1

		if len(element.pointData) == closestPointPastX+1 {
			adjustment = 0
		}

		y0 := float32(element.pointData[closestPointPastX])
		y1 := float32(element.pointData[closestPointPastX+adjustment])
		x0 := float32(closestPointPastX-skippedNodes) * float32(element.dataMargin)
		x1 := float32(closestPointPastX+adjustment-skippedNodes)*float32(element.dataMargin) + 1
		t := (float32(m[0]-offsetX) - x0) / (x1 - x0)

		finalY := float64(ScalarLerp(y0, y1, t))
		scaledFinalY := int32(float64(finalY-smallestValue) * float64(scaleY))

		// vertical line
		rl.DrawLine(
			offsetX,
			offsetY+height-scaledFinalY,
			offsetX+width,
			offsetY+height-scaledFinalY,
			rl.Red,
		)

		txt := fmt.Sprintf("%.02f %s", finalY, element.ValueSuffix)
		rl.DrawText(txt, m[0]+2, m[1]-10, 10, rl.RayWhite)
	}

	// Maxima
	rl.DrawText(
		fmt.Sprintf("%.02f %s", biggestValue, element.ValueSuffix),
		offsetX+width+5,
		offsetY-10,
		10,
		rl.RayWhite,
	)

	// Minima
	rl.DrawText(
		fmt.Sprintf("%.02f %s", smallestValue, element.ValueSuffix),
		offsetX+width+5,
		offsetY+height-10,
		10,
		rl.RayWhite,
	)

	// Average
	rl.DrawText(
		fmt.Sprintf("avg. %.02f %s", avgValue, element.ValueSuffix),
		offsetX+width+5,
		offsetY+(height/2)-10,
		10,
		rl.RayWhite,
	)

	scaledAvgValue := int32(float32(avgValue-smallestValue) * scaleY)
	scaledAvgY := scaledAvgValue

	rl.DrawLine(offsetX, offsetY+height-scaledAvgY, offsetX+width, offsetY+height-scaledAvgY, rl.RayWhite)

	return height + 5
}

func drawEditorElement(element *editorElement, offsetX, offsetY int32) int32 {
	color := rl.White
	var ext int32 = 10
	var textWidth = rl.MeasureText(element.text, 10)

	isInRectangle := IsMouseInRectangle(offsetX, offsetY, textWidth, 10)

	if element.isCollapsed != nil && isInRectangle {
		color = rl.Red

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			*element.isCollapsed = !*element.isCollapsed
		}
	} else if isInRectangle && element.callback != nil && element.isButton == false {
		color = rl.Yellow

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			element.callback()
		}
	} else if element.isButton {
		if isInRectangle {
			if rl.IsMouseButtonDown(rl.MouseLeftButton) {
				rl.DrawRectangle(
					offsetX-2,
					offsetY-2,
					textWidth+4,
					14,
					element.buttonPressColor,
				)
			} else {
				rl.DrawRectangle(
					offsetX-2,
					offsetY-2,
					textWidth+4,
					14,
					element.buttonHoverColor,
				)
			}

			if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				element.callback()
			}
		} else {
			rl.DrawRectangle(
				offsetX-2,
				offsetY-2,
				textWidth+4,
				14,
				element.buttonColor,
			)
		}

		ext += 2
	}

	rl.DrawText(element.text, offsetX+1, offsetY+1, 10, rl.Black)
	rl.DrawText(element.text, offsetX, offsetY, 10, color)

	if element.graphEnabled {
		ext += drawGraph(element, offsetX+5, offsetY+ext)
	}

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
