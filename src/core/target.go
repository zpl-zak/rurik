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
	rl "github.com/zaklaus/raylib-go/raylib"
)

// NewTarget dummy target point
func (o *Object) NewTarget() {
	o.Draw = func(o *Object) {
		if !DebugMode || !o.DebugVisible {
			return
		}

		rl.DrawCircle(int32(o.Position.X), int32(o.Position.Y), 5, rl.White)
		DrawTextCentered(o.Name, int32(o.Position.X), int32(o.Position.Y)+5, 10, rl.White)
	}
}
