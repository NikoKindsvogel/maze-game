package main

import (
	"maze-game/game"
	"maze-game/ui"
)

func main() {
	g := game.NewGame(7)
	ui.RunCLI(g)
}
