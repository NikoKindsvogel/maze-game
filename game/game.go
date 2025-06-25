package game

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"maze-game/maze"
	"maze-game/mazegen"
)

type Game struct {
	Maze                   *maze.Maze
	Players                []*Player
	current                int
	ShowVisibilityMessages bool
	RiverMoveLength        int
}

func NewGame() *Game {
	var size int = 7
	cfg := mazegen.MazeConfig{
		Size:                    size,
		NumHoles:                2,
		NumArmories:             1,
		NumHospitals:            1,
		NumDragons:              1,
		RiverLength:             size - 1,
		ExtraOpenings:           size * 3,
		MinTreasureExitDistance: size - 2,
	}

	m := mazegen.GenerateMaze(cfg)

	players := PlacePlayers(m, 2)

	return &Game{
		Maze:                   m,
		Players:                players,
		current:                0,
		ShowVisibilityMessages: true,
		RiverMoveLength:        2,
	}
}

func (g *Game) PerformAction(cmd string) (string, string) {
	cmd = strings.ToUpper(cmd)

	var res string
	var usedTurn bool
	p := g.CurrentPlayer()

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

			// After a valid shot, check if standing on a hole
			if g.Maze.Grid[p.Row][p.Col].Type == maze.Hole {
				g.teleportPlayerFromHole(p)
				res += " You shot and got teleported through the hole!"
			}
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
		if cell.Type == maze.Hole {
			g.teleportPlayerFromHole(p)
			return p.ID + ": You hit a wall and got teleported through the hole!"
		}
		return p.ID + ": You hit a wall."
	}

	nr, nc := maze.Neighbor(p.Row, p.Col, dir)
	if !g.Maze.InBounds(nr, nc) {
		return p.ID + ": Out of bounds."
	}

	var visibilityMsgs []string

	target := g.Maze.Grid[nr][nc]

	// Move player
	p.Row, p.Col = nr, nc
	status := p.ID + ": "
	var treasure string

	// Check if the player moves onto the treasure
	if g.Maze.TreasureOnMap && nr == g.Maze.TreasureRow && nc == g.Maze.TreasureCol {
		if p.Hurt {
			treasure += " You found the treasure but can't pick it up because you are hurt!"
		} else {
			treasure += " You found the treasure!"
			p.HasTreasure = true
			g.Maze.TreasureOnMap = false
		}
	}

	switch target.Type {
	case maze.Exit:
		if p.HasTreasure && !p.Hurt {
			status += "You reached the exit with the treasure. You win!"
		} else if p.Hurt {
			status += "You're hurt and can't escape. Go to a hospital first."
		} else {
			status += "You reached the exit but don't have the treasure."
		}
	case maze.Hole:
		g.teleportPlayerFromHole(p)
		status += "You fell into a hole and got teleported!"
	case maze.Dragon:
		if p.HasTreasure {
			return p.ID + ": You stepped on the dragon while carrying the treasure. The dragon wins! Game over."
		}
		if p.HasTreasure {
			g.Maze.TreasureRow = p.Row
			g.Maze.TreasureCol = p.Col
			g.Maze.TreasureOnMap = true
			p.HasTreasure = false
		}
		p.Hurt = true
		status += "The dragon burned you. You're hurt and dropped the treasure."
	case maze.Hospital:
		if p.Hurt {
			p.Hurt = false
			status += "You reached the hospital and are healed!"
		} else {
			status += "You visited the hospital, but you're already fine."
		}
	case maze.Armory:
		p.Bullet = true
		status += "You found an armory and received a bullet!"
	case maze.River:
		status += "You stepped into a river. "
		status += g.moveAlongRiver(p)
	case maze.Estuary:
		status += "You stepped directly on the estuary."
	default:
		status += "You moved successfully."
	}

	// Dragon visibility (4 directions)
	for r := p.Row - 1; r >= 0; r-- {
		cell := g.Maze.Grid[r][p.Col]
		if g.Maze.Grid[r+1][p.Col].Walls[maze.Up] {
			break
		}
		if cell.Type == maze.Dragon {
			visibilityMsgs = append(visibilityMsgs, "The dragon sees you!")
			break
		}
	}
	for r := p.Row + 1; r < g.Maze.Size; r++ {
		cell := g.Maze.Grid[r][p.Col]
		if g.Maze.Grid[r-1][p.Col].Walls[maze.Down] {
			break
		}
		if cell.Type == maze.Dragon {
			visibilityMsgs = append(visibilityMsgs, "The dragon sees you!")
			break
		}
	}
	for c := p.Col - 1; c >= 0; c-- {
		cell := g.Maze.Grid[p.Row][c]
		if g.Maze.Grid[p.Row][c+1].Walls[maze.Left] {
			break
		}
		if cell.Type == maze.Dragon {
			visibilityMsgs = append(visibilityMsgs, "The dragon sees you!")
			break
		}
	}
	for c := p.Col + 1; c < g.Maze.Size; c++ {
		cell := g.Maze.Grid[p.Row][c]
		if g.Maze.Grid[p.Row][c-1].Walls[maze.Right] {
			break
		}
		if cell.Type == maze.Dragon {
			visibilityMsgs = append(visibilityMsgs, "The dragon sees you!")
			break
		}
	}

	if g.ShowVisibilityMessages {
		// Treasure visibility (if on map)
		if g.Maze.TreasureOnMap {
			for r := p.Row - 1; r >= 0; r-- {
				if g.Maze.Grid[r+1][p.Col].Walls[maze.Up] {
					break
				}
				if r == g.Maze.TreasureRow && p.Col == g.Maze.TreasureCol {
					visibilityMsgs = append(visibilityMsgs, "You see the treasure!")
					break
				}
			}
			for r := p.Row + 1; r < g.Maze.Size; r++ {
				if g.Maze.Grid[r-1][p.Col].Walls[maze.Down] {
					break
				}
				if r == g.Maze.TreasureRow && p.Col == g.Maze.TreasureCol {
					visibilityMsgs = append(visibilityMsgs, "You see the treasure!")
					break
				}
			}
			for c := p.Col - 1; c >= 0; c-- {
				if g.Maze.Grid[p.Row][c+1].Walls[maze.Left] {
					break
				}
				if p.Row == g.Maze.TreasureRow && c == g.Maze.TreasureCol {
					visibilityMsgs = append(visibilityMsgs, "You see the treasure!")
					break
				}
			}
			for c := p.Col + 1; c < g.Maze.Size; c++ {
				if g.Maze.Grid[p.Row][c-1].Walls[maze.Right] {
					break
				}
				if p.Row == g.Maze.TreasureRow && c == g.Maze.TreasureCol {
					visibilityMsgs = append(visibilityMsgs, "You see the treasure!")
					break
				}
			}
		}
	}

	if len(visibilityMsgs) > 0 {
		status += " " + strings.Join(visibilityMsgs, " ")
	}
	return status + treasure
}

