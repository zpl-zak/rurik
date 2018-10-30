package core

import (
	"fmt"

	ry "github.com/gen2brain/raylib-go/raymath"

	"../system"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/solarlune/GoAseprite"
	"github.com/solarlune/resolv/resolv"
)

type player struct {
	Texture rl.Texture2D
	Ase     goaseprite.File
	Locked  bool
}

// NewPlayer instance
func (p *Object) NewPlayer() {
	p.Ase = goaseprite.Load("assets/gfx/player.json")
	p.Texture = system.GetTexture("assets/gfx/player.png")
	p.Size = []int32{p.Ase.FrameWidth, p.Ase.FrameHeight}
	p.Update = updatePlayer
	p.Draw = drawPlayer
	p.GetAABB = getPlayerAABB
	p.HandleCollision = handlePlayerCollision

	LocalPlayer = p

	p.Ase.Play("Stand")
}

func updatePlayer(p *Object, dt float32) {
	p.Ase.Update(dt)

	var moveSpeed float32 = 120

	p.Movement.X = 0
	p.Movement.Y = 0

	if !p.Locked {
		p.Movement.X = system.GetAxis("horizontal")
		p.Movement.Y = system.GetAxis("vertical")
	}

	if ry.Vector2Length(p.Movement) > 0 {
		//ry.Vector2Normalize(&p.Movement)
		ry.Vector2Scale(&p.Movement, moveSpeed)

		p.Facing.X = p.Movement.X
		p.Facing.Y = p.Movement.Y
		ry.Vector2Normalize(&p.Facing)
	}

	p.Movement.X *= dt
	p.Movement.Y *= dt

	resX, okX := CheckForCollision(p, int32(p.Movement.X), 0)
	resY, okY := CheckForCollision(p, 0, int32(p.Movement.Y))

	if okX {
		p.Movement.X = float32(resX.ResolveX)
	}

	if okY {
		p.Movement.Y = float32(resY.ResolveY)
	}

	p.Position.X += p.Movement.X
	p.Position.Y += p.Movement.Y
}

func drawPlayer(p *Object) {
	sourceX, sourceY := p.Ase.GetFrameXY()
	source := rl.NewRectangle(float32(sourceX), float32(sourceY), float32(p.Ase.FrameWidth), float32(p.Ase.FrameHeight))

	dest := rl.NewRectangle(p.Position.X-float32(p.Ase.FrameWidth/2), p.Position.Y-float32(p.Ase.FrameHeight/2), float32(p.Ase.FrameWidth), float32(p.Ase.FrameHeight))

	rl.DrawTexturePro(p.Texture, source, dest, rl.Vector2{}, 0, rl.White)

	if DebugMode {
		c := getPlayerAABB(p)
		rl.DrawRectangleLinesEx(c.ToFloat32(), 1, rl.Blue)
		drawTextCentered(p.Name, c.X+c.Width/2, c.Y+c.Height+2, 1, rl.White)
	}
}

func getPlayerAABB(p *Object) rl.RectangleInt32 {
	return rl.RectangleInt32{
		X:      int32(p.Position.X) - int32(float32(p.Ase.FrameWidth/2)) + int32(float32(p.Ase.FrameWidth/4)),
		Y:      int32(p.Position.Y),
		Width:  p.Ase.FrameWidth / 2,
		Height: p.Ase.FrameHeight / 2,
	}
}

func handlePlayerCollision(res *resolv.Collision, p, other *Object) {
	fmt.Println("Collision has happened!")
}
