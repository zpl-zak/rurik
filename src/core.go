package main

var (
	// MainCamera is the primary camera used for the viewport
	MainCamera *Object

	// LocalPlayer is player's object
	LocalPlayer *Object

	// DebugMode switch
	DebugMode = true
)

// Init initializes the game engine
func Init() {
	initObjectTypes()
}
