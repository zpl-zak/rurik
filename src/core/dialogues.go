/*
   Copyright 2019 V4 Games

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

	jsoniter "github.com/json-iterator/go"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

var (
	dialogues = make(map[string]Dialogue)
)

// Dialogue defines connversation flow
type Dialogue struct {
	Name       string `json:"name"`
	Avatar     *rl.Texture2D
	AvatarFile string    `json:"avatar"`
	Text       string    `json:"text"`
	Choices    []*Choice `json:"choices"`
	Target     string    `json:"target"`
	Event      string    `json:"event"`
	EventArgs  string    `json:"eventArgs"`
	SkipPrompt bool      `json:"skipPrompt"`
	Next       *Dialogue `json:"next"`
}

// Choice is a selection from dialogue branches
type Choice struct {
	Text string    `json:"text"`
	Next *Dialogue `json:"next"`
}

// InitText initializes the dialogue's text
func InitText(t *Dialogue) {
	if t.AvatarFile != "" {
		t.Avatar = system.GetTexture(t.AvatarFile)
	}

	if t.Next != nil {
		InitText(t.Next)
	}

	if t.Choices != nil {
		for _, ch := range t.Choices {
			if ch.Next != nil {
				InitText(ch.Next)
			}
		}
	}
}

// GetDialogue retrieves texts for a dialogue
func GetDialogue(name string) *Dialogue {
	dia, ok := dialogues[name]

	if ok {
		return &dia
	}

	data := system.GetFile(fmt.Sprintf("texts/%s", name), false)
	err := jsoniter.Unmarshal(data, &dia)

	if err != nil {
		log.Printf("Dialogue '%s' is broken!\n", name)
		return &Dialogue{}
	}

	dialogues[name] = dia
	return &dia
}
