package core

import "math"

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

// SignFloat returns a value's sign
func SignFloat(x float32) float32 {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	} else {
		return 0
	}
}

// RoundFloat rounds a value
func RoundFloat(x float32) float32 {
	return float32(math.Round(float64(x)))
}

// RoundInt32ToFloat rounds a value
func RoundInt32ToFloat(x int32) float32 {
	return float32(math.Round(float64(x)))
}

// RoundFloatToInt32 rounds a value and casts it into int32
func RoundFloatToInt32(x float32) int32 {
	return int32(math.Round(float64(x)))
}
