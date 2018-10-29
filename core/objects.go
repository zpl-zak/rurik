package core

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/solarlune/resolv/resolv"

	"github.com/gen2brain/raylib-go/raylib"

	"github.com/lafriks/go-tiled"
)

var (
	// ObjectTypes contains all definitions of objects in the game
	ObjectTypes map[string]func(objectData *tiled.Object) *Object

	// Objects contains all spawned objects in the game
	Objects []*Object

	// GlobalIndex is globally tracked object allocation index
	GlobalIndex int

	objTypes map[string]string
)

// Object is map object with logic and data
type Object struct {
	GID        int
	Name       string
	Class      string
	Visible    bool
	Position   rl.Vector2
	Movement   rl.Vector2
	Facing     rl.Vector2
	Size       []int32
	Meta       *tiled.Object
	Depends    []*Object
	Target     *Object
	Started    bool
	WasUpdated bool

	// TODO: Improve this
	Finish          func(o *Object)
	Update          func(o *Object, dt float32)
	Draw            func(o *Object)
	DrawUI          func(o *Object)
	Trigger         func(o, inst *Object)
	HandleCollision func(res *resolv.Collision, o, other *Object)
	GetAABB         func(o *Object) rl.RectangleInt32

	// TODO: figure out better way to do this
	player
	collision
	camera
	wait
	script
	talk
}

func flushObjects() {
	Objects = []*Object{}
}

func initObjectTypes() {
	ObjectTypes = make(map[string]func(objectData *tiled.Object) *Object)
	flushObjects()

	objTypes = map[string]string{
		"player": "Player",
		"col":    "Collision",
		"cam":    "Camera",
		"target": "Target",
		"wait":   "Wait",
		"script": "Script",
		"talk":   "Talk",
	}

	for k := range objTypes {
		ObjectTypes[k] = func(o *tiled.Object) *Object {
			inst := NewObject(o)

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
func GetObjectsOfType(name string, avoidType bool) (ret []*Object) {
	for _, o := range Objects {
		if !avoidType && o.Class == name {
			ret = append(ret, o)
		} else if avoidType && o.Class != name {
			ret = append(ret, o)
		}
	}

	return ret
}

// NewObject creates a new object
func NewObject(o *tiled.Object) *Object {
	if o == nil {
		o = &tiled.Object{}
	}

	idx := GlobalIndex
	GlobalIndex++

	return &Object{
		GID:             idx,
		Name:            o.Name,
		Class:           o.Type,
		Visible:         true,
		Meta:            o,
		Depends:         nil,
		Finish:          func(o *Object) {},
		Update:          func(o *Object, dt float32) {},
		Trigger:         func(o, inst *Object) {},
		Draw:            func(o *Object) {},
		DrawUI:          func(o *Object) {},
		HandleCollision: func(res *resolv.Collision, o, other *Object) {},
		GetAABB:         func(o *Object) rl.RectangleInt32 { return rl.RectangleInt32{} },
	}
}

func spawnObject(objectData *tiled.Object) {
	objType, ok := ObjectTypes[objectData.Type]

	if !ok {
		log.Printf("Object type: %s not found!\n", objectData.Type)
		return
	}

	obj := objType(objectData)

	if obj == nil {
		log.Printf("Object creation failed!\n")
		return
	}

	obj.Position = rl.NewVector2(float32(objectData.X), float32(objectData.Y))
	obj.Movement = rl.NewVector2(0, 0)

	Objects = append(Objects, obj)
}

func postProcessObjects() {
	for _, o := range Objects {
		resolveObjectDependencies(o)
		findTargets(o)
		o.Finish(o)
	}
}

func resolveObjectDependencies(o *Object) {
	depName := o.Meta.Properties.GetString("depends")

	if depName != "" {
		names := strings.Split(depName, ";")
		o.Depends = []*Object{}

		for _, x := range names {
			dep, _ := FindObject(x)

			if o == dep {
				log.Fatalf("Object depends on self: '%s' !\n", o.Name)
				return
			}

			o.Depends = append(o.Depends, dep)
		}
	}
}

func findTargets(o *Object) {
	target := o.Meta.Properties.GetString("target")

	if target != "" {
		o.Target, _ = FindObject(target)
	}
}

// FindObject looks up an object with specified name
func FindObject(name string) (*Object, int) {
	for _, o := range Objects {
		if o.Name == name {
			return o, o.GID
		}
	}

	return nil, 0
}

func getObject(gid int) *Object {
	for _, o := range Objects {
		if o.GID == gid {
			return o
		}
	}

	return nil
}

// UpdateObjects performs an update on all objects
func UpdateObjects() {
	for _, o := range Objects {
		o.WasUpdated = false
	}

	for _, o := range Objects {
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

	o.Update(o, rl.GetFrameTime())
	o.WasUpdated = true
}

// DrawObjects draws all drawable objects on the screen
func DrawObjects() {
	for _, o := range Objects {
		if o.Visible {
			o.Draw(o)
		}
	}
}

// DrawObjectUI draws all drawable objects's UI on the screen
func DrawObjectUI() {
	for _, o := range Objects {
		if o.Visible {
			o.DrawUI(o)
		}
	}
}
