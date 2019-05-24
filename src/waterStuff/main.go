package main

import (
	"fmt"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/raylib-go/raymath"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

const (
	screenW = 640
	screenH = 480

	// Apply 2x upscaling
	windowW = screenW * 2
	windowH = screenH * 2

	defaultCameraZoom = 1.79
)

func init() {
	rl.SetCallbackFunc(main)
}

type piece struct {
	temp float32
}

func main() {
	system.InitAssets([]string{"gfx.dta"}, true)

	rl.InitWindow(800, 800, "raylib [core] example - basic window")

	grid := []piece{}
	/* var blyat float32 = 0.0 */

	var worldSize int32 = 200
	var worldTileSize int32 = 4
	mousePickId := 0

	for idx := 0; idx < int(worldSize*worldSize); idx++ {
		grid = append(grid, piece{
			temp: 50, //rand.Float32() * 100,
		})
	}

	hot := rl.NewVector3(float32(rl.Red.R), float32(rl.Red.G), float32(rl.Red.B))
	cold := rl.NewVector3(float32(rl.Blue.R), float32(rl.Blue.G), float32(rl.Blue.B))

	var lastFrame int32
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		currentTime := rl.GetTime()

		var y int32 = 1
		var x int32 = 1
		for ; y < worldSize-1; y++ {
			for ; x < worldSize-1; x++ {
				idx := (y * (worldSize)) + x
				m := &grid[idx]
				col := raymath.Vector3Lerp(cold, hot, clamp(m.temp, 0, 100)/100)
				rl.DrawRectangle(
					int32(x*worldTileSize),
					int32(y*worldTileSize),
					worldTileSize,
					worldTileSize,
					rl.NewColor(uint8(col.X), uint8(col.Y), uint8(col.Z), 255),
				)

				ie := x - 1
				iw := x + 1
				is := y + 1
				in := y - 1
				/*
					if ie < 0 {
						ie = worldSize - 1
					}

					if iw >= worldSize {
						iw = 0
					}

					if is >= worldSize {
						is = 0
					}

					if in < 0 {
						in = worldSize - 1
					} */

				//fmt.Println(ie, iw, is, in)
				ve := &grid[(y*worldSize)+ie]
				vw := &grid[(y*worldSize)+iw]
				vs := &grid[(is*worldSize)+x]
				vn := &grid[(in*worldSize)+x]

				m.temp = core.ScalarLerp(
					m.temp,
					(ve.temp+vw.temp+vs.temp+vn.temp+m.temp)/5,
					1,
				)

				/* cyka := (ve.temp + vw.temp + vs.temp + vn.temp + m.temp) / 5
				/* if ve.temp > 0 {
					fmt.Println(ve.temp, vw.temp, vs.temp, vn.temp, m.temp, cyka)
				} */
				//m.temp = m.temp - (m.temp - cyka)
			}

			x = 1
		}

		mousePos := rl.GetMousePosition()
		mousePos.X = clamp(mousePos.X, 0, float32(worldSize)*float32(worldTileSize)-1) / float32(worldTileSize)
		mousePos.Y = clamp(mousePos.Y, 0, float32(worldSize)*float32(worldTileSize)-1) / float32(worldTileSize)

		//grid[mousePickId].temp = 50

		mousePickId = int((int32(mousePos.Y) * worldSize) + int32(mousePos.X))
		//fmt.Println(mousePos, mousePickId)

		if rl.IsMouseButtonDown(0) {
			grid[mousePickId].temp = 1000
		}

		if rl.IsMouseButtonDown(1) {
			grid[mousePickId].temp = -1000
		}

		/* for idx := 0; idx < int(worldSize*worldSize); idx++ {
			k := &grid[idx]
			k.temp = core.ScalarLerp(k.temp, rand.Float32()*100, 0.09)
		} */
		/*
			az := &grid[(34*worldSize)+worldSize/2+int32(blyat)]
			az.temp = 0
			bz := &grid[(78*worldSize)+worldSize/2+int32(blyat)]
			bz.temp = 100

			blyat += float32(math.Sin(float64(rl.GetFrameTime()))) * 400
		*/
		// hot water and cold as well cause fuck you that's why
		var ix int32
		for ix = 0; ix < worldSize; ix++ {
			grid[ix].temp = 100
			grid[((worldSize-1)*worldSize)+ix].temp = 0
		}

		if lastFrame%10 == 0 {
			fmt.Println((rl.GetTime() - currentTime) * 1000)
		}

		lastFrame++

		rl.EndDrawing()
		rl.PollInputEvents()
	}

	rl.CloseWindow()
}

func clamp(c, a, b float32) float32 {
	if c > b {
		c = b
	} else if c < a {
		c = a
	}

	return c
}
