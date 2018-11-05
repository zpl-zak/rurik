package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
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
func NewProfiler(name string) Profiler {
	return Profiler{
		profilerName:  name,
		isCollapsed:   false,
		displayString: name + ": 0 ms",
	}
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
