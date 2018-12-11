/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:27:04
 * @Last Modified by: Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 03:38:51
 */

package core

import (
	"fmt"
	"log"

	rl "github.com/zaklaus/raylib-go/raylib"
	jsoniter "github.com/json-iterator/go"
	"madaraszd.net/zaklaus/rurik/src/system"
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
