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

		isInRectangle := IsMouseInRectangle(offsetX-graphXTreshold+(int32(x)*element.dataMargin)-2, offsetY+height-scaledValue-2, 4, 4)

		if isInRectangle {
			txt := fmt.Sprintf("%.02f", v)

			rl.DrawRectangle(offsetX-graphXTreshold+8+(int32(x)*element.dataMargin), offsetY+height-scaledValue-2, rl.MeasureText(txt, 10)+2, 12, rl.NewColor(30, 30, 30, 120))

			rl.DrawText(txt, offsetX-graphXTreshold+10+(int32(x)*element.dataMargin), offsetY+height-scaledValue, 10, rl.RayWhite)
		}

		oldValue = scaledValue
	}

	// Maxima
	rl.DrawText(
		fmt.Sprintf("%.02f", biggestValue),
		offsetX+width+5,
		offsetY-10,
		10,
		rl.RayWhite,
	)

	// Minima
	rl.DrawText(
		fmt.Sprintf("%.02f", smallestValue),
		offsetX+width+5,
		offsetY+height-10,
		10,
		rl.RayWhite,
	)

	// Average
	rl.DrawText(
		fmt.Sprintf("avg. %.02f", avgValue),
		offsetX+width+5,
		offsetY+(height/2)-10,
		10,
		rl.RayWhite,
	)

	scaledAvgValue := int32(float32(avgValue-smallestValue) * scaleY)
	scaledAvgY := offsetY + scaledAvgValue

	rl.DrawLine(offsetX, scaledAvgY, offsetX+width, scaledAvgY, rl.RayWhite)

	return height + 5
}

func drawEditorElement(element *editorElement, offsetX, offsetY int32) int32 {
	color := rl.White
	var ext int32 = 10

	isInRectangle := IsMouseInRectangle(offsetX, offsetY, rl.MeasureText(element.text, 10), 10)

	if element.isCollapsed != nil && isInRectangle {
		color = rl.Red

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			*element.isCollapsed = !*element.isCollapsed
		}
	} else if isInRectangle && element.callback != nil {
		color = rl.Yellow

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			element.callback()
		}
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
