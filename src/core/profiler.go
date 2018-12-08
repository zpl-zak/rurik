/*
 * @Author: V4 Games
 * @Date: 2018-11-09 02:14:45
 * @Last Modified by:   Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-11-09 02:14:45
 */

package core

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	profilers []*Profiler
)

// Profiler tracks the timing of particular operations.
type Profiler struct {
	profilerName    string
	startTime       float64
	passedTime      float64
	invocationCount int32
	displayString   string

	isCollapsed bool
}

// NewProfiler returns a new profiler
func NewProfiler(name string) *Profiler {
	prof := &Profiler{
		profilerName:  name,
		isCollapsed:   false,
		displayString: name + ": 0 ms",
	}

	profilers = append(profilers, prof)
	return prof
}

// StartInvocation start timing this block
func (p *Profiler) StartInvocation() {
	p.startTime = float64(rl.GetTime())
}

// StopInvocation stops timing this block
func (p *Profiler) StopInvocation() {
	p.passedTime += float64(rl.GetTime()) - p.startTime
	p.startTime = 0
	p.invocationCount++
}

// GetTime returns the passed time and resets the profiler
func (p *Profiler) GetTime(divisor float64) (result float64) {
	if divisor == 0 {
		divisor = float64(p.invocationCount)
	}

	if p.passedTime == 0 && divisor == 0 {
		result = 0
	} else {
		result = p.passedTime / float64(divisor) * 1000
	}

	p.displayString = fmt.Sprintf("%s: %.02f ms", p.profilerName, result)
	p.passedTime = 0
	p.invocationCount = 0
	return
}
