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

package system

import (
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
	vsCodeAsset := FindAsset(vsFileName)
	if vsCodeAsset == nil && vsFileName != "" {
		log.Fatalf("Shader %s was not found!\n", vsFileName)
		return Program{}
	}

	var vsCode string

	if vsFileName != "" {
		vsCode = string(vsCodeAsset.Data)
	}

	fsCodeAsset := FindAsset(fsFileName)
	if fsCodeAsset == nil {
		log.Fatalf("Shader %s was not found!\n", fsFileName)
		return Program{}
	}

	fsCode := string(fsCodeAsset.Data)
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
func (prog *Program) RenderToTexture(source, target RenderTarget) {
	rl.BeginTextureMode(target)
	rl.ClearBackground(rl.Black)
	rl.BeginShaderMode(prog.Shader)
	prog.UpdateDefaultUniforms()
	rl.DrawTexturePro(
		source.Texture,
		rl.NewRectangle(0, 0, float32(source.Texture.Width), float32(source.Texture.Height)),
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
	prog.SetShaderValue("size", []float32{float32(ScreenWidth), float32(ScreenHeight)}, 2)
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
