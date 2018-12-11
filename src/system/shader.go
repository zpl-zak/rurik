/*
 * @Author: V4 Games
 * @Date: 2018-12-10 21:58:20
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-11 04:02:17
 */

package system

import (
	"io/ioutil"
	"log"

	rl "github.com/zaklaus/raylib-go/raylib"
)

// Program describes shader pair
type Program struct {
	Shader   rl.Shader
	Uniforms map[string]int32
}

// ShaderPipeline describes the pipeline used for drawing things
type ShaderPipeline interface {
	Apply()
}

// NewProgram creates a new shader program
func NewProgram(vsFileName, fsFileName string) Program {
	vsCode, _ := ioutil.ReadFile(vsFileName)
	fsCode, err := ioutil.ReadFile(fsFileName)
	if err != nil {
		log.Fatalf("Shader %s was not found!\n", fsFileName)
		return Program{}
	}
	shader := rl.LoadShaderCode(string(vsCode), string(fsCode))

	return Program{
		Shader:   shader,
		Uniforms: make(map[string]int32),
	}
}

// NewProgramFromCode creates a new shader program from the code
func NewProgramFromCode(vsCode, fsCode string) Program {
	shader := rl.LoadShaderCode(vsCode, fsCode)

	return Program{
		Shader:   shader,
		Uniforms: make(map[string]int32),
	}
}

// RenderToTexture renders the current render target to an output
func (prog *Program) RenderToTexture(source, target *rl.RenderTexture2D) {
	if nullTarget == nil {
		nullTarget = CreateRenderTarget(ScreenWidth, ScreenHeight)
	}

	inp := source
	if source == nil {
		inp = nullTarget
	}

	prog.UpdateDefaultUniforms()
	rl.BeginTextureMode(*target)
	rl.BeginShaderMode(prog.Shader)
	rl.DrawTexturePro(
		inp.Texture,
		rl.NewRectangle(0, 0, float32(inp.Texture.Width), float32(inp.Texture.Height)),
		rl.NewRectangle(0, 0, float32(target.Texture.Width), float32(target.Texture.Height)),
		rl.Vector2{},
		0,
		rl.White,
	)
	rl.EndShaderMode()
	rl.EndTextureMode()
}

// UpdateDefaultUniforms updates the shader with default game engine values
func (prog *Program) UpdateDefaultUniforms() {
	prog.SetShaderValue("time", []float32{rl.GetTime()}, 1)
}

// SetShaderValue sets shader uniform value (float)
func (prog *Program) SetShaderValue(locName string, values []float32, count int32) {
	locID := prog.GetShaderLocation(locName)

	if locID != -1 {
		rl.SetShaderValue(prog.Shader, locID, values, count)
	}
}

// SetShaderValuei sets shader uniform value (int)
func (prog *Program) SetShaderValuei(locName string, values []int32, count int32) {
	locID := prog.GetShaderLocation(locName)

	if locID != -1 {
		rl.SetShaderValuei(prog.Shader, locID, values, count)
	}
}

// SetShaderValueMatrix sets shader uniform value (mat4x4)
func (prog *Program) SetShaderValueMatrix(locName string, mat rl.Matrix) {
	locID := prog.GetShaderLocation(locName)

	if locID != -1 {
		rl.SetShaderValueMatrix(prog.Shader, locID, mat)
	}
}

// GetShaderLocation retrieves the uniform location
func (prog *Program) GetShaderLocation(locName string) int32 {
	locID, ok := prog.Uniforms[locName]

	if !ok {
		locID = rl.GetShaderLocation(prog.Shader, locName)
	}

	prog.Uniforms[locName] = locID

	return locID
}
