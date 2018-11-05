package main

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	goaseprite "github.com/solarlune/GoAseprite"
	"github.com/solarlune/resolv/resolv"
)

var (
	// ObjectTypes contains all definitions of objects in the game
	ObjectTypes map[string]func(w *World, objectData *tiled.Object) *Object

	// Worlds are all the worlds loaded within the game
	Worlds []*World

	objTypes   map[string]string
	worldIndex int
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
	GID           int
	World         *World
	Name          string
	Class         string
	Visible       bool
	Position      rl.Vector2
	Movement      rl.Vector2
	Facing        rl.Vector2
	Size          []int32
	Meta          *tiled.Object
	Depends       []*Object
	Target        *Object
	FileName      string
	Texture       rl.Texture2D
	Ase           goaseprite.File
	LastTrigger   float32
	AutoStart     bool
	IsCollidable  bool
	CollisionType string
	Started       bool

	// Internal fields
	WasUpdated bool

	// Callbacks
	Finish          func(o *Object)
	Update          func(o *Object, dt float32)
	Draw            func(o *Object)
	DrawUI          func(o *Object)
	Trigger         func(o, inst *Object)
	HandleCollision func(res *resolv.Collision, o, other *Object)
	GetAABB         func(o *Object) rl.RectangleInt32

	// Specialized data
	player
	collision
	camera
	wait
	script
	talk
	anim
	area
}

func (w *World) flushObjects() {
	w.Objects = []*Object{}
	w.GlobalIndex = 0
}

func initObjectTypes() {
	ObjectTypes = make(map[string]func(w *World, objectData *tiled.Object) *Object)

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
	}

	for k := range objTypes {
		ObjectTypes[k] = func(w *World, o *tiled.Object) *Object {
			inst := w.NewObject(o)

			class := objTypes[o.Type]
			methodName := fmt.Sprintf("New%s", class)

			method := reflect.ValueOf(inst).MethodByName(methodName)

			if !method.IsValid() {
				log.Fatalf("Object type creation of '%s' not found!\n", class)
				return nil
			}

			method.Call([]reflect.Value{})

			return inst
		}
	}
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

	return &Object{
		GID:           idx,
		World:         w,
		Name:          o.Name,
		Class:         o.Type,
		Visible:       true,
		Meta:          o,
		Depends:       []*Object{},
		CollisionType: o.Properties.GetString("colType"),
		AutoStart:     o.Properties.GetString("autostart") == "1",
		FileName:      o.Properties.GetString("file"),

		// Callbacks
		Finish:          func(o *Object) {},
		Update:          func(o *Object, dt float32) {},
		Trigger:         func(o, inst *Object) {},
		Draw:            func(o *Object) {},
		DrawUI:          func(o *Object) {},
		HandleCollision: func(res *resolv.Collision, o, other *Object) {},
		GetAABB:         func(o *Object) rl.RectangleInt32 { return rl.RectangleInt32{} },
	}
}

// AddObject adds object to the world
func (w *World) AddObject(o *Object) {
	w.Objects = append(w.Objects, o)
}

func (w *World) spawnObject(objectData *tiled.Object) {
	objType, ok := ObjectTypes[objectData.Type]

	if !ok {
		log.Printf("Object type: %s not found!\n", objectData.Type)
		return
	}

	obj := objType(w, objectData)

	if obj == nil {
		log.Printf("Object creation failed!\n")
		return
	}

	obj.Position = rl.NewVector2(float32(objectData.X), float32(objectData.Y))
	obj.Movement = rl.NewVector2(0, 0)

	if obj.CollisionType != "" {
		obj.IsCollidable = obj.CollisionType != "none"
	}

	w.Objects = append(w.Objects, obj)
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

	o.Update(o, FrameTime*float32(TimeScale))
	o.WasUpdated = true
}

// DrawObjects draws all drawable objects on the screen
// It sorts all objects by Y position
func (w *World) DrawObjects() {
	sort.Slice(w.Objects, func(i, j int) bool {
		return w.Objects[i].Position.Y < w.Objects[j].Position.Y
	})

	for _, o := range w.Objects {
		if o.Visible {

			if !isPointWithinFrustum(o.Position) {
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
