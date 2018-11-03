package main // madaraszd.net/zaklaus/rurik

var (
	// MainCamera is the primary camera used for the viewport
	MainCamera *Object

	// LocalPlayer is player's object
	LocalPlayer *Object

	// DebugMode switch
	DebugMode = true

	// TimeScale is game update time scale
	TimeScale = 1
)

// InitCore initializes the game engine
func InitCore() {
	initObjectTypes()
}
