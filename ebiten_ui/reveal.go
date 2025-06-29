package ebiten_ui

import (
	"maze-game/maze"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type RevealScreen struct {
	MazeStart *maze.Maze
	MazeEnd   *maze.Maze
}

func NewRevealScreen(start, end *maze.Maze) *RevealScreen {
	return &RevealScreen{
		MazeStart: start,
		MazeEnd:   end,
	}
}

func (r *RevealScreen) Update() {
	// Optional: support switching views
}

func (r *RevealScreen) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Maze Reveal Screen\n(You could show start + end maze here)")
}
