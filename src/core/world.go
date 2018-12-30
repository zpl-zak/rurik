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
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	goaseprite "github.com/solarlune/GoAseprite"
	"github.com/solarlune/resolv/resolv"
	tiled "github.com/zaklaus/go-tiled"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

var (
	// Worlds are all the worlds loaded within the game
	Worlds []*World

	objTypes   map[string]string
	worldIndex int
	objCtors   = make(map[string]func(o *Object))
)

// World represents the simulation region of the current map
type World struct {
	// Objects contains all spawned objects in the game
	Objects []*Object

	// GlobalIndex is globally tracked object allocation index
	GlobalIndex int
}

// Object is map object with logic and data
type Object struct {
	GID              int
	Name             string
	Class            string
	Visible          bool
	DebugVisible     bool
	Position         rl.Vector2
	Movement         rl.Vector2
	Rotation         float32
	Facing           rl.Vector2
	Size             []int32
	Meta             *tiled.Object
	Depends          []*Object
	Target           *Object
	Proxy            *Object
	ProxyName        string
	FileName         string
	Texture          *rl.Texture2D
	Ase              *goaseprite.File
	LastTrigger      float32
	AutoStart        bool
	IsCollidable     bool
	CollisionType    string
	Started          bool
	WasExecuted      bool
	CanRepeat        bool
	IsPersistent     bool
	Fullbright       bool
	TintColor        rl.Color
	Color            rl.Color
	Attenuation      float32
	Radius           float32
	HasLight         bool
	HasSpecularLight bool
	Offset           rl.Vector2
	LocalTileset     *tilesetData
	UserData         ObjectUserData

	// Internal fields
	WasUpdated bool
	world      *World

	// Callbacks
	Init            func(o *Object)
	Finish          func(o *Object)
	Update          func(o *Object, dt float32)
	Draw            func(o *Object)
	DrawUI          func(o *Object)
	Trigger         func(o, inst *Object)
	HandleCollision func(res *resolv.Collision, o, other *Object)
	GetAABB         func(o *Object) rl.RectangleInt32
	Serialize       func(o *Object) string
	Deserialize     func(o *Object, data string)

	// Specialized data
	player
	collision
	camera
	wait
	script
	talk
	anim
	area
	tile
}

// ObjectUserData describes custom data used by game's classes
type ObjectUserData interface {
	Serialize() string
	Deserialize(input string)
}

func (w *World) flushObjects() {
	w.Objects = []*Object{}
	w.GlobalIndex = 0
}

func initObjectTypes() {
	objTypes = map[string]string{
		"player": "Player",
		"col":    "Collision",
		"cam":    "Camera",
		"target": "Target",
		"wait":   "Wait",
		"script": "Script",
		"talk":   "Talk",
		"anim":   "Anim",
		"area":   "Area",
		"tile":   "Tile",
	}
}

// RegisterClass adds a new object type
func RegisterClass(class, methodName string, ctor func(o *Object)) error {
	_, ok := objCtors[class]

	if ok {
		return fmt.Errorf("can't register already existing class '%s'", class)
	}

	objCtors[class] = ctor

	return nil
}