func (g *Game) moveAlongRiver(p *Player) string {
	for i := 0; i < g.RiverMoveLength; i++ {
		cell := g.Maze.Grid[p.Row][p.Col]
		if cell.Type == maze.Estuary {
			return "You arrived at the estuary."
		}
		dir := cell.RiverDir
		if cell.Walls[dir] {
			break
		}
		nr, nc := maze.Neighbor(p.Row, p.Col, dir)
		if !g.Maze.InBounds(nr, nc) {
			break
		}
		p.Row, p.Col = nr, nc
		cell = g.Maze.Grid[p.Row][p.Col]
		if cell.Type == maze.Estuary {
			return "You arrived at the estuary."
		}
	}
	return ""
}
func (g *Game) teleportPlayerFromHole(p *Player) {
	size := g.Maze.Size
	r0, c0 := p.Row, p.Col

	// Collect all hole positions
	type pos struct{ r, c int }
	var holes []pos
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if g.Maze.Grid[r][c].Type == maze.Hole {
				holes = append(holes, pos{r, c})
			}
		}
	}

	// If less than 2 holes, do nothing
	if len(holes) < 2 {
		return
	}

	// Find index of current hole
	current := -1
	for i, h := range holes {
		if h.r == r0 && h.c == c0 {
			current = i
			break
		}
	}

	if current == -1 {
		return // not currently on a hole (shouldn't happen)
	}

	// Teleport to next hole (clockwise in list)
	next := (current + 1) % len(holes)
	p.Row, p.Col = holes[next].r, holes[next].c
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
			shooter.Bullet = false
			return "Your bullet hit a wall and stopped."
		}

		// Move to next cell in direction
		nr, nc := maze.Neighbor(r, c, dir)
		if !m.InBounds(nr, nc) {
			shooter.Bullet = false
			return "Your bullet flew out of bounds."
		}

		// Check if a player is in the next cell
		for _, p := range g.Players {
			if p.Row == nr && p.Col == nc {
				// Hit player
				p.Hurt = true

				if p.HasTreasure {
					p.HasTreasure = false
					m.TreasureRow = p.Row
					m.TreasureCol = p.Col
					m.TreasureOnMap = true
				}

				shooter.Bullet = false
				return fmt.Sprintf("You shot player %s! They are now hurt and dropped the treasure.", p.ID)
			}
		}

		// No player hit, continue to next cell
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

func (g *Game) SaveToFile(filename string) error {
	// Ensure the save directory exists
	if err := os.MkdirAll("saved", 0755); err != nil {
		return err
	}
	file, err := os.Create(filepath.Join("saved", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(g)
}

func LoadFromFile(filename string) (*Game, error) {
	file, err := os.Open(filepath.Join("saved", filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var g Game
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&g)
	return &g, err
}
