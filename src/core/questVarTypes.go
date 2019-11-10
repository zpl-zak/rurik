package core

import (
	"fmt"
	"math"

	rl "github.com/zaklaus/raylib-go/raylib"
)

type QuestVarNumber struct {
	Value float64
}

func (v *QuestVarNumber) Str() string {
	if math.Floor(v.Value) == v.Value {
		return fmt.Sprintf("%d", int64(v.Value))
	}

	return fmt.Sprintf("%f", v.Value)
}

type QuestVarVector struct {
	Value rl.Vector2
}

func (v *QuestVarVector) Str() string {
	return fmt.Sprintf("[%f, %f]", v.Value)
}
