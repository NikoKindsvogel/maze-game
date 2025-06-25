package game

import (
	"fmt"
	"math/rand"

	"maze-game/maze"
)

type Player struct {
	ID          string
	Row, Col    int
	Hurt        bool
	HasTreasure bool
	Bullet      bool
}

func PlacePlayers(m *maze.Maze, count int) []*Player {
	placed := make(map[[2]int]bool)
	players := make([]*Player, 0, count)

	for len(players) < count {
		r := rand.Intn(m.Size)
		c := rand.Intn(m.Size)
		pos := [2]int{r, c}

		if m.Grid[r][c].Type == maze.Empty &&
			!(r == m.TreasureRow && c == m.TreasureCol) &&
			!placed[pos] {

			player := &Player{
				ID:     fmt.Sprintf("P%d", len(players)+1),
				Row:    r,
				Col:    c,
				Hurt:   false,
				Bullet: true,
			}
			players = append(players, player)
			placed[pos] = true
		}
	}

	return players
}
