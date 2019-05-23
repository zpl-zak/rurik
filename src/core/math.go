package core

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
