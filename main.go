package main

import (
	"maze-game/game"
	"maze-game/ui"
)

func main() {
	g := game.NewGame()
	ui.RunCLI(g)
}