// BuildObject builds already-prepared object
func BuildObject(w *World, o *tiled.Object, savegameData *defaultObjectData) (*Object, error) {
	inst := w.NewObject(o)

	if inst.Name == "" && inst.Class == "" {
		inst.Name = fmt.Sprintf("tile_%d", inst.GID)
		inst.Class = "tile"
		o.Type = "tile"
	}

	if inst.Class == "col" && inst.Name == "" {
		inst.Name = fmt.Sprintf("col_%d", inst.GID)
	}

	className := "Unknown"

	if o != nil {
		className = o.Type
	} else if savegameData != nil {
		className = savegameData.Type
	}

	class, ok := objTypes[className]

	if !ok {
		// custom type check
		ctor, ctorOk := objCtors[className]

		if !ctorOk {
			return nil, fmt.Errorf("class %s is undefined", class)
		}

		ctor(inst)
		return inst, nil
	}

	methodName := fmt.Sprintf("New%s", class)

	method := reflect.ValueOf(inst).MethodByName(methodName)

	if !method.IsValid() {
		return nil, fmt.Errorf("internal error")
	}

	method.Call([]reflect.Value{})

	return inst, nil
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

	idx := w.GlobalIndex
	w.GlobalIndex++

	tintHex := o.Properties.GetString("tint")
	var tint rl.Color

	if tintHex != "" {
		tintVec, _ := getColorFromHex(tintHex)
		tint = vec3ToColor(tintVec)
	} else {
		tint = rl.Blank
	}

	hexColor := o.Properties.GetString("color")
	var color rl.Color

	if hexColor != "" {
		col, _ := getColorFromHex(hexColor)
		color = vec3ToColor(col)
	}

	txtRadius := o.Properties.GetString("radius")
	var radius float32

	if txtRadius != "" {
		rad, _ := strconv.ParseFloat(txtRadius, 32)
		radius = float32(rad)
	}

	txtAttenuation := o.Properties.GetString("atten")
	var attenuation float32

	if txtAttenuation != "" {
		atten, _ := strconv.ParseFloat(txtAttenuation, 32)
		attenuation = float32(atten)
	}

	txtOffset := o.Properties.GetString("offset")
	var offset rl.Vector2

	if txtOffset != "" {
		offset = stringToVec2(txtOffset)
	}

	return &Object{
		GID:              idx,
		world:            w,
		Name:             o.Name,
		Class:            o.Type,
		Visible:          true,
		DebugVisible:     DefaultDebugShowAll,
		Meta:             o,
		Depends:          []*Object{},
		Rotation:         float32(o.Rotation),
		CollisionType:    o.Properties.GetString("colType"),
		AutoStart:        o.Properties.GetString("autostart") == "1",
		CanRepeat:        o.Properties.GetString("canRepeat") == "1",
		Fullbright:       o.Properties.GetString("fullbright") == "1",
		HasLight:         o.Properties.GetString("light") == "1",
		HasSpecularLight: o.Properties.GetString("specular") == "1",
		FileName:         o.Properties.GetString("file"),
		TintColor:        tint,
		Color:            color,
		Attenuation:      attenuation,
		Radius:           radius,
		Offset:           offset,
		IsPersistent:     true,

		// Callbacks
		Finish:          func(o *Object) {},
		Init:            func(o *Object) {},
		Update:          func(o *Object, dt float32) {},
		Trigger:         func(o, inst *Object) {},
		Draw:            func(o *Object) {},
		DrawUI:          func(o *Object) {},
		HandleCollision: func(res *resolv.Collision, o, other *Object) {},
		GetAABB:         func(o *Object) rl.RectangleInt32 { return rl.RectangleInt32{} },
		Serialize:       func(o *Object) string { return "{}" },
		Deserialize:     func(o *Object, data string) {},
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
	target := o.Meta.Properties.GetString("target")

	if target != "" {
		o.Target, _ = w.FindObject(target)
	}

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
		updateObject(o, o)
	}
}

// InitObjects initializes all objects
func (w *World) InitObjects() {
	for _, o := range w.Objects {
		o.Init(o)
	}
}

func updateObject(o, orig *Object) {
	if o.WasUpdated {
		return
	}

	if len(o.Depends) > 0 {
		for _, x := range o.Depends {
			if x == orig {
				log.Fatalf("Cyclic dependency on object update detected: '%s' !\n", orig.Name)
				return
			}

			updateObject(x, orig)
		}
	}

	o.Update(o, system.FrameTime*float32(TimeScale))
	o.WasUpdated = true
}

// DrawObjects draws all drawable objects on the screen
// It sorts all objects by Y position
func (w *World) DrawObjects() {
	sortRenderProfiler.StartInvocation()
	drawObjects := make([]*Object, len(w.Objects))
	copy(drawObjects, w.Objects)
	sort.Slice(drawObjects, func(i, j int) bool {
		return drawObjects[i].Position.Y < drawObjects[j].Position.Y
	})
	sortRenderProfiler.StopInvocation()

	for _, o := range drawObjects {
		if o.Visible {
			rec := o.GetAABB(o)
			orig := o.Position
			orig.X += float32(rec.Width / 2.0)
			orig.Y += float32(rec.Height / 2.0)

			if !isPointWithinFrustum(orig) && cullingEnabled {
				continue
			}

			o.Draw(o)
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
