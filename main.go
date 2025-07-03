package main

import (
	"log"

	"maze-game/ebiten_ui"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := ebiten_ui.NewUIManager()
	ebiten.SetWindowSize(1200, 900)
	ebiten.SetWindowTitle("Maze Game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
