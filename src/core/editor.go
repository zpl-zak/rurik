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
	"github.com/zaklaus/raylib-go/raymath"
	"github.com/zaklaus/rurik/src/system"
)

var (
	rootIsCollapsed = false
	rootElement     = &EditorElement{
		Text:        "editor",
		IsCollapsed: &rootIsCollapsed,
		Children:    []*EditorElement{},
	}

	editorElementCounter uint64

	sliderHandleID uint64
)

const (
	defaultGraphHeight     = 50
	defaultGraphWidth      = 200
	defaultGraphDataMargin = 5

	defaultSliderWidth             = 100
	defaultSliderHeight            = 10
	defaultSliderHandleWidth       = 10
	defaultSliderHandleVisualWidth = 10
	defaultSliderHandleHeight      = 10

	defaultSliderValueMin  = 0
	defaultSliderValueMax  = 100
	defaultSliderValueStep = 1
)

const (
	elementTypeStandard = iota
	elementTypeButton
	elementTypeSlider
)

// EditorElement represents an editor UI element
type EditorElement struct {
	ID           uint64
	Class        uint8
	Text         string
	IsCollapsed  *bool
	IsHorizontal bool
	Padding      rl.RectangleInt32
	Children     []*EditorElement
	Callback     func()

	// Graphs
	GraphEnabled bool
	LineColor    rl.Color
	DataMargin   int32
	GraphHeight  int32
	GraphWidth   int32
	PointData    []float64
	UseCurves    bool
	ValueSuffix  string

	// Buttons
	ButtonColor      rl.Color
	ButtonHoverColor rl.Color
	ButtonPressColor rl.Color

	// Sliders
	SliderValue          *float64
	SliderValueMin       float64
	SliderValueMax       float64
	SliderValueRounding  int64
	SliderValueLimitless bool
}

// GetRootElement retrieves root UI element
func GetRootElement() *EditorElement {
	return rootElement
}

// PushEditorElement adds new UI element to the editor
func PushEditorElement(element *EditorElement, Text string, IsCollapsed *bool) *EditorElement {
	return PushEditorElementEx(element, Text, IsCollapsed, func() {})
}

// PushEditorElementEx adds new UI element to the editor using extended features
func PushEditorElementEx(element *EditorElement, Text string, IsCollapsed *bool, Callback func()) *EditorElement {

	editorElementCounter++
	child := &EditorElement{
		ID:          editorElementCounter,
		Text:        Text,
		IsCollapsed: IsCollapsed,
		Children:    []*EditorElement{},
		Callback:    Callback,
	}
	element.Children = append(element.Children, child)

	return child
}

// SetUpButton initializes the UI element as a button
func SetUpButton(element *EditorElement, Callback func(), IsHorizontal bool) {
	element.Class = elementTypeButton
	element.ButtonColor = rl.DarkPurple
	element.ButtonHoverColor = rl.Purple
	element.ButtonPressColor = rl.Pink
	element.IsHorizontal = IsHorizontal
	element.Callback = Callback
}

// SetUpSlider initializes the UI element as a slider
func SetUpSlider(element *EditorElement, value *float64, min, max float64) {
	element.Class = elementTypeSlider
	element.SliderValue = value
	if min != max {
		element.SliderValueMax = max
		element.SliderValueMin = min
	} else {
		element.SliderValueLimitless = true
	}
	element.SliderValueRounding = 2
}

var (
	editorPickMode     = false
	editorPickObject   *Object
	editorPickMousePos [2]int32
	editorPickOffset   rl.Vector2
	editorPickVector   rl.Vector2
	editorPickCalm     = 0
)

func editorHandleObjectTransform(o *Object) {
	if editorPickObject == nil || !rl.IsKeyDown(rl.KeyLeftShift) {
		return
	}

	mo := GetMousePosition2D()

	if editorPickMousePos != mo {
		editorPickVector = raymath.Vector2Subtract(
			IntArrayToVector2(mo),
			IntArrayToVector2(editorPickMousePos),
		)

		editorPickCalm = 0
	} else {
		editorPickCalm++
	}

	editorPickObject.Position = rl.NewVector2(
		float32(mo[0])+editorPickOffset.X,
		float32(mo[1])+editorPickOffset.Y,
	)

	editorPickMousePos = mo

	editorPickObject.Movement = rl.Vector2{}
}

func drawObjectDebug2D(o *Object) {
	rect := o.GetAABB(o)
	col := rl.RayWhite

	mouseOver := IsMouseInRectangle2D(rect)

	if mouseOver {
		col = rl.Yellow

		mo := IntArrayToVector2(GetMousePosition2D())
		rl.DrawLineV(o.Position, mo, rl.Yellow)
	}

	rl.DrawRectangleLines(
		rect.X,
		rect.Y,
		rect.Width,
		rect.Height,
		col,
	)

	DrawTextCentered(o.Name, rect.X+rect.Width/2, rect.Y+rect.Height+2, 1, rl.White)
}

