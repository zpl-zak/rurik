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
	"encoding/gob"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/solarlune/resolv/resolv"
	tiled "github.com/zaklaus/go-tiled"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

var (
	// Worlds are all the worlds loaded within the game
	Worlds []*World

	objTypes    map[string]string
	worldIndex  int
	objCtors    = make(map[string]func(o *Object))
	drawObjects []*Object
)

// World represents the simulation region of the current map
type World struct {
	// Objects contains all spawned objects in the game
	Objects []*Object

	// GlobalIndex is globally tracked object allocation index
	GlobalIndex int
}

func (w *World) flushObjects() {
	w.Objects = []*Object{}
	w.GlobalIndex = 0
}

// GetObjectsOfType returns all objects of a given type
func (w *World) GetObjectsOfType(name string, avoidType bool) (ret []*Object) {
	for _, o := range w.Objects {
		if !avoidType && o.Class == name {
			ret = append(ret, o)
		} else if avoidType && o.Class != name {
			ret = append(ret, o)
		}
	}

	return ret
}

// NewObject creates a new object
func (w *World) NewObject(o *tiled.Object) *Object {
	if o == nil {
		o = &tiled.Object{}
	}

	return &Object{
		GID: func() int {
			idx := w.GlobalIndex
			w.GlobalIndex++
			return idx
		}(),
		world:   w,
		Name:    o.Name,
		Class:   o.Type,
		Visible: true,
		DebugVisible: func() bool {
			dbgVis := o.Properties.GetString("dbgShow") == "1"
			return dbgVis
		}(),
		Meta:         o,
		Depends:      []*Object{},
		Rotation:     float32(o.Rotation),
		IsPersistent: true,

		// Properties
		CollisionType:    o.Properties.GetString("colType"),
		AutoStart:        o.Properties.GetString("autostart") == "1",
		CanRepeat:        o.Properties.GetString("canRepeat") == "1",
		Fullbright:       o.Properties.GetString("fullbright") == "1",
		HasLight:         o.Properties.GetString("light") == "1",
		HasSpecularLight: o.Properties.GetString("specular") == "1",
		FileName:         o.Properties.GetString("file"),
		IsOverlay:        o.Properties.GetString("overlay") == "1",
		EventName:        o.Properties.GetString("event"),
		EventArgs:        CompileEventArgs(o.Properties.GetString("eventArgs")),
		TintColor:        GetColorFromProperty(o, "tint"),
		Color:            GetColorFromProperty(o, "color"),
		Attenuation:      GetFloatFromProperty(o, "atten"),
		Radius:           GetFloatFromProperty(o, "radius"),
		Offset:           GetVec2FromProperty(o, "offset"),
		PolyLines:        o.PolyLines,

		// Callbacks
		Finish:               func(o *Object) {},
		Init:                 func(o *Object) {},
		Update:               func(o *Object, dt float32) {},
		Trigger:              func(o, inst *Object) {},
		Draw:                 func(o *Object) {},
		DrawUI:               func(o *Object) {},
		HandleCollision:      func(res *resolv.Collision, o, other *Object) {},
		HandleCollisionEnter: func(res *resolv.Collision, o, other *Object) {},
		HandleCollisionLeave: func(res *resolv.Collision, o, other *Object) {},
		InsideArea:           func(o, a *Object) bool { return false },
		GetAABB:              func(o *Object) rl.RectangleInt32 { return rl.RectangleInt32{} },
		Serialize:            func(o *Object, enc *gob.Encoder) {},
		Deserialize:          func(o *Object, dec *gob.Decoder) {},
	}
}

// NewObjectPro creates a new object by providing its class dynamically.
func (w *World) NewObjectPro(name, class string) *Object {
	objectData := tiled.Object{
		Name: name,
		Type: class,
	}

	obj := w.spawnObject(&objectData)

	return obj
}

// AddObject adds object to the world
func (w *World) AddObject(o *Object) {
	if o == nil {
		return
	}

	if o.Name == "" {
		o.Name = fmt.Sprintf("unknown_%d", o.GID)
	}

	duplicateObject, _ := w.FindObject(o.Name)

	if duplicateObject == nil {
		w.Objects = append(w.Objects, o)
	} else {
		log.Printf("You can't add duplicate object to the world! Object name: %s\n", o.Name)
		return
	}
}

// FinalizeObject fully initializes the object and adds it to the world
func (w *World) FinalizeObject(o *Object) {
	if o == nil {
		return
	}

	w.AddObject(o)
	w.resolveObjectDependencies(o)
	w.findTargets(o)
	o.Finish(o)
}

