/*
   Copyright 2018 V4 Games

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

	"github.com/zaklaus/rurik/src/system"
)

var (
	updateProfiler     *system.Profiler
	collisionProfiler  *system.Profiler
	musicProfiler      *system.Profiler
	weatherProfiler    *system.Profiler
	gameModeProfiler   *system.Profiler
	drawProfiler       *system.Profiler
	sortRenderProfiler *system.Profiler
	lightingProfiler   *system.Profiler
	scriptingProfiler  *system.Profiler

	isProfilerCollapsed bool

	frameRateString = ""
	otherTimeString = ""
)

// InitGameProfilers initializes all game profilers used within the engine
func InitGameProfilers() {
	updateProfiler = system.NewProfiler("update")
	collisionProfiler = system.NewProfiler("collision")
	musicProfiler = system.NewProfiler("music")
	weatherProfiler = system.NewProfiler("weather")
	gameModeProfiler = system.NewProfiler("gameMode")
	drawProfiler = system.NewProfiler("draw")
	sortRenderProfiler = system.NewProfiler("sortRender")
	lightingProfiler = system.NewProfiler("lighting")
	scriptingProfiler = system.NewProfiler("scripting")

	frameRateString = "total time: 0 ms (0 FPS)"
	otherTimeString = "measured time: 0 ms"
}

func updateProfiling(frameCounter, frames float64) {
	totalTime := ((1000 * frameCounter) / (float64(frames)))
	var totalMeasuredTime float64

	for _, x := range system.Profilers {
		totalMeasuredTime += x.GetTime(frames)
	}

	frameRateString = fmt.Sprintf("total time: %.02f ms (%.02f FPS)", totalTime, 1000/totalTime)
	otherTimeString = fmt.Sprintf("measured time: %.02f ms", totalMeasuredTime)
}

func drawProfiling() {
	profilerNode := pushEditorElement(rootElement, "profiler", &isProfilerCollapsed)

	if !isProfilerCollapsed {
		pushEditorElement(profilerNode, frameRateString, nil)
		pushEditorElement(profilerNode, otherTimeString, nil)
		updateNode := pushEditorElement(profilerNode, updateProfiler.DisplayString, &updateProfiler.IsCollapsed)

		if !updateProfiler.IsCollapsed {
			pushEditorElement(updateNode, collisionProfiler.DisplayString, nil)
			pushEditorElement(updateNode, scriptingProfiler.DisplayString, nil)
		}
		pushEditorElement(profilerNode, musicProfiler.DisplayString, nil)
		pushEditorElement(profilerNode, weatherProfiler.DisplayString, nil)
		pushEditorElement(profilerNode, gameModeProfiler.DisplayString, nil)

		renderNode := pushEditorElement(profilerNode, drawProfiler.DisplayString, &drawProfiler.IsCollapsed)

		if !drawProfiler.IsCollapsed {
			pushEditorElement(renderNode, sortRenderProfiler.DisplayString, nil)
			pushEditorElement(renderNode, lightingProfiler.DisplayString, nil)
		}
	}
}
