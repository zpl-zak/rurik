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

package system

import (
	"fmt"

	rl "github.com/zaklaus/raylib-go/raylib"
)

var (
	// Profilers consists of list of all profilers
	Profilers []*Profiler
)

// Profiler tracks the timing of particular operations.
type Profiler struct {
	ProfilerName    string
	StartTime       float64
	PassedTime      float64
	InvocationCount int32
	DisplayString   string

	IsCollapsed bool
}

// NewProfiler returns a new profiler
func NewProfiler(name string) *Profiler {
	prof := &Profiler{
		ProfilerName:  name,
		IsCollapsed:   false,
		DisplayString: name + ": 0 ms",
	}

	Profilers = append(Profilers, prof)
	return prof
}

// StartInvocation start timing this block
func (p *Profiler) StartInvocation() {
	p.StartTime = float64(rl.GetTime())
}

// StopInvocation stops timing this block
func (p *Profiler) StopInvocation() {
	p.PassedTime += float64(rl.GetTime()) - p.StartTime
	p.StartTime = 0
	p.InvocationCount++
}

// GetTime returns the passed time and resets the profiler
func (p *Profiler) GetTime(divisor float64) (result float64) {
	if divisor == 0 {
		divisor = float64(p.InvocationCount)
	}

	if p.PassedTime == 0 && divisor == 0 {
		result = 0
	} else {
		result = p.PassedTime / float64(divisor) * 1000
	}

	p.DisplayString = fmt.Sprintf("%s: %.02f ms", p.ProfilerName, result)
	p.PassedTime = 0
	p.InvocationCount = 0
	return
}
