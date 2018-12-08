/*
 * @Author: V4 Games
 * @Date: 2018-11-09 02:14:28
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-08 21:30:08
 */

package core

import (
	"fmt"

	"madaraszd.net/zaklaus/rurik/src/system"
)

var (
	updateProfiler     *system.Profiler
	collisionProfiler  *system.Profiler
	musicProfiler      *system.Profiler
	weatherProfiler    *system.Profiler
	gameModeProfiler   *system.Profiler
	drawProfiler       *system.Profiler
	sortRenderProfiler *system.Profiler

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
	gameModeProfiler = system.NewProfiler("custom")
	drawProfiler = system.NewProfiler("draw")
	sortRenderProfiler = system.NewProfiler("sortRender")

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
		}
		pushEditorElement(profilerNode, musicProfiler.DisplayString, nil)
		pushEditorElement(profilerNode, weatherProfiler.DisplayString, nil)
		pushEditorElement(profilerNode, gameModeProfiler.DisplayString, nil)

		renderNode := pushEditorElement(profilerNode, drawProfiler.DisplayString, &drawProfiler.IsCollapsed)

		if !drawProfiler.IsCollapsed {
			pushEditorElement(renderNode, sortRenderProfiler.DisplayString, nil)
		}
	}
}
