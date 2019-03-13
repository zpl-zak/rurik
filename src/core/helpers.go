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
	"strconv"
	"strings"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/solarlune/resolv/resolv"
	tiled "github.com/zaklaus/go-tiled"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/raylib-go/raymath"
	"github.com/zaklaus/rurik/src/system"
)

const (
	// FrustumSafeMargin safe margin to be considered safe to render off-screen
	FrustumSafeMargin = 32.0
)

// Bits represent bitflags
type Bits uint64

// BitsSet sets a bit
func BitsSet(b, flag Bits) Bits { return b | flag }

// BitsClear clears a bit
func BitsClear(b, flag Bits) Bits { return b &^ flag }

// BitsToggle toggles a bit on/off
func BitsToggle(b, flag Bits) Bits { return b ^ flag }

// BitsHas checks if a bit is set on
func BitsHas(b, flag Bits) bool { return b&flag != 0 }

// CompileEventArgs returns cooked event arguments
func CompileEventArgs(args string) []string {
	evntArglist := args
	evntArgs := []string{evntArglist}

	if strings.Contains(evntArglist, ";") {
		evntArgs = strings.Split(evntArglist, ";")
	}

	return evntArgs
}

// RayRectangleInt32ToResolv conv
func rayRectangleInt32ToResolv(rec *resolv.Rectangle, i rl.RectangleInt32) {
	*rec = resolv.Rectangle{
		BasicShape: resolv.BasicShape{
			X:           i.X,
			Y:           i.Y,
			Collideable: true,
		},
		W: i.Width,
		H: i.Height,
	}
}

// DrawTextCentered draws a text that is centered
func DrawTextCentered(text string, posX, posY, fontSize int32, color rl.Color) {
	if fontSize < 10 {
		fontSize = 10
	}

	rl.DrawText(text, posX-rl.MeasureText(text, fontSize)/2, posY, fontSize, color)
}

// Vector2Lerp lerps vec2
func Vector2Lerp(v1, v2 rl.Vector2, amount float32) (result rl.Vector2) {
	result.X = v1.X + amount*(v2.X-v1.X)
	result.Y = v1.Y + amount*(v2.Y-v1.Y)

	return result
}

// ScalarLerp lerps a scalar value
func ScalarLerp(v1, v2 float32, amount float32) (result float32) {
	result = v1 + amount*(v2-v1)

	return result
}

// StringToVec2 conv
func StringToVec2(inp string) rl.Vector2 {
	comps := strings.Split(inp, " ")
	x, _ := strconv.ParseFloat(comps[0], 32)
	y, _ := strconv.ParseFloat(comps[1], 32)

	return rl.NewVector2(float32(x), float32(y))
}

// LerpColor lerps Color
func LerpColor(a, b rl.Vector3, t float64) rl.Vector3 {
	return raymath.Vector3Lerp(a, b, float32(t))
}

// GetColorFromHex conv
func GetColorFromHex(hex string) (rl.Vector3, error) {
	if hex == "" {
		return rl.Vector3{}, fmt.Errorf("hex not specified")
	}

	c, err := colorful.Hex("#" + hex[3:])

	if err != nil {
		return rl.Vector3{}, err
	}

	d := rl.NewVector3(
		float32(c.R),
		float32(c.G),
		float32(c.B),
	)

	return d, nil
}

// Vec3ToColor conv
func Vec3ToColor(a rl.Vector3) rl.Color {
	return rl.NewColor(
		uint8(a.X*255),
		uint8(a.Y*255),
		uint8(a.Z*255),
		255,
	)
}

// ColorToVec3 conv
func ColorToVec3(a rl.Color) rl.Vector3 {
	return rl.NewVector3(
		float32(a.R)/255.0,
		float32(a.G)/255.0,
		float32(a.B)/255.0,
	)
}

// MixColor mixes two colors together
func MixColor(a, b rl.Color) rl.Color {
	return Vec3ToColor(raymath.Vector3Lerp(
		ColorToVec3(a),
		ColorToVec3(b),
		0.5,
	))
}

