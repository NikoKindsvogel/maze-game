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
	MoveHistory            []string
}

func NewGame() *Game {
	var size int = 7
	cfg := mazegen.MazeConfig{
		Size:                    size,
		NumHoles:                2,
		NumArmories:             1,
		NumHospitals:            1,
		NumDragons:              1,
		RiverLength:             size + 2,
		ExtraOpenings:           0,
		MinTreasureExitDistance: size - 2,
	}

	var m *maze.Maze
	var players []*Player

	for {
		m = mazegen.GenerateMaze(cfg)
		players = PlacePlayers(m, 2)

		if AllPlayersCanReachTreasureAndExit(m, players) && CanReachTreasureFromEstuary(m, m.TreasureRow, m.TreasureCol) && HospitalReachableFromExit(m) {
			break
		}
	}

	return &Game{
		Maze:                   m,
		Players:                players,
		current:                0,
		ShowVisibilityMessages: true,
		RiverMoveLength:        2,
	}
}

func NewGameWithConfig(size, holes, riverLength, riverPush int, names []string) *Game {
	cfg := mazegen.MazeConfig{
		Size:                    size,
		NumHoles:                holes,
		NumArmories:             1,
		NumHospitals:            1,
		NumDragons:              1,
		RiverLength:             riverLength,
		ExtraOpenings:           15,
		MinTreasureExitDistance: size - 2,
	}

	var m *maze.Maze
	var players []*Player

	for {
		m = mazegen.GenerateMaze(cfg)
		players = PlacePlayersByName(m, names)

		if AllPlayersCanReachTreasureAndExit(m, players) &&
			CanReachTreasureFromEstuary(m, m.TreasureRow, m.TreasureCol) &&
			HospitalReachableFromExit(m) {
			break
		}
	}

	return &Game{
		Maze:                   m,
		Players:                players,
		current:                0,
		ShowVisibilityMessages: true,
		RiverMoveLength:        riverPush,
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
		g.MoveHistory = append(g.MoveHistory, cmd)
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

	status := p.ID + ": "
	var treasure string

	// Handle wall
	if cell.Walls[dir] {
		// Hole teleportation
		switch cell.Type {
		case maze.Hole:
			g.teleportPlayerFromHole(p)
			status += "You hit a wall and got teleported through the hole! "
		case maze.River:
			status += "You hit a wall but the river pushes you! "
			status += g.moveAlongRiver(p)
		default:
			// Plain wall hit — return early with visibility from current position
			vis := g.computeVisibilityMessages(p)
			return status + "You hit a wall. " + strings.Join(vis, " ")
		}
	} else {
		// Valid move
		nr, nc := maze.Neighbor(p.Row, p.Col, dir)
		if !g.Maze.InBounds(nr, nc) {
			vis := g.computeVisibilityMessages(p)
			return status + "Out of bounds. " + strings.Join(vis, " ")
		}

		target := g.Maze.Grid[nr][nc]
		p.Row, p.Col = nr, nc

		switch target.Type {
		case maze.Exit:
			if p.HasTreasure && !p.Hurt {
				status += "You reached the exit with the treasure. You win!"
			} else if p.Hurt {
				status += "You reached the exit but you're hurt and can't escape. Go to a hospital first."
			} else {
				status += "You reached the exit but don't have the treasure."
			}
		case maze.Hole:
			g.teleportPlayerFromHole(p)
			status += "You fell into a hole and got teleported!"
		case maze.Dragon:
			if p.Hurt {
				status += "The dragon burned you. You're still hurt."
			} else {
				p.Hurt = true
				status += "The dragon burned you. You're hurt now."
			}
			if p.HasTreasure {
				g.Maze.TreasureRow = g.Maze.TreasureStartRow
				g.Maze.TreasureCol = g.Maze.TreasureStartCol
				g.Maze.TreasureOnMap = true
				p.HasTreasure = false
				status += " You lost the treasure and it was returned to its starting position."
			}
		case maze.Hospital:
			if p.Hurt {
				p.Hurt = false
				status += "You reached the hospital and are healed!"
			} else {
				status += "You visited the hospital, but you're already fine."
			}
		case maze.Armory:
			if p.Bullet {
				status += "You found an armory but already had a bullet!"
			} else {
				p.Bullet = true
				status += "You found an armory and received a bullet!"
			}
		case maze.River:
			status += "You stepped into a river. "
			status += g.moveAlongRiver(p)
		case maze.Estuary:
			status += "You stepped directly on the estuary."
		default:
			status += "You moved successfully."
		}

		// Check treasure
		if g.Maze.TreasureOnMap && nr == g.Maze.TreasureRow && nc == g.Maze.TreasureCol {
			if p.Hurt {
				treasure += " You found the treasure but can't pick it up because you are hurt!"
			} else {
				treasure += " You found the treasure!"
				p.HasTreasure = true
				g.Maze.TreasureOnMap = false
			}
		}
	}

	// Recompute visibility from final position
	visibilityMsgs := g.computeVisibilityMessages(p)

	if len(visibilityMsgs) > 0 {
		status += " " + strings.Join(visibilityMsgs, " ")
	}
	return status + treasure
}

func (g *Game) computeVisibilityMessages(p *Player) []string {
	var msgs []string

	// Dragon visibility
	for r := p.Row - 1; r >= 0; r-- {
		if g.Maze.Grid[r+1][p.Col].Walls[maze.Up] {
			break
		}
		if g.Maze.Grid[r][p.Col].Type == maze.Dragon {
			msgs = append(msgs, "The dragon sees you!")
			break
		}
	}
	for r := p.Row + 1; r < g.Maze.Size; r++ {
		if g.Maze.Grid[r-1][p.Col].Walls[maze.Down] {
			break
		}
		if g.Maze.Grid[r][p.Col].Type == maze.Dragon {
			msgs = append(msgs, "The dragon sees you!")
			break
		}
	}
	for c := p.Col - 1; c >= 0; c-- {
		if g.Maze.Grid[p.Row][c+1].Walls[maze.Left] {
			break
		}
		if g.Maze.Grid[p.Row][c].Type == maze.Dragon {
			msgs = append(msgs, "The dragon sees you!")
			break
		}
	}
	for c := p.Col + 1; c < g.Maze.Size; c++ {
		if g.Maze.Grid[p.Row][c-1].Walls[maze.Right] {
			break
		}
		if g.Maze.Grid[p.Row][c].Type == maze.Dragon {
			msgs = append(msgs, "The dragon sees you!")
			break
		}
	}

	if g.ShowVisibilityMessages && g.Maze.TreasureOnMap {
		if p.Col == g.Maze.TreasureCol {
			for r := p.Row - 1; r >= 0; r-- {
				if g.Maze.Grid[r+1][p.Col].Walls[maze.Up] {
					break
				}
				if r == g.Maze.TreasureRow {
					msgs = append(msgs, "You see the treasure!")
					break
				}
			}
			for r := p.Row + 1; r < g.Maze.Size; r++ {
				if g.Maze.Grid[r-1][p.Col].Walls[maze.Down] {
					break
				}
				if r == g.Maze.TreasureRow {
					msgs = append(msgs, "You see the treasure!")
					break
				}
			}
		}
		if p.Row == g.Maze.TreasureRow {
			for c := p.Col - 1; c >= 0; c-- {
				if g.Maze.Grid[p.Row][c+1].Walls[maze.Left] {
					break
				}
				if c == g.Maze.TreasureCol {
					msgs = append(msgs, "You see the treasure!")
					break
				}
			}
			for c := p.Col + 1; c < g.Maze.Size; c++ {
				if g.Maze.Grid[p.Row][c-1].Walls[maze.Right] {
					break
				}
				if c == g.Maze.TreasureCol {
					msgs = append(msgs, "You see the treasure!")
					break
				}
			}
		}
	}

	return msgs
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
					shooter.Bullet = false
					return fmt.Sprintf("You shot player %s! They are now hurt and dropped the treasure.", p.ID)
				} else {
					shooter.Bullet = false
					return fmt.Sprintf("You shot player %s! They are now hurt.", p.ID)
				}
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

func (g *Game) GetCurrent() int {
	return g.current
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

func (g *Game) Copy() *Game {
	// Deep copy players
	playersCopy := make([]*Player, len(g.Players))
	for i, p := range g.Players {
		playersCopy[i] = &Player{
			ID:     p.ID,
			Row:    p.Row,
			Col:    p.Col,
			Hurt:   p.Hurt,
			Bullet: p.Bullet,
		}
	}

	// Deep copy Maze
	newMaze := maze.CopyMaze(g.Maze)

	return &Game{
		Maze:                   newMaze,
		Players:                playersCopy,
		current:                g.current,
		ShowVisibilityMessages: false, // suppress output during sim
		RiverMoveLength:        g.RiverMoveLength,
	}
}