// DrawEditor draws debug UI
func DrawEditor() {
	if DebugMode {
		handleEditorElement(rootElement, 5, 5)
	}

	if CurrentMap != nil {
		for _, o := range CurrentMap.World.Objects {
			if !DebugMode || !o.DebugVisible {
				continue
			}

			rect := o.GetAABB(o)

			if IsMouseInRectangle2D(rect) {
				pos := WorldToScreenPosRec(rect)
				mo := GetMousePosition2D()
				dir := raymath.Vector2Subtract(
					IntArrayToVector2(mo),
					o.Position,
				)
				offset := Vector2ToIntArray(dir)

				rl.DrawText(
					fmt.Sprintf(
						"X: %.02f\nY: %.02f\nW: %d\nH: %d\nOX: %d\n OY: %d\n",
						o.Position.X,
						o.Position.Y,
						rect.Width,
						rect.Height,
						offset[0],
						offset[1],
					),
					pos[0]+int32(float32(rect.Width+2)*MainCamera.Zoom),
					pos[1],
					10,
					rl.Yellow,
				)

				if rl.IsKeyDown(rl.KeyLeftShift) && editorPickObject == nil {
					editorPickObject = o

					editorPickOffset = raymath.Vector2Subtract(
						o.Position,
						IntArrayToVector2(mo),
					)
				}
			}

			editorHandleObjectTransform(o)
		}

		if !rl.IsKeyDown(rl.KeyLeftShift) && editorPickObject != nil {
			if editorPickCalm < 2 {
				editorPickObject.Movement = editorPickVector
			}
			editorPickObject = nil
		}
	}
}

// UpdateEditor updates editor debug UI
func UpdateEditor() {
	flushEditorElement()

	DebugShowAll = !rootIsCollapsed
}

