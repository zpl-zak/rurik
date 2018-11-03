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

	useTimeCycle    bool
	skyStageName    string
	skyTime         float64
	skyTargetTime   float64
	skyStageIndex   int
	skyStages       []weatherStage
	skyLastColor    rl.Vector3
	skyCurrentColor rl.Vector3
	skyTargetColor  rl.Vector3

	weatherIsCollapsed = true
)

type weatherStage struct {
	name     string
	color    rl.Vector3
	duration float64
}

// WeatherInit sets up the mood by initializing sky color tint and other properties
func WeatherInit() {
	var err error
	skyCurrentColor, err = getColorFromHex(CurrentMap.tilemap.Properties.GetString("skyColor"))

	if err != nil {
		SkyColor = rl.White
	} else {
		SkyColor = vec3ToColor(skyCurrentColor)
	}

	skyStages = []weatherStage{}

	appendSkyStage("skyRiseColor", "riseDuration")
	appendSkyStage("skyDayColor", "dayDuration")
	appendSkyStage("skyDawnColor", "dawnDuration")
	appendSkyStage("skyNightColor", "nightDuration")

	if len(skyStages) > 0 {
		skyLastColor = skyCurrentColor
		skyStageName = skyStages[0].name
		skyTargetColor = skyStages[0].color
		skyTime = skyStages[0].duration
		skyTargetTime = skyTime
		SkyColor = vec3ToColor(skyCurrentColor)
		skyStageIndex = 0

		if err != nil {
			skyCurrentColor = skyTargetColor
			nextSkyStage()
		}
	}

	weatherIsCollapsed = true
}

// UpdateWeather updates the time cycle and weather effects
func UpdateWeather() {
	if useTimeCycle {
		if skyTime <= 0 {
			nextSkyStage()
		} else {
			skyTime -= float64(rl.GetFrameTime()) * WeatherTimeScale
		}

		if skyTargetTime != 0 {
			skyCurrentColor = lerpColor(skyLastColor, skyTargetColor, 1-skyTime/skyTargetTime)
		} else {
			skyCurrentColor = skyTargetColor
		}

		SkyColor = vec3ToColor(skyCurrentColor)
	}

	if DebugMode {
		weatherElement := pushEditorElement(rootElement, "weather", &weatherIsCollapsed)
		pushEditorElement(weatherElement, fmt.Sprintf("Sky: %s (%d)", skyStageName, skyStageIndex), nil)
		pushEditorElement(weatherElement, fmt.Sprintf("Sky time: %d/%d", int(skyTargetTime-skyTime), int(skyTargetTime)), nil)
	}
}

// DrawWeather draws weather effects
func DrawWeather() {

}

func nextSkyStage() {
	skyStageIndex++

	if skyStageIndex >= len(skyStages) {
		skyStageIndex = 0
	}

	stage := skyStages[skyStageIndex]
	skyStageName = stage.name
	skyTime = stage.duration
	skyTargetTime = skyTime
	skyTargetColor = stage.color
	skyLastColor = skyCurrentColor
}

func appendSkyStage(skyName, stageName string) {
	color, err := getColorFromHex(CurrentMap.tilemap.Properties.GetString(skyName))

	if err == nil {
		useTimeCycle = true
	} else {
		return
	}

	duration, _ := strconv.ParseFloat(CurrentMap.tilemap.Properties.GetString(stageName), 10)

	skyStages = append(skyStages, weatherStage{
		name:     skyName,
		color:    color,
		duration: duration * 60,
	})
}
