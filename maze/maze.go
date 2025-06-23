package maze

type Cell struct {
	Row, Col int
	Walls    map[Direction]bool // true means wall exists in that direction
	Type     CellType
}

type Maze struct {
	Size          int
	Grid          [][]*Cell
	TreasureRow   int
	TreasureCol   int
	TreasureOnMap bool
}

// CreateMaze initializes an empty maze with border walls
func CreateMaze(size int, treasureRow, treasureCol int) *Maze {
	grid := make([][]*Cell, size)
	for r := 0; r < size; r++ {
		grid[r] = make([]*Cell, size)
		for c := 0; c < size; c++ {
			grid[r][c] = &Cell{
				Row:   r,
				Col:   c,
				Walls: map[Direction]bool{Up: false, Right: false, Down: false, Left: false},
				Type:  Empty,
			}
		}
	}

	m := &Maze{Grid: grid, Size: size}

	m.TreasureRow = treasureCol
	m.TreasureCol = treasureRow
	m.TreasureOnMap = true

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
