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
	"reflect"

	"github.com/solarlune/resolv/resolv"
	goaseprite "github.com/zaklaus/GoAseprite"
	tiled "github.com/zaklaus/go-tiled"
	rl "github.com/zaklaus/raylib-go/raylib"
)

// Object is a world object with logic and data
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
	EventName        string
	EventArgs        []string
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
	CanTrigger       bool
	IsPersistent     bool
	Fullbright       bool
	TintColor        rl.Color
	Color            rl.Color
	Attenuation      float32
	Radius           float32
	HasLight         bool
	HasSpecularLight bool
	IsOverlay        bool
	Offset           rl.Vector2
	LocalTileset     *tilesetData
	PolyLines        []*tiled.PolyLine
	UserData         ObjectUserData

	// Internal fields
	WasUpdated bool
	world      *World

	// Callbacks
	Init                 func(o *Object)
	Finish               func(o *Object)
	Update               func(o *Object, dt float32)
	Draw                 func(o *Object)
	DrawUI               func(o *Object)
	Trigger              func(o, inst *Object)
	HandleCollision      func(res *resolv.Collision, o, other *Object)
	HandleCollisionEnter func(res *resolv.Collision, o, other *Object)
	HandleCollisionLeave func(res *resolv.Collision, o, other *Object)
	InsideArea           func(o, a *Object) bool
	GetAABB              func(o *Object) rl.RectangleInt32
	Serialize            func(o *Object, enc *gob.Encoder)
	Deserialize          func(o *Object, dec *gob.Decoder)

	// Specialized data
	collision
	camera
	wait
	script
	anim
	area
	tile
}

// ObjectUserData describes custom data used by game's classes
type ObjectUserData interface {
	Serialize(enc *gob.Encoder)
	Deserialize(dec *gob.Decoder)
}

func initObjectTypes() {
	objTypes = map[string]string{
		"col":    "Collision",
		"cam":    "Camera",
		"target": "Target",
		"wait":   "Wait",
		"script": "Script",
		"anim":   "Anim",
		"area":   "Area",
		"tile":   "Tile",
	}
}

// RegisterClass adds a new object type
func RegisterClass(class string, ctor func(o *Object)) error {
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

	if o != nil && o.Type != "" {
		className = o.Type
	} else if savegameData != nil && savegameData.Type != "" {
		className = savegameData.Type
	}

	class, ok := objTypes[className]

	if !ok {
		// custom type check
		ctor, ctorOk := objCtors[className]

		if !ctorOk && className != "" {
			return nil, fmt.Errorf("class '%s' is undefined", className)
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
