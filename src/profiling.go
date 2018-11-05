package main

import "fmt"

var (
	updateProfiler    Profiler
	collisionProfiler Profiler
	animsProfiler     Profiler
	musicProfiler     Profiler
	weatherProfiler   Profiler
	customProfiler    Profiler
	drawProfiler      Profiler

	isProfilerCollapsed bool

	frameRateString = ""
	otherTimeString = ""
)

func initGameProfilers() {
	updateProfiler = NewProfiler("update")
	animsProfiler = NewProfiler("anims")
	collisionProfiler = NewProfiler("collision")
	musicProfiler = NewProfiler("music")
	weatherProfiler = NewProfiler("weather")
	customProfiler = NewProfiler("custom")
	drawProfiler = NewProfiler("draw")
}

func updateProfiling(frameCounter, frames float64) {
	totalTime := ((1000 * frameCounter) / (float64(frames)))
	var totalMeasuredTime float64

	totalMeasuredTime += updateProfiler.GetTime(frames)
	totalMeasuredTime += collisionProfiler.GetTime(frames)
	totalMeasuredTime += animsProfiler.GetTime(frames)
	totalMeasuredTime += musicProfiler.GetTime(frames)
	totalMeasuredTime += weatherProfiler.GetTime(frames)
	totalMeasuredTime += customProfiler.GetTime(frames)
	totalMeasuredTime += drawProfiler.GetTime(frames)

	frameRateString = fmt.Sprintf("total time: %.02f ms (%.02f FPS)", totalTime, 1000/totalTime)
	otherTimeString = fmt.Sprintf("measured time: %.02f ms", totalMeasuredTime)
}

func drawProfiling() {
	profilerNode := pushEditorElement(rootElement, "profiler", &isProfilerCollapsed)
	{
		pushEditorElement(profilerNode, frameRateString, nil)
		pushEditorElement(profilerNode, otherTimeString, nil)
		updateNode := pushEditorElement(profilerNode, updateProfiler.displayString, &updateProfiler.isCollapsed)
		{
			pushEditorElement(updateNode, collisionProfiler.displayString, nil)
			pushEditorElement(updateNode, animsProfiler.displayString, nil)
		}
		pushEditorElement(profilerNode, musicProfiler.displayString, nil)
		pushEditorElement(profilerNode, weatherProfiler.displayString, nil)
		pushEditorElement(profilerNode, customProfiler.displayString, nil)
		pushEditorElement(profilerNode, drawProfiler.displayString, nil)
	}
}
