package core

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/raylib-go/raymath"
)

// Vec2 2-component vector
type Vec2 struct {
	X int32
	Y int32
}

// NewVec2 creates a new Vec2
func NewVec2(x, y int32) Vec2 {
	return Vec2{X: x, Y: y}
}

// Vec2FromRayVector2 rounds raylib Vector2 and casts it to Vec2
func Vec2FromRayVector2(v2 rl.Vector2) Vec2 {
	return Vec2{
		X: RoundFloatToInt32(v2.X),
		Y: RoundFloatToInt32(v2.Y),
	}
}

// RayVector2FromVec2 casts Vec2 to rl.Vector2
func RayVector2FromVec2(v2 Vec2) rl.Vector2 {
	return rl.NewVector2(
		float32(v2.X),
		float32(v2.Y),
	)
}

// Vec2Normalize normalizes a Vec2
func Vec2Normalize(v2 Vec2) Vec2 {
	oldVec := RayVector2FromVec2(v2)
	raymath.Vector2Normalize(&oldVec)
	return Vec2FromRayVector2(oldVec)
}

// SignInt32 returns a value's sign
func SignInt32(x int32) int32 {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	} else {
		return 0
	}
}
