package main

import (
	"github.com/lucasb-eyer/go-colorful"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// SkyColor is the tint color used for drawn sprites/tiles
	SkyColor rl.Color
)

// WeatherInit sets up the mood by initializing sky color tint and other properties
func WeatherInit() {
	skyColorText := tilemap.Properties.GetString("skyColor")

	if skyColorText != "" {
		c, _ := colorful.Hex(skyColorText)
		d := rl.NewColor(
			uint8(c.B*255),
			uint8(c.G*255),
			uint8(c.R*255),
			255,
		)
		SkyColor = d
	} else {
		SkyColor = rl.White
	}
}
