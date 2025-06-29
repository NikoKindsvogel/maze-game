package ebiten_ui

import (
	"maze-game/maze"
	"maze-game/mazegen"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameScreen struct {
	MazeStart *maze.Maze
	MazeEnd   *maze.Maze
	Finished  bool
}

func NewGameScreen() *GameScreen {
	cfg := mazegen.MazeConfig{
		Size:                    7,
		NumHoles:                2,
		NumArmories:             1,
		NumHospitals:            1,
		NumDragons:              1,
		RiverLength:             7 - 1,
		ExtraOpenings:           0,
		MinTreasureExitDistance: 7 - 2,
	}
	start := mazegen.GenerateMaze(cfg)
	return &GameScreen{
		MazeStart: start,
		MazeEnd:   maze.CopyMaze(start), // Simulate real-time game progression
	}
}

func (g *GameScreen) Update() {
	// Placeholder
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.Finished = true
	}
}

func (g *GameScreen) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Game Screen\nPress ESC to finish game and reveal maze...")
}
