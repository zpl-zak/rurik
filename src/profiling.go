/*
 * @Author: V4 Games
 * @Date: 2018-11-09 02:14:28
 * @Last Modified by:   Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-11-09 02:14:28
 */

package main

import "fmt"

var (
	updateProfiler     *Profiler
	collisionProfiler  *Profiler
	musicProfiler      *Profiler
	weatherProfiler    *Profiler
	customProfiler     *Profiler
	drawProfiler       *Profiler
	sortRenderProfiler *Profiler

	isProfilerCollapsed bool

	frameRateString = ""
	otherTimeString = ""
)

func initGameProfilers() {
	updateProfiler = NewProfiler("update")
	collisionProfiler = NewProfiler("collision")
	musicProfiler = NewProfiler("music")
	weatherProfiler = NewProfiler("weather")
	customProfiler = NewProfiler("custom")
	drawProfiler = NewProfiler("draw")
	sortRenderProfiler = NewProfiler("sortRender")

	frameRateString = "total time: 0 ms (0 FPS)"
	otherTimeString = "measured time: 0 ms"
}

func updateProfiling(frameCounter, frames float64) {
	totalTime := ((1000 * frameCounter) / (float64(frames)))
	var totalMeasuredTime float64

	for _, x := range profilers {
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
		updateNode := pushEditorElement(profilerNode, updateProfiler.displayString, &updateProfiler.isCollapsed)

		if !updateProfiler.isCollapsed {
			pushEditorElement(updateNode, collisionProfiler.displayString, nil)
		}
		pushEditorElement(profilerNode, musicProfiler.displayString, nil)
		pushEditorElement(profilerNode, weatherProfiler.displayString, nil)
		pushEditorElement(profilerNode, customProfiler.displayString, nil)

		renderNode := pushEditorElement(profilerNode, drawProfiler.displayString, &drawProfiler.isCollapsed)

		if !drawProfiler.isCollapsed {
			pushEditorElement(renderNode, sortRenderProfiler.displayString, nil)
		}
	}
}
