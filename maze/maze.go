package maze

type Cell struct {
	Type     CellType
	Walls    map[Direction]bool
	RiverDir Direction // Only used if Type == River
}

type Maze struct {
	Size             int
	Grid             [][]*Cell
	TreasureRow      int
	TreasureCol      int
	TreasureOnMap    bool
	TreasureStartRow int
	TreasureStartCol int
}

// CreateMaze initializes an empty maze with border walls
func CreateMaze(size int, treasureRow, treasureCol int) *Maze {
	grid := make([][]*Cell, size)
	for r := 0; r < size; r++ {
		grid[r] = make([]*Cell, size)
		for c := 0; c < size; c++ {
			grid[r][c] = &Cell{
				Walls: map[Direction]bool{Up: true, Right: true, Down: true, Left: true},
				Type:  Empty,
			}
		}
	}

	m := &Maze{Grid: grid, Size: size}

	// Add border walls
	for r := 0; r < size; r++ {
		m.AddWall(r, 0, Left)
		m.AddWall(r, size-1, Right)
	}
	for c := 0; c < size; c++ {
		m.AddWall(0, c, Up)
		m.AddWall(size-1, c, Down)
	}

	return m
}

func (m *Maze) InBounds(r, c int) bool {
	return r >= 0 && r < m.Size && c >= 0 && c < m.Size
}

func (m *Maze) AddWall(r, c int, dir Direction) {
	if !m.InBounds(r, c) {
		return
	}

	m.Grid[r][c].Walls[dir] = true
	nr, nc := Neighbor(r, c, dir)
	if m.InBounds(nr, nc) {
		m.Grid[nr][nc].Walls[Opposite(dir)] = true
	}
}

func Neighbor(r, c int, dir Direction) (int, int) {
	switch dir {
	case Up:
		return r - 1, c
	case Down:
		return r + 1, c
	case Left:
		return r, c - 1
	case Right:
		return r, c + 1
	}
	return r, c
}

func Delta(dir Direction) (int, int) {
	switch dir {
	case Up:
		return -1, 0
	case Down:
		return 1, 0
	case Left:
		return 0, -1
	case Right:
		return 0, 1
	default:
		return 0, 0
	}
}

func DirectionFromDelta(dr, dc int) Direction {
	switch {
	case dr == -1 && dc == 0:
		return Up
	case dr == 1 && dc == 0:
		return Down
	case dr == 0 && dc == -1:
		return Left
	case dr == 0 && dc == 1:
		return Right
	}
	return Up // Default fallback
}

func (m *Maze) RemoveWallBetween(r, c int, dir Direction) {
	m.Grid[r][c].Walls[dir] = false
	nr, nc := Neighbor(r, c, dir)
	if m.InBounds(nr, nc) {
		m.Grid[nr][nc].Walls[Opposite(dir)] = false
	}
}

func FindExit(m *Maze) (row, col int, found bool) {
	for r := 0; r < m.Size; r++ {
		for c := 0; c < m.Size; c++ {
			if m.Grid[r][c].Type == Exit {
				return r, c, true
			}
		}
	}
	return 0, 0, false
}

func CopyMaze(original *Maze) *Maze {
	copyGrid := make([][]*Cell, original.Size)
	for r := 0; r < original.Size; r++ {
		copyGrid[r] = make([]*Cell, original.Size)
		for c := 0; c < original.Size; c++ {
			origCell := original.Grid[r][c]
			copyWalls := make(map[Direction]bool)
			for dir, hasWall := range origCell.Walls {
				copyWalls[dir] = hasWall
			}

			copyGrid[r][c] = &Cell{
				Type:     origCell.Type,
				Walls:    copyWalls,
				RiverDir: origCell.RiverDir,
			}
		}
	}

	return &Maze{
		Size:          original.Size,
		Grid:          copyGrid,
		TreasureRow:   original.TreasureRow,
		TreasureCol:   original.TreasureCol,
		TreasureOnMap: original.TreasureOnMap,
	}
}
