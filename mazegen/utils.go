package mazegen

import (
	"math/rand"

	"maze-game/maze"
)

func carveMaze(m *maze.Maze) {
	size := m.Size
	visited := make([][]bool, size)
	for i := range visited {
		visited[i] = make([]bool, size)
	}

	var dfs func(r, c int)
	dfs = func(r, c int) {
		visited[r][c] = true
		dirs := rand.Perm(4)
		for _, d := range dirs {
			dr, dc := maze.Delta(maze.Direction(d))
			nr, nc := r+dr, c+dc
			if m.InBounds(nr, nc) && !visited[nr][nc] {
				m.RemoveWallBetween(r, c, maze.Direction(d))
				dfs(nr, nc)
			}
		}
	}

	dfs(0, 0)
}

func placeRandomCellOfType(m *maze.Maze, t maze.CellType) {
	for {
		r := rand.Intn(m.Size)
		c := rand.Intn(m.Size)
		if m.Grid[r][c].Type == maze.Empty {
			m.Grid[r][c].Type = t
			return
		}
	}
}

func placeTreasure(m *maze.Maze, minDist int) {
	type point struct{ r, c int }

	var exit point
	found := false
	for r := 0; r < m.Size && !found; r++ {
		for c := 0; c < m.Size; c++ {
			if m.Grid[r][c].Type == maze.Exit {
				exit = point{r, c}
				found = true
				break
			}
		}
	}

	if !found {
		// fallback: no exit found, place treasure randomly
		for {
			r := rand.Intn(m.Size)
			c := rand.Intn(m.Size)
			if m.Grid[r][c].Type == maze.Empty {
				m.TreasureRow = r
				m.TreasureCol = c
				m.TreasureOnMap = true
				return
			}
		}
	}

	// Find a valid position at least minDist away from the exit
	for tries := 0; tries < 1000; tries++ {
		r := rand.Intn(m.Size)
		c := rand.Intn(m.Size)

		if m.Grid[r][c].Type == maze.Empty &&
			abs(r-exit.r)+abs(c-exit.c) >= minDist {

			m.TreasureRow = r
			m.TreasureCol = c
			m.TreasureOnMap = true
			return
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func openUpMaze(m *maze.Maze, extraOpenings int) {
	size := m.Size
	for i := 0; i < extraOpenings; {
		r := rand.Intn(size)
		c := rand.Intn(size)
		dirs := rand.Perm(4)
		for _, d := range dirs {
			dir := maze.Direction(d)
			nr, nc := maze.Neighbor(r, c, dir)
			if !m.InBounds(nr, nc) {
				continue
			}
			if m.Grid[r][c].Walls[dir] {
				m.RemoveWallBetween(r, c, dir)
				i++
				break
			}
		}
	}
}

func placeSmartRiver(m *maze.Maze, length int) {
	dirs := []maze.Direction{maze.Up, maze.Right, maze.Down, maze.Left}

	for attempt := 0; attempt < 100000; attempt++ {
		r := rand.Intn(m.Size)
		c := rand.Intn(m.Size)

		if m.Grid[r][c].Type != maze.Empty {
			continue
		}

		path := [][2]int{{r, c}}
		dir := dirs[rand.Intn(4)]

		for i := 1; i < length+2; i++ {
			// Occasionally change direction
			if rand.Float64() < 0.5 {
				dir = dirs[rand.Intn(4)]
			}

			dr, dc := maze.Delta(dir)
			nr, nc := r+dr, c+dc

			// Stop if out of bounds, blocked, or not empty
			if !m.InBounds(nr, nc) ||
				m.Grid[nr][nc].Type != maze.Empty ||
				m.Grid[r][c].Walls[dir] ||
				m.Grid[nr][nc].Walls[maze.Opposite(dir)] {
				break
			}

			path = append(path, [2]int{nr, nc})
			r, c = nr, nc
		}

		if len(path) < length {
			continue
		}

		// Now set directions and types
		for i := 0; i < len(path); i++ {
			r, c := path[i][0], path[i][1]
			if i == len(path)-1 {
				m.Grid[r][c].Type = maze.Estuary
			} else {
				m.Grid[r][c].Type = maze.River
				// Direction is from current to *next* tile
				nr, nc := path[i+1][0], path[i+1][1]
				m.Grid[r][c].RiverDir = maze.DirectionFromDelta(nr-r, nc-c)
			}
		}
		return
	}
}
