package mazegen

import (
	"math/rand"
	"time"

	"maze-game/maze"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type MazeConfig struct {
	Size          int
	NumHoles      int
	NumArmories   int
	NumHospitals  int
	NumDragons    int
	RiverLength   int
	ExtraOpenings int
}

// GenerateMaze creates a solvable maze with given features.
func GenerateMaze(cfg MazeConfig) *maze.Maze {
	m := maze.CreateMaze(cfg.Size, 0, 0)
	carveMaze(m)
	openUpMaze(m, cfg.ExtraOpenings)

	placeRandomCellOfType(m, maze.Exit)
	placeTreasure(m)

	for i := 0; i < cfg.NumHoles; i++ {
		placeRandomCellOfType(m, maze.Hole)
	}
	for i := 0; i < cfg.NumHospitals; i++ {
		placeRandomCellOfType(m, maze.Hospital)
	}
	for i := 0; i < cfg.NumArmories; i++ {
		placeRandomCellOfType(m, maze.Armory)
	}
	for i := 0; i < cfg.NumDragons; i++ {
		placeRandomCellOfType(m, maze.Dragon)
	}

	placeSmartRiver(m, cfg.RiverLength)

	return m
}
