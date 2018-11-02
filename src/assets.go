package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	textures = make(map[string]rl.Texture2D)
)

// GetTexture retrieves a cached texture from disk
func GetTexture(texturePath string) rl.Texture2D {
	tx, ok := textures[texturePath]

	if ok {
		return tx
	}

	tx = rl.LoadTexture(texturePath)

	textures[texturePath] = tx

	return tx
}