func drawGraph(element *EditorElement, offsetX, offsetY int32) int32 {
	var height, width int32

	height = element.GraphHeight
	width = element.GraphWidth

	rl.DrawRectangle(offsetX, offsetY, width, height, rl.NewColor(40, 40, 40, 140))

	// draw grid
	gridColumns := int(width / element.DataMargin)
	for x := 0; x < gridColumns; x++ {
		posX := int32(x * int(element.DataMargin))
		rl.DrawLine(
			offsetX+posX,
			offsetY,
			offsetX+posX,
			offsetY+height,
			rl.NewColor(255, 255, 255, 40),
		)
	}
	gridRows := int(height / element.DataMargin)
	for x := 0; x < gridRows; x++ {
		posY := int32(x * int(element.DataMargin))
		rl.DrawLine(
			offsetX,
			offsetY+posY,
			offsetX+width,
			offsetY+posY,
			rl.NewColor(255, 255, 255, 40),
		)
	}

	if element.PointData == nil || (element.PointData != nil && len(element.PointData) < 1) {
		return height + 5
	}

	// Value mapping
	var biggestValue float64
	var smallestValue = math.MaxFloat64
	var sum float64
	var avgValue float64
	var nodeCount int32

	var graphXTreshold int32
	actualGraphWidth := int32(len(element.PointData)) * element.DataMargin

	if actualGraphWidth > width {
		graphXTreshold = actualGraphWidth - width - element.DataMargin
	}

	for x, v := range element.PointData {
		if (int32(x) * element.DataMargin) < graphXTreshold {
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
	for x, v := range element.PointData {
		if (int32(x) * element.DataMargin) < graphXTreshold {
			continue
		}

		scaledValue := int32(float32(v-smallestValue) * float32(scaleY))
		rl.DrawCircle(offsetX-graphXTreshold+(int32(x)*element.DataMargin), offsetY+height-scaledValue, 1, element.LineColor)

		if oldValue != -1 {
			if element.UseCurves {
				rl.DrawLineBezier(
					rl.NewVector2(
						float32(offsetX-graphXTreshold+(int32(x-1)*element.DataMargin)),
						float32(offsetY+height-oldValue),
					),
					rl.NewVector2(
						float32(offsetX-graphXTreshold+(int32(x)*element.DataMargin)),
						float32(offsetY+height-scaledValue),
					),
					1,
					element.LineColor,
				)
			} else {
				rl.DrawLine(
					offsetX-graphXTreshold+(int32(x-1)*element.DataMargin),
					offsetY+height-oldValue,
					offsetX-graphXTreshold+(int32(x)*element.DataMargin),
					offsetY+height-scaledValue,
					element.LineColor,
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
		m := system.GetMousePosition()

		// horizontal line
		rl.DrawLine(
			m[0],
			offsetY,
			m[0],
			offsetY+height,
			rl.Red,
		)

		var closestPointPastX int

		for x := range element.PointData {
			if (x*int(element.DataMargin) - int(graphXTreshold)) < int(m[0]-offsetX) {
				closestPointPastX = x
			} else {
				break
			}
		}

		adjustment := 1

		if len(element.PointData) == closestPointPastX+1 {
			adjustment = 0
		}

		y0 := float32(element.PointData[closestPointPastX])
		y1 := float32(element.PointData[closestPointPastX+adjustment])
		x0 := float32(closestPointPastX) * float32(element.DataMargin)
		x1 := float32(closestPointPastX+adjustment)*float32(element.DataMargin) + 1
		t := (float32(m[0]-offsetX) + float32(graphXTreshold) - x0) / (x1 - x0)
		if t > 1 {
			t = 1
		} else if t < 0 {
			t = 0
		}

		finalY := float64(ScalarLerp(y0, y1, t))
		scaledFinalY := int32(float64(finalY-smallestValue) * float64(scaleY))

		// vertical line (fixed)
		rl.DrawLine(
			offsetX,
			offsetY+height-scaledFinalY,
			offsetX+width,
			offsetY+height-scaledFinalY,
			rl.Red,
		)
		// vertical line (free)
		rl.DrawLine(
			offsetX,
			m[1],
			offsetX+width,
			m[1],
			rl.NewColor(255, 0, 0, 140),
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

func drawButton(element *EditorElement, offsetX, offsetY, TextWidth int32, isInRectangle bool) {
	if isInRectangle {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			rl.DrawRectangle(
				offsetX-2,
				offsetY-2,
				TextWidth+4,
				14,
				element.ButtonPressColor,
			)
		} else {
			rl.DrawRectangle(
				offsetX-2,
				offsetY-2,
				TextWidth+4,
				14,
				element.ButtonHoverColor,
			)
		}

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			element.Callback()
		}
	} else {
		rl.DrawRectangle(
			offsetX-2,
			offsetY-2,
			TextWidth+4,
			14,
			element.ButtonColor,
		)
	}
}

func drawSlider(element *EditorElement, offsetX, offsetY int32, TextWidth int32) {
	offsetX += TextWidth + 3

	rl.DrawRectangle(
		offsetX,
		offsetY,
		defaultSliderWidth,
		defaultSliderHeight,
		rl.NewColor(80, 40, 80, 200),
	)

	rl.DrawText(
		fmt.Sprintf("%.02f", *element.SliderValue),
		offsetX+defaultSliderWidth+3,
		offsetY,
		10,
		rl.RayWhite,
	)

	if !element.SliderValueLimitless {
		rl.DrawText(
			fmt.Sprintf("%.02f", element.SliderValueMin),
			offsetX,
			offsetY,
			10,
			rl.DarkPurple,
		)

		maxTxt := fmt.Sprintf("%.02f", element.SliderValueMax)
		maxTxtWidth := rl.MeasureText(maxTxt, 10)
		rl.DrawText(
			maxTxt,
			offsetX+defaultSliderWidth-maxTxtWidth,
			offsetY,
			10,
			rl.DarkPurple,
		)
	}

	scaledPositionX := float64(defaultSliderWidth / 2)
	scaleX := 1.0

	if !element.SliderValueLimitless {
		scaleX = float64(defaultSliderWidth) / float64(element.SliderValueMax-element.SliderValueMin)
		scaledPositionX = float64((*element.SliderValue - element.SliderValueMin) * scaleX)
	}

	isInRectangle := IsMouseInRectangle(
		offsetX+int32(scaledPositionX)-defaultSliderHandleWidth-defaultSliderHandleVisualWidth/8,
		offsetY,
		defaultSliderHandleWidth*2,
		defaultSliderHandleHeight,
	)

	if isInRectangle || sliderHandleID == element.ID {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			sliderHandleID = element.ID
		} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			sliderHandleID = 0
		}

		if element.ID == sliderHandleID {
			if element.SliderValueLimitless {
				m := system.MouseDelta
				scaledPositionX = *element.SliderValue + float64(m[0])
				*element.SliderValue = scaledPositionX
			} else {
				m := system.GetMousePosition()
				scaledPositionX = float64(m[0]-offsetX) + defaultSliderHandleVisualWidth/4
				*element.SliderValue = scaledPositionX - (defaultSliderWidth / 2)
			}

		}

		var extraSliderOffset int32

		if element.SliderValueLimitless && sliderHandleID != 0 {
			extraSliderOffset = +defaultSliderHandleVisualWidth/4 + defaultSliderWidth/2
		}

		rl.DrawRectangle(
			offsetX+int32(scaledPositionX)-defaultSliderHandleVisualWidth/4+extraSliderOffset,
			offsetY,
			defaultSliderHandleVisualWidth/4,
			defaultSliderHandleHeight,
			rl.Pink,
		)
	} else {
		rl.DrawRectangle(
			offsetX+int32(scaledPositionX)-defaultSliderHandleVisualWidth/4,
			offsetY,
			defaultSliderHandleVisualWidth/4,
			defaultSliderHandleHeight,
			rl.Purple,
		)
	}

	if !element.SliderValueLimitless {
		*element.SliderValue = scaledPositionX/scaleX + element.SliderValueMin
		if *element.SliderValue < element.SliderValueMin {
			*element.SliderValue = element.SliderValueMin
		} else if *element.SliderValue > element.SliderValueMax {
			*element.SliderValue = element.SliderValueMax
		}
	}

	minSteps := []float64{1.0, 0.1, 0.01, 0.001, 0.0001, 0.00001, 0.000001, 0.0000001, 0.00000001, 0.000000001}
	var decimalPrecision float64
	if element.SliderValueRounding >= 0 && element.SliderValueRounding < 10 {
		decimalPrecision = minSteps[element.SliderValueRounding]
	} else {
		decimalPrecision = math.Pow10(int(-element.SliderValueRounding))
	}

	*element.SliderValue = math.Round(*element.SliderValue/decimalPrecision) * decimalPrecision
}

func handleEditorElement(element *EditorElement, offsetX, offsetY int32) (int32, int32, int32) {
	color := rl.White
	var ext int32 = 10
	var TextWidth = rl.MeasureText(element.Text, 10)
	var ext2 = TextWidth
	var totale2 = ext2 + offsetX - 5

	offsetX += element.Padding.X
	offsetY += element.Padding.Y
	ext += element.Padding.Height
	ext2 += element.Padding.Width
	var buttonWidth int32
	var buttonHeight int32

	if element.Class == elementTypeButton {
		buttonWidth = 4
		buttonHeight = 8
	}

	isInRectangle := IsMouseInRectangle(offsetX-buttonWidth, offsetY, TextWidth+buttonWidth, 10+buttonHeight)

	if element.IsCollapsed != nil && isInRectangle {
		color = rl.Red

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			*element.IsCollapsed = !*element.IsCollapsed
		}
	} else if isInRectangle && element.Callback != nil && element.Class == elementTypeStandard {
		color = rl.Yellow

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			element.Callback()
		}
	} else if element.Class == elementTypeButton {
		offsetY += 5
		drawButton(element, offsetX, offsetY, TextWidth, isInRectangle)
		ext += 8
		ext2 += 5
		totale2 += 5
	} else if element.Class == elementTypeSlider {
		offsetY += 5
		drawSlider(element, offsetX, offsetY, TextWidth)
		ext += 8
		ext2 += 5
		totale2 += 5
	}

	rl.DrawText(element.Text, offsetX+1, offsetY+1, 10, rl.Black)
	rl.DrawText(element.Text, offsetX, offsetY, 10, color)

	if element.GraphEnabled && element.IsCollapsed != nil && *element.IsCollapsed == false {
		ext += drawGraph(element, offsetX+5, offsetY+ext)
	}

	if element.IsCollapsed != nil && *element.IsCollapsed {
		return ext, ext2, totale2
	}

	var lastChildWidth int32
	var lastChildHeight int32

	for x, v := range element.Children {
		if x == 0 && v.IsHorizontal && v.Class != elementTypeStandard {
			v.IsHorizontal = false
		}

		var extraOffsetX int32
		var extraOffsetY int32
		if v.IsHorizontal {
			extraOffsetX = lastChildWidth + 5

			if v.Class != elementTypeStandard {
				extraOffsetY = lastChildHeight
			} else {
				extraOffsetX = 0
				if x != 0 && element.Children[x-1].IsHorizontal {
					extraOffsetX = totale2
				}
			}
		} else {
			extraOffsetX = 0
			extraOffsetY = 0
		}
		rext, rext2, rtotal := handleEditorElement(v, offsetX+5+extraOffsetX, offsetY+ext-extraOffsetY)
		if !v.IsHorizontal {
			lastChildWidth = rext2
			ext += rext
		} else {
			lastChildWidth += rext2 + 5
		}
		if rtotal > totale2 {
			totale2 = rtotal
		}
		lastChildHeight = rext
	}

	return ext, ext2, totale2
}

func flushEditorElement() {
	rootElement = &EditorElement{
		Text:        "editor",
		IsCollapsed: &rootIsCollapsed,
		Children:    []*EditorElement{},
	}

	editorElementCounter = 0
}