func (w *World) spawnObject(objectData *tiled.Object) *Object {
	obj, err := BuildObject(w, objectData, nil)

	if err != nil {
		log.Printf("Object creation failed: %s!\n", err.Error())
		return nil
	}

	obj.Position = rl.NewVector2(float32(objectData.X), float32(objectData.Y))
	obj.Movement = rl.NewVector2(0, 0)

	if obj.CollisionType != "" {
		obj.IsCollidable = obj.CollisionType != "none"
	}

	return obj
}

func (w *World) postProcessObjects() {
	for _, o := range w.Objects {
		w.resolveObjectDependencies(o)
		w.findTargets(o)
		o.Finish(o)
	}
}

func (w *World) resolveObjectDependencies(o *Object) {
	depName := o.Meta.Properties.GetString("depends")

	if depName != "" {
		names := strings.Split(depName, ";")

		for _, x := range names {
			dep, _ := w.FindObject(x)

			if o == dep {
				log.Fatalf("Object depends on self: '%s' !\n", o.Name)
				return
			}

			o.Depends = append(o.Depends, dep)
		}
	}
}

func (w *World) findTargets(o *Object) {
	if o.ProxyName == "" {
		o.ProxyName = o.Meta.Properties.GetString("proxy")
	}

	if o.ProxyName != "" {
		o.Proxy, _ = w.FindObject(o.ProxyName)
	}
}

// FindObject looks up an object with specified name
func (w *World) FindObject(name string) (*Object, int) {
	for _, o := range w.Objects {
		if o.Name == name {
			return o, o.GID
		}
	}

	return nil, 0
}

func (w *World) getObject(gid int) *Object {
	for _, o := range w.Objects {
		if o.GID == gid {
			return o
		}
	}

	return nil
}

// UpdateObjects performs an update on all objects
func (w *World) UpdateObjects() {
	for _, o := range w.Objects {
		o.WasUpdated = false
	}

	for _, o := range w.Objects {
		w.updateObject(o, o)
	}
}

// InitObjects initializes all objects
func (w *World) InitObjects() {
	for _, o := range w.Objects {
		o.Init(o)
	}
}

func (w *World) updateObject(o, orig *Object) {
	if o.WasUpdated {
		return
	}

	if len(o.Depends) > 0 {
		for _, x := range o.Depends {
			if x == orig {
				log.Fatalf("Cyclic dependency on object update detected: '%s' !\n", orig.Name)
				return
			}

			if x == nil {
				log.Fatalf("Object '%s' depends on nil entity!\n", o.Name)
				return
			}

			w.updateObject(x, orig)
		}
	}

	o.updateTriggerArea()
	o.Update(o, system.FrameTime*float32(TimeScale))
	o.WasUpdated = true
}

// DrawObjects draws all drawable objects on the screen
// It sorts all objects by Y position
func (w *World) DrawObjects() {
	cullRenderProfiler.StartInvocation()
	drawObjects = []*Object{}

	for _, v := range w.Objects {
		if !v.Visible || v.IsOverlay {
			continue
		}

		rec := v.GetAABB(v)
		orig := v.Position
		orig.X += float32(rec.Width / 2.0)
		orig.Y += float32(rec.Height / 2.0)

		if cullingEnabled && !IsPointWithinFrustum(orig) {
			continue
		}

		drawObjects = append(drawObjects, v)
	}
	cullRenderProfiler.StopInvocation()

	sortRenderProfiler.StartInvocation()
	sort.Slice(drawObjects, func(i, j int) bool {
		return drawObjects[i].Position.Y < drawObjects[j].Position.Y
	})
	sortRenderProfiler.StopInvocation()

	for _, v := range w.Objects {
		if v.IsOverlay {
			drawObjects = append(drawObjects, v)
		}
	}

	for _, o := range drawObjects {
		o.Draw(o)

		if DebugMode && o.DebugVisible {
			rect := o.GetAABB(o)
			rl.DrawRectangleLines(
				rect.X,
				rect.Y,
				rect.Width,
				rect.Height,
				rl.RayWhite,
			)

			DrawTextCentered(o.Name, rect.X+rect.Width/2, rect.Y+rect.Height+2, 1, rl.White)
		}
	}
}

// DrawObjectUI draws all drawable objects's UI on the screen
func (w *World) DrawObjectUI() {
	for _, o := range w.Objects {
		if o.Visible {
			o.DrawUI(o)
		}
	}
}

// SetPosition sets the object's position
func (o *Object) SetPosition(x, y float32) {
	o.Position = rl.NewVector2(x, y)
}

// GetWorld returns the active world
func (o *Object) GetWorld() *World {
	return o.world
}
