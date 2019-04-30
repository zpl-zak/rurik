/*
   Copyright 2019 Dominik Madarász <zaklaus@madaraszd.net>

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
	"log"
	"math"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/raylib-go/raymath"
)

const (
	// CameraModeStatic stationary camera
	CameraModeStatic = 1

	// CameraModeFollow camera following an object
	CameraModeFollow = 2

	// CameraModeLerp camera transiting between two objects
	CameraModeLerp = 3
)

type camera struct {
	Follow     *Object
	Start      *Object
	End        *Object
	Speed      float32
	Progress   float32
	Zoom       float32
	TargetZoom float32
	ZoomSpeed  float32
	Smoothing  float32
	Mode       int
	First      bool
	FollowName string
}

type cameraData struct {
	Follow     string
	Start      string
	End        string
	Speed      float32
	Progress   float32
	TargetZoom float32
	ZoomSpeed  float32
	Smoothing  float32
	Mode       int
	First      bool
}

// NewCamera game camera
func (c *Object) NewCamera() {
	c.Update = updateCamera
	c.Finish = finishCamera
	c.Zoom = 1
	c.TargetZoom = c.Zoom
	c.ZoomSpeed = 0.8
	c.Smoothing = 0.8
	c.DebugVisible = false
	c.First = true
	followName := c.Meta.Properties.GetString("follow")
	if followName == "" {
		c.FollowName = "player"
	} else {
		c.FollowName = followName
	}
	c.Mode = CameraModeFollow
	strMode := c.Meta.Properties.GetString("mode")
	if strMode != "" {
		c.SetCameraMode(strMode)
	}
	spd, _ := strconv.ParseFloat(c.Meta.Properties.GetString("speed"), 32)

	if spd == 0 {
		spd = 1
	}

	c.Serialize = func(o *Object) string {
		val, _ := jsoniter.MarshalToString(&cameraData{
			Follow:     o.FollowName,
			Start:      "",
			End:        "",
			Speed:      o.Speed,
			Progress:   o.Progress,
			TargetZoom: o.TargetZoom,
			ZoomSpeed:  o.ZoomSpeed,
			Smoothing:  o.Smoothing,
			Mode:       o.Mode,
			First:      o.First,
		})

		return val
	}

	c.Deserialize = func(o *Object, v string) {
		var dat cameraData
		jsoniter.UnmarshalFromString(v, &dat)

		/* o.Follow = dat.Follow
		o.Start = dat.Start
		o.End = dat.End */
		o.Speed = dat.Speed
		o.Progress = dat.Progress
		o.TargetZoom = dat.TargetZoom
		o.ZoomSpeed = dat.ZoomSpeed
		o.Smoothing = dat.Smoothing
		o.Mode = dat.Mode
		o.First = dat.First
		o.FollowName = dat.Follow
	}

	c.Speed = float32(spd)

	MainCamera = c

	c.Draw = func(o *Object) {
		if !DebugMode || !o.DebugVisible {
			return
		}

		mode := "static"

		switch o.Mode {
		case 2:
			mode = "follow"
		case 3:
			mode = "lerp"
		}

		rl.DrawCircle(int32(o.Position.X), int32(o.Position.Y), 2, rl.White)
		DrawTextCentered(fmt.Sprintf("%s", o.Name), int32(o.Position.X), int32(o.Position.Y)+5, 10, rl.White)
		DrawTextCentered(fmt.Sprintf("Mode: %s", mode), int32(o.Position.X), int32(o.Position.Y)+15, 10, rl.White)
	}
}

func finishCamera(c *Object) {
	if c.Mode == CameraModeFollow {
		c.Follow, _ = c.world.FindObject(c.FollowName)

		if c.First && c.Follow != nil {
			c.Position = c.Follow.Position
		}
	} else if c.Mode == CameraModeLerp {
		c.Start, _ = c.world.FindObject(c.Meta.Properties.GetString("start"))
		c.End, _ = c.world.FindObject(c.Meta.Properties.GetString("end"))
	}
}

func updateCamera(c *Object, dt float32) {
	var dest rl.Vector2

	CanSave = BitsSet(CanSave, IsSequenceHappening)

	if c.Mode == CameraModeFollow {
		if c.Follow == nil {
			log.Println("Camera object follows nil reference.")
			return
		}

		if c.Follow == LocalPlayer {
			CanSave = BitsClear(CanSave, IsSequenceHappening)
		}

		dest = Vector2Lerp(c.Position, c.Follow.Position, c.Smoothing)
	} else if c.Mode == CameraModeLerp {
		if c.Start == nil || c.End == nil {
			log.Println("Camera object lerps between nil references.")
			return
		}

		if c.First {
			c.Position = c.Start.Position
		}

		dest = c.End.Position
	} else {
		dest = c.Position
	}

	dest.X = float32(math.Round(float64(dest.X)))
	dest.Y = float32(math.Round(float64(dest.Y)))

	if !c.First || c.Mode == CameraModeLerp {
		t := c.Speed

		if c.Mode == CameraModeLerp {
			vd := raymath.Vector2Distance(dest, c.Position)

			if c.Progress > 1 {
				c.Progress = 1

				c.Mode = CameraModeStatic

				if c.End.EventName != "" {
					FireEvent(c.End.EventName, c.End.EventArgs)
				}
			}

			if t > vd {
				t = vd
			}

			c.Position = Vector2Lerp(c.Start.Position, dest, c.Progress)
			c.Progress += t * dt
		} else {
			c.Position = Vector2Lerp(c.Position, dest, t)
		}
	} else {
		c.Position = dest
	}

	c.Zoom = ScalarLerp(c.Zoom, c.TargetZoom, c.ZoomSpeed*dt)

	c.First = false
}

// SetCameraZoom overrides camera zoom
func (c *Object) SetCameraZoom(t float32) {
	c.Zoom = t
	c.TargetZoom = t
}

// SetCameraMode sets the camera behavior mode
func (c *Object) SetCameraMode(strMode string) {
	switch strMode {
	default:
		fallthrough
	case "static":
		c.Mode = 1
	case "follow":
		c.Mode = 2
	case "lerp":
		c.Mode = 3
	}
}
