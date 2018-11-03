package main

import (
	"fmt"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	// SkyColor is the tint color used for drawn sprites/tiles
	SkyColor rl.Color

	useTimeCycle    bool
	skyTime         float64
	skyTargetTime   float64
	skyStageIndex   int
	skyStages       []weatherStage
	skyLastColor    rl.Vector3
	skyCurrentColor rl.Vector3
	skyTargetColor  rl.Vector3
)

type weatherStage struct {
	color    rl.Vector3
	duration float64
}

// WeatherInit sets up the mood by initializing sky color tint and other properties
func WeatherInit() {
	var err error
	skyCurrentColor, err = getColorFromHex(tilemap.Properties.GetString("skyColor"))

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

	if len(skyStages) > 1 {
		skyLastColor = skyCurrentColor
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
}

// UpdateWeather updates the time cycle and weather effects
func UpdateWeather() {
	if useTimeCycle {
		if skyTime <= 0 {
			nextSkyStage()
		} else {
			skyTime -= float64(rl.GetFrameTime())
		}

		if skyTargetTime != 0 {
			skyCurrentColor = lerpColor(skyLastColor, skyTargetColor, 1-skyTime/skyTargetTime)
		} else {
			skyCurrentColor = skyTargetColor
		}

		SkyColor = vec3ToColor(skyCurrentColor)
	}
}

// DrawWeather draws weather effects
func DrawWeather() {
	if DebugMode {
		rl.DrawText(fmt.Sprintf("Sky stage: id %d val %v", skyStageIndex, skyStages[skyStageIndex]), 5, 20, 10, rl.White)
		rl.DrawText(fmt.Sprintf("Sky time: %f", skyTime), 5, 30, 10, rl.White)
	}
}

func nextSkyStage() {
	skyStageIndex++

	if skyStageIndex >= len(skyStages) {
		skyStageIndex = 0
	}

	stage := skyStages[skyStageIndex]
	skyTime = stage.duration
	skyTargetTime = skyTime
	skyTargetColor = stage.color
	skyLastColor = skyCurrentColor
}

func vec3ToColor(a rl.Vector3) rl.Color {
	return rl.NewColor(
		uint8(a.X*255),
		uint8(a.Y*255),
		uint8(a.Z*255),
		255,
	)
}

func lerpColor(a, b rl.Vector3, t float64) rl.Vector3 {
	return raymath.Vector3Lerp(a, b, float32(t))
}

func getColorFromHex(hex string) (rl.Vector3, error) {
	if hex == "" {
		return rl.Vector3{}, fmt.Errorf("hex not specified")
	}

	c, err := colorful.Hex("#" + hex[3:])

	if err != nil {
		return rl.Vector3{}, err
	}

	d := rl.NewVector3(
		float32(c.R),
		float32(c.G),
		float32(c.B),
	)

	return d, nil
}

func appendSkyStage(skyName, stageName string) {
	color, err := getColorFromHex(tilemap.Properties.GetString(skyName))

	if err == nil {
		useTimeCycle = true
	} else {
		return
	}

	duration, _ := strconv.ParseFloat(tilemap.Properties.GetString(stageName), 10)

	skyStages = append(skyStages, weatherStage{
		color:    color,
		duration: duration * 60,
	})
}
