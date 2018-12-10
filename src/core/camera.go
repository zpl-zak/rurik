/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:26:32
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 14:23:14
 */

package core

import (
	"fmt"
	"log"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"
	"github.com/json-iterator/go"
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
	Mode       int
	First      bool
}

type cameraData struct {
	Follow     string
	Start      string
	End        string
	Speed      float32
	Progress   float32
	TargetZoom float32
	ZoomSpeed  float32
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
	c.DebugVisible = false
	c.First = true
	strMode := c.Meta.Properties.GetString("mode")
	spd, _ := strconv.ParseFloat(c.Meta.Properties.GetString("speed"), 32)

	c.SetCameraMode(strMode)

	if spd == 0 {
		spd = 1
	}

	c.Serialize = func(o *Object) string {
		val, _ := jsoniter.MarshalToString(&cameraData{
			Follow:     "",
			Start:      "",
			End:        "",
			Speed:      o.Speed,
			Progress:   o.Progress,
			TargetZoom: o.TargetZoom,
			ZoomSpeed:  o.ZoomSpeed,
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
		o.Mode = dat.Mode
		o.First = dat.First
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
		drawTextCentered(fmt.Sprintf("%s", o.Name), int32(o.Position.X), int32(o.Position.Y)+5, 10, rl.White)
		drawTextCentered(fmt.Sprintf("Mode: %s", mode), int32(o.Position.X), int32(o.Position.Y)+15, 10, rl.White)
	}
}

func finishCamera(c *Object) {
	if c.Mode == CameraModeFollow {
		c.Follow, _ = c.world.FindObject(c.Meta.Properties.GetString("follow"))
	} else if c.Mode == CameraModeLerp {
		c.Start, _ = c.world.FindObject(c.Meta.Properties.GetString("start"))
		c.End, _ = c.world.FindObject(c.Meta.Properties.GetString("end"))
	}
}

func updateCamera(c *Object, dt float32) {
	var dest rl.Vector2

	CanSave = bitsSet(CanSave, isSequenceHappening)

	if c.Mode == CameraModeFollow {
		if c.Follow == nil {
			log.Println("Camera object follows nil reference.")
			return
		}

		if c.Follow == LocalPlayer {
			CanSave = bitsClear(CanSave, isSequenceHappening)
		}

		dest = c.Follow.Position
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

	if !c.First || c.Mode == CameraModeLerp {
		t := c.Speed

		if c.Mode == CameraModeLerp {
			vd := raymath.Vector2Distance(dest, c.Position)

			if c.Progress > 1 {
				c.Progress = 1

				c.Mode = CameraModeStatic

				if c.Target != nil {
					c.Target.Trigger(c.Target, c)
				}
			}

			if t > vd {
				t = vd
			}

			c.Position = vector2Lerp(c.Start.Position, dest, c.Progress)
			c.Progress += t * dt
		} else {
			c.Position = vector2Lerp(c.Position, dest, t)
		}
	} else {
		c.Position = dest
	}

	c.Zoom = scalarLerp(c.Zoom, c.TargetZoom, c.ZoomSpeed*dt)

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
