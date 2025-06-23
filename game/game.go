package game

import (
	"fmt"
	"strings"

	"maze-game/maze"
)

type Game struct {
	Maze    *maze.Maze
	Players []*Player
	current int // index of current player's turn
}

func NewGame(size int) *Game {
	m := maze.CreateMaze(size)

	// Example maze setup:
	m.Grid[2][2].Type = maze.Treasure
	m.Grid[0][4].Type = maze.Exit
	m.Grid[3][1].Type = maze.Hospital
	m.Grid[1][3].Type = maze.Dragon
	m.Grid[2][4].Type = maze.Hole

	// Internal walls
	m.AddWall(1, 1, maze.Right)
	m.AddWall(2, 3, maze.Down)

	players := []*Player{
		{ID: "P1", Row: 1, Col: 1, Hurt: false, Bullet: true},
		{ID: "P2", Row: 4, Col: 4, Hurt: false, Bullet: true},
	}

	return &Game{
		Maze:    m,
		Players: players,
		current: 0,
	}
}

func (g *Game) PerformAction(cmd string) (string, string) {
	cmd = strings.ToUpper(cmd)

	var res string
	var usedTurn bool

	switch {
	case cmd == "UP" || cmd == "DOWN" || cmd == "LEFT" || cmd == "RIGHT":
		res = g.moveCurrentPlayerInDirection(cmd)
		usedTurn = true
	case strings.HasPrefix(cmd, "SHOOT "):
		dir := strings.TrimPrefix(cmd, "SHOOT ")
		res = g.Shoot(dir)
		if res == "You have no bullets to shoot." || res == "Invalid shooting direction." {
			usedTurn = false
		} else {
			usedTurn = true
		}
	default:
		return "Unknown command.", g.CurrentPlayer().ID
	}

	if usedTurn {
		g.NextPlayer()
	}

	return res, g.CurrentPlayer().ID
}

func (g *Game) moveCurrentPlayerInDirection(dirStr string) string {
	dir := parseDirection(dirStr)
	if dir == -1 {
		return "Invalid direction."
	}

	p := g.CurrentPlayer()

	cell := g.Maze.Grid[p.Row][p.Col]
	if cell.Walls[dir] {
		return p.ID + ": You hit a wall."
	}

	nr, nc := maze.Neighbor(p.Row, p.Col, dir)
	if !g.Maze.InBounds(nr, nc) {
		return p.ID + ": Out of bounds."
	}

	target := g.Maze.Grid[nr][nc]
	p.Row, p.Col = nr, nc

	status := p.ID + ": "

	switch target.Type {
	case maze.Treasure:
		if p.Hurt {
			status += "You're hurt and can't pick up the treasure."
		} else {
			p.HasTreasure = true
			status += "You found the treasure!"
		}
	case maze.Exit:
		if p.HasTreasure && !p.Hurt {
			status += "You reached the exit with the treasure. You win!"
		} else if p.Hurt {
			status += "You're hurt and can't escape. Go to a hospital first."
		} else {
			status += "You reached the exit but don't have the treasure."
		}
	case maze.Hole:
		p.Hurt = true
		p.HasTreasure = false
		status += "You fell into a hole, got hurt, and dropped the treasure."
	case maze.Dragon:
		p.Hurt = true
		p.HasTreasure = false
		status += "The dragon burned you. You're hurt and dropped the treasure."
	case maze.Hospital:
		if p.Hurt {
			p.Hurt = false
			status += "You reached the hospital and are healed!"
		} else {
			status += "You visited the hospital, but you're already fine."
		}
	default:
		status += "You moved successfully."
	}

	return status
}

func (g *Game) Shoot(dirStr string) string {
	dir := parseDirection(dirStr)
	if dir == -1 {
		return "Invalid shooting direction."
	}

	shooter := g.CurrentPlayer()
	if !shooter.Bullet {
		return "You have no bullets to shoot."
	}

	r, c := shooter.Row, shooter.Col
	m := g.Maze

	for {
		// Check if wall blocks shooting out of current cell
		if m.Grid[r][c].Walls[dir] {
			return "Your bullet hit a wall and stopped."
		}

		// Move to next cell in direction
		nr, nc := maze.Neighbor(r, c, dir)
		if !m.InBounds(nr, nc) {
			return "Your bullet flew out of bounds."
		}

		// Check if a player is in next cell
		for _, p := range g.Players {
			if p.Row == nr && p.Col == nc {
				// Hit player
				p.Hurt = true
				p.HasTreasure = false
				shooter.Bullet = false
				return fmt.Sprintf("You shot player %s! They are now hurt and dropped the treasure.", p.ID)
			}
		}

		// No player hit, continue next cell
		r, c = nr, nc
	}
}

func parseDirection(input string) maze.Direction {
	switch strings.ToUpper(input) {
	case "UP":
		return maze.Up
	case "DOWN":
		return maze.Down
	case "LEFT":
		return maze.Left
	case "RIGHT":
		return maze.Right
	default:
		return -1
	}
}

func (g *Game) CurrentPlayer() *Player {
	return g.Players[g.current]
}

func (g *Game) NextPlayer() {
	g.current = (g.current + 1) % len(g.Players)
}

func (g *Game) GetMaze() *maze.Maze {
	return g.Maze
}

func (g *Game) GetPlayers() []*Player {
	return g.Players
}