// IsMouseInRectangle checks whether a mouse is inside of a rectangle
func IsMouseInRectangle(x, y, x2, y2 int32) bool {
	x2 = x + x2
	y2 = y + y2

	mo := rl.GetMousePosition()
	m := [2]int32{
		int32(mo.X) / system.ScaleRatio,
		int32(mo.Y) / system.ScaleRatio,
	}

	if m[0] > x && m[0] < x2 &&
		m[1] > y && m[1] < y2 {
		return true
	}

	return false
}

// GetSpriteAABB retrieves Aseprite boundaries
func GetSpriteAABB(o *Object) rl.RectangleInt32 {
	if o.Ase == nil {
		return rl.RectangleInt32{
			X:      int32(o.Position.X),
			Y:      int32(o.Position.Y - 32),
			Width:  32,
			Height: 32,
		}
	}

	return rl.RectangleInt32{
		X:      int32(o.Position.X) - int32(float32(o.Ase.FrameWidth/2)) + int32(float32(o.Ase.FrameWidth/4)),
		Y:      int32(o.Position.Y),
		Width:  o.Ase.FrameWidth / 2,
		Height: o.Ase.FrameHeight / 2,
	}
}

// PlayAnim plays an animation for a given object
func PlayAnim(p *Object, animName string) {
	if p.Ase.GetAnimation(animName) != nil {
		p.Ase.Play(animName)
	} else {
		//log.Println("Animation name:", animName, "not found!")
	}
}

// GetSpriteRectangle retrieves sprite's bounds
func GetSpriteRectangle(o *Object) rl.Rectangle {
	sourceX, sourceY := o.Ase.GetFrameXY()
	return rl.NewRectangle(float32(sourceX), float32(sourceY), float32(o.Ase.FrameWidth), float32(o.Ase.FrameHeight))
}

// GetSpriteOrigin retrieves sprite's origin
func GetSpriteOrigin(o *Object) rl.Rectangle {
	return rl.NewRectangle(o.Position.X-float32(o.Ase.FrameWidth/2), o.Position.Y-float32(o.Ase.FrameHeight/2), float32(o.Ase.FrameWidth), float32(o.Ase.FrameHeight))
}

// IsPointWithinRectangle checks whether a point is within a rectangle
func IsPointWithinRectangle(p rl.Vector2, r rl.Rectangle) bool {
	if p.X > r.X && p.X < (r.X+r.Width) &&
		p.Y > r.Y && p.Y < (r.Y+r.Height) {
		return true
	}

	return false
}

// GetColorFromProperty retrieves a color from property
func GetColorFromProperty(o *tiled.Object, name string) rl.Color {
	colorHex := o.Properties.GetString(name)
	var color rl.Color

	if colorHex != "" {
		colorVec, _ := GetColorFromHex(colorHex)
		color = Vec3ToColor(colorVec)
	} else {
		color = rl.Blank
	}

	return color
}

// GetVector2FromProperty retrieves a Vector2 from property
func GetVector2FromProperty(o *tiled.Object, name string) rl.Vector2 {
	txtVec := o.Properties.GetString(name)
	var vec rl.Vector2

	if txtVec != "" {
		vec = StringToVec2(txtVec)
	}

	return vec
}

// GetFloatFromProperty retrieves a float from property
func GetFloatFromProperty(o *tiled.Object, name string) float32 {
	fltString := o.Properties.GetString(name)
	var flt float32

	if fltString != "" {
		fltRaw, _ := strconv.ParseFloat(fltString, 32)
		flt = float32(fltRaw)
	} else {
		flt = 0
	}

	return flt
}

// IsPointWithinFrustum checks whether a point is within camera's frustum
func IsPointWithinFrustum(p rl.Vector2) bool {
	if MainCamera == nil {
		return false
	}

	camOffset := rl.Vector2{
		X: float32(int(MainCamera.Position.X - float32(system.ScreenWidth)/2/MainCamera.Zoom)),
		Y: float32(int(MainCamera.Position.Y - float32(system.ScreenHeight)/2/MainCamera.Zoom)),
	}

	cam := rl.Rectangle{
		X:      camOffset.X - FrustumSafeMargin,
		Y:      camOffset.Y - FrustumSafeMargin,
		Width:  float32(system.ScreenWidth)/MainCamera.Zoom + FrustumSafeMargin*2,
		Height: float32(system.ScreenHeight)/MainCamera.Zoom + FrustumSafeMargin*2,
	}

	return IsPointWithinRectangle(p, cam)
}
