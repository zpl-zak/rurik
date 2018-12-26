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

	jsoniter "github.com/json-iterator/go"
	rl "github.com/zaklaus/raylib-go/raylib"
)

type light struct {
	Color       rl.Color
	Attenuation float32
}

type lightData struct {
	Color       rl.Color `json:"col"`
	Attenuation float32  `json:"atten"`
	Radius      float32  `json:"rad"`
}

// NewLight light instance
func (o *Object) NewLight() {
	o.Serialize = func(o *Object) string {
		data := lightData{
			Color:       o.Color,
			Attenuation: o.Attenuation,
			Radius:      o.Radius,
		}

		ret, _ := jsoniter.MarshalToString(&data)
		return ret
	}

	o.Deserialize = func(o *Object, in string) {
		var data lightData
		jsoniter.UnmarshalFromString(in, &data)

		o.Color = data.Color
		o.Attenuation = data.Attenuation
		o.Radius = data.Radius
	}

	o.Init = func(o *Object) {
		hexColor := o.Meta.Properties.GetString("color")

		if hexColor != "" {
			col, _ := getColorFromHex(hexColor)
			o.Color = vec3ToColor(col)
		}

		radius := o.Meta.Properties.GetString("radius")

		if radius != "" {
			rad, _ := strconv.ParseFloat(radius, 32)
			o.Radius = float32(rad)
		}

		attenuation := o.Meta.Properties.GetString("atten")

		if attenuation != "" {
			rad, _ := strconv.ParseFloat(attenuation, 32)
			o.Attenuation = float32(rad)
		}
	}
}
