/*
 * @Author: V4 Games
 * @Date: 2018-12-10 21:28:01
 * @Last Modified by:   Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-12-10 21:28:01
 */

package core

// GameMode describes main game rules and subsystems
type GameMode interface {
	Init()
	Shutdown()
	Update()
	Draw()
	DrawUI()
	PostDraw()
	IgnoreUpdate() bool
	Serialize() string
	Deserialize(data string)
}
