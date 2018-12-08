package core

// GameMode describes main game rules and subsystems.
type GameMode interface {
	Init()
	Shutdown()
	Reload()
	Update()
	Draw()
}
