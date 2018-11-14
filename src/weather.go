package main

import (
	"fmt"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// SkyColor is the tint color used for drawn sprites/tiles
	SkyColor rl.Color

	// WeatherTimeScale specifies time cycle scale
	WeatherTimeScale float64

	weatherIsCollapsed = true
)

type weatherStage struct {
	Name     string
	Color    rl.Vector3
	Duration float64
}

// Weather represents the map time and weather
type Weather struct {
	UseTimeCycle    bool
	SkyStageName    string
	SkyTime         float64
	SkyTargetTime   float64
	SkyStageIndex   int
	SkyStages       []weatherStage
	SkyLastColor    rl.Vector3
	SkyCurrentColor rl.Vector3
	SkyTargetColor  rl.Vector3
}

// WeatherInit sets up the mood by initializing Sky color tint and other properties
func (w *Weather) WeatherInit(cmap *Map) {
	var err error
	w.SkyCurrentColor, err = getColorFromHex(cmap.tilemap.Properties.GetString("skyColor"))

	if err != nil {
		SkyColor = rl.White
	} else {
		SkyColor = vec3ToColor(w.SkyCurrentColor)
	}

	w.SkyStages = []weatherStage{}

	w.appendSkyStage(cmap, "skyRiseColor", "riseDuration")
	w.appendSkyStage(cmap, "skyDayColor", "dayDuration")
	w.appendSkyStage(cmap, "skyDawnColor", "dawnDuration")
	w.appendSkyStage(cmap, "skyNightColor", "nightDuration")

	if len(w.SkyStages) > 0 {
		w.SkyLastColor = w.SkyCurrentColor
		w.SkyStageName = w.SkyStages[0].Name
		w.SkyTargetColor = w.SkyStages[0].Color
		w.SkyTime = w.SkyStages[0].Duration
		w.SkyTargetTime = w.SkyTime
		SkyColor = vec3ToColor(w.SkyCurrentColor)
		w.SkyStageIndex = 0

		if err != nil {
			w.SkyCurrentColor = w.SkyTargetColor
			w.nextSkyStage()
		}
	}

	weatherIsCollapsed = true
}

// UpdateWeather updates the time cycle and weather effects
func (w *Weather) UpdateWeather() {
	if w.UseTimeCycle {
		if w.SkyTime <= 0 {
			w.nextSkyStage()
		} else {
			w.SkyTime -= float64(FrameTime) * WeatherTimeScale
		}

		if w.SkyTargetTime != 0 {
			w.SkyCurrentColor = lerpColor(w.SkyLastColor, w.SkyTargetColor, 1-w.SkyTime/w.SkyTargetTime)
		} else {
			w.SkyCurrentColor = w.SkyTargetColor
		}

		SkyColor = vec3ToColor(w.SkyCurrentColor)
	}

	if DebugMode {
		weatherElement := pushEditorElement(rootElement, "weather", &weatherIsCollapsed)
		pushEditorElement(weatherElement, fmt.Sprintf("sky: %s (%d)", w.SkyStageName, w.SkyStageIndex), nil)
		pushEditorElement(weatherElement, fmt.Sprintf("sky time: %d/%d", int(w.SkyTargetTime-w.SkyTime), int(w.SkyTargetTime)), nil)
	}
}

// DrawWeather draws weather effects
func (w *Weather) DrawWeather() {

}

func (w *Weather) nextSkyStage() {
	w.SkyStageIndex++

	if w.SkyStageIndex >= len(w.SkyStages) {
		w.SkyStageIndex = 0
	}

	stage := w.SkyStages[w.SkyStageIndex]
	w.SkyStageName = stage.Name
	w.SkyTime = stage.Duration
	w.SkyTargetTime = w.SkyTime
	w.SkyTargetColor = stage.Color
	w.SkyLastColor = w.SkyCurrentColor
}

func (w *Weather) appendSkyStage(cmap *Map, SkyName, stageName string) {
	color, err := getColorFromHex(cmap.tilemap.Properties.GetString(SkyName))

	if err == nil {
		w.UseTimeCycle = true
	} else {
		return
	}

	duration, _ := strconv.ParseFloat(cmap.tilemap.Properties.GetString(stageName), 10)

	w.SkyStages = append(w.SkyStages, weatherStage{
		Name:     SkyName,
		Color:    color,
		Duration: duration * 60,
	})
}
