package main

import (
	"fmt"

	"github.com/zaklaus/rurik/src/core"

	"github.com/solarlune/resolv/resolv"
	rl "github.com/zaklaus/raylib-go/raylib"
	ry "github.com/zaklaus/raylib-go/raymath"
	"github.com/zaklaus/rurik/src/system"
)

const ()

type player struct{}

// NewPlayer player
func NewPlayer(p *core.Object) {
	aseData := system.GetAnimData("player")
	p.Ase = &aseData
	p.Texture = system.GetTexture("player.png")
	p.Size = []int32{p.Ase.FrameWidth, p.Ase.FrameHeight}
	p.Update = updatePlayer
	p.Draw = drawPlayer
	p.GetAABB = core.GetSpriteAABB
	p.HandleCollision = handlePlayerCollision
	p.Facing = rl.NewVector2(1, 0)
	p.IsCollidable = true

	core.LocalPlayer = p

	core.PlayAnim(p, "StandE")
}

func updatePlayer(p *core.Object, dt float32) {
	p.Ase.Update(dt)

	var moveSpeed float32 = 120

	p.Movement.X = 0
	p.Movement.Y = 0

	if core.CanSave == 0 || core.BitsHas(core.CanSave, core.IsInChallenge) {
		p.Movement.X = system.GetAxis("horizontal")
		p.Movement.Y = system.GetAxis("vertical")
	}

	var tag string

	if ry.Vector2Length(p.Movement) > 0 {
		//ry.Vector2Normalize(&p.Movement)
		ry.Vector2Scale(&p.Movement, moveSpeed)

		p.Facing.X = p.Movement.X
		p.Facing.Y = p.Movement.Y
		ry.Vector2Normalize(&p.Facing)

		tag = "Walk"
	} else {
		tag = "Stand"
	}

	if p.Facing.Y > 0 {
		tag += "N"
	} else if p.Facing.Y < 0 {
		tag += "S"
	}

	if p.Facing.X > 0 {
		tag += "E"
	} else if p.Facing.X < 0 {
		tag += "W"
	}

	core.PlayAnim(p, tag)

	p.Movement.X *= dt
	p.Movement.Y *= dt

	resX, okX := core.CheckForCollision(p, int32(p.Movement.X), 0)
	resY, okY := core.CheckForCollision(p, 0, int32(p.Movement.Y))

	if okX {
		p.Movement.X = float32(resX.ResolveX)
	}

	if okY {
		p.Movement.Y = float32(resY.ResolveY)
	}

	p.Position.X += p.Movement.X
	p.Position.Y += p.Movement.Y
}

func drawPlayer(p *core.Object) {
	source := core.GetSpriteRectangle(p)
	dest := core.GetSpriteOrigin(p)

	if core.DebugMode && p.DebugVisible {
		c := core.GetSpriteAABB(p)
		rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
		core.DrawTextCentered(p.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
	}

	rl.DrawTexturePro(*p.Texture, source, dest, rl.Vector2{}, 0, core.SkyColor)
}

func handlePlayerCollision(res *resolv.Collision, p, other *core.Object) {
	fmt.Println("Collision has happened!")
}
