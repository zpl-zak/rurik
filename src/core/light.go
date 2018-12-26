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
	"strconv"

	rl "github.com/zaklaus/raylib-go/raylib"
)

type light struct {
	color  rl.Color
	radius float32
}

// NewLight light instance
func (o *Object) NewLight() {
	o.Init = func(o *Object) {
		hexColor := o.Meta.Properties.GetString("color")

		if hexColor != "" {
			col, _ := getColorFromHex(hexColor)
			o.color = vec3ToColor(col)
		}

		radius := o.Meta.Properties.GetString("radius")

		if radius != "" {
			rad, _ := strconv.ParseFloat(radius, 32)
			o.radius = float32(rad)
		}
	}
}
