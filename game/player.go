package game

import (
	"fmt"
	"math/rand"
	"strings"

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

func CanReachUsingActions(g *Game, startRow, startCol, targetRow, targetCol int) bool {
	startState := &Game{
		Maze: g.Maze,
		Players: []*Player{{
			ID:     g.CurrentPlayer().ID,
			Row:    startRow,
			Col:    startCol,
			Hurt:   g.CurrentPlayer().Hurt,
			Bullet: g.CurrentPlayer().Bullet,
		}},
		current:                0,
		ShowVisibilityMessages: false,
		RiverMoveLength:        g.RiverMoveLength,
	}

	stack := []*Game{startState}
	visited := make(map[[2]int]bool)

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		p := current.CurrentPlayer()
		pos := [2]int{p.Row, p.Col}

		if pos == [2]int{targetRow, targetCol} {
			return true
		}

		if visited[pos] {
			continue
		}
		visited[pos] = true

		for _, dir := range []string{"UP", "DOWN", "LEFT", "RIGHT"} {
			nextGame := current.Copy()
			res, _ := nextGame.PerformAction(dir)

			// Skip if move was invalid
			if strings.Contains(res, "can't") || strings.Contains(strings.ToLower(res), "invalid") {
				continue
			}

			np := nextGame.CurrentPlayer()
			cell := nextGame.Maze.Grid[np.Row][np.Col]

			// Skip if new cell is a dragon
			if cell.Type == maze.Dragon {
				continue
			}

			stack = append(stack, nextGame)
		}
	}

	return false
}

func AllPlayersCanReachTreasureAndExit(m *maze.Maze, players []*Player) bool {
	// Save original treasure position
	treasureRow, treasureCol := m.TreasureRow, m.TreasureCol

	// Find exit position
	exitRow, exitCol, found := maze.FindExit(m)
	if !found {
		return false
	}

	for _, p := range players {
		// Create a temporary game instance with a deep copy of the player
		tempPlayer := &Player{
			ID:     p.ID,
			Row:    p.Row,
			Col:    p.Col,
			Hurt:   p.Hurt,
			Bullet: p.Bullet,
		}
		tempGame := &Game{
			Maze:                   m,
			Players:                []*Player{tempPlayer},
			current:                0,
			ShowVisibilityMessages: false,
			RiverMoveLength:        2,
		}

		// Step 1: Check if the player can reach the treasure
		if !CanReachUsingActions(tempGame, tempPlayer.Row, tempPlayer.Col, treasureRow, treasureCol) {
			return false
		}

		// Move player to treasure and restore treasure position
		tempPlayer.Row = treasureRow
		tempPlayer.Col = treasureCol
		m.TreasureRow = treasureRow
		m.TreasureCol = treasureCol
		m.TreasureOnMap = true

		// Step 2: Check if the player can reach the exit from the treasure
		if !CanReachUsingActions(tempGame, treasureRow, treasureCol, exitRow, exitCol) {
			return false
		}

		m.TreasureRow = treasureRow
		m.TreasureCol = treasureCol
		m.TreasureOnMap = true
	}

	return true
}

func CanReachTreasureFromEstuary(m *maze.Maze, treasureRow, treasureCol int) bool {
	// Find the estuary tile
	var estuaryRow, estuaryCol int
	found := false
	for r := 0; r < m.Size; r++ {
		for c := 0; c < m.Size; c++ {
			if m.Grid[r][c].Type == maze.Estuary {
				estuaryRow = r
				estuaryCol = c
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return false
	}

	// Create a game with a single player starting on the estuary
	player := &Player{
		ID:     "RiverTester",
		Row:    estuaryRow,
		Col:    estuaryCol,
		Hurt:   false,
		Bullet: true,
	}
	g := &Game{
		Maze:                   m,
		Players:                []*Player{player},
		current:                0,
		ShowVisibilityMessages: false,
		RiverMoveLength:        2,
	}

	if CanReachUsingActions(g, estuaryRow, estuaryCol, treasureRow, treasureCol) {
		m.TreasureRow = treasureRow
		m.TreasureCol = treasureCol
		m.TreasureOnMap = true
		return true
	} else {
		return false
	}
}

func HospitalReachableFromExit(m *maze.Maze) bool {
	treasureRow := m.TreasureRow
	treasureCol := m.TreasureCol
	exitRow, exitCol, exitFound := maze.FindExit(m)
	if !exitFound {
		return false
	}

	// Find the hospital
	var hospitalRow, hospitalCol int
	hospitalFound := false
	for r := 0; r < m.Size; r++ {
		for c := 0; c < m.Size; c++ {
			if m.Grid[r][c].Type == maze.Hospital {
				hospitalRow = r
				hospitalCol = c
				hospitalFound = true
				break
			}
		}
		if hospitalFound {
			break
		}
	}
	if !hospitalFound {
		return false
	}

	// Create a test player starting at exit
	exitPlayer := &Player{
		ID:     "HospitalTester",
		Row:    exitRow,
		Col:    exitCol,
		Hurt:   false,
		Bullet: true,
	}
	gameFromExit := &Game{
		Maze:                   m,
		Players:                []*Player{exitPlayer},
		current:                0,
		ShowVisibilityMessages: false,
		RiverMoveLength:        2,
	}

	// Check exit -> hospital
	if !CanReachUsingActions(gameFromExit, exitRow, exitCol, hospitalRow, hospitalCol) {
		return false
	}

	// Now check hospital -> exit
	hospitalPlayer := &Player{
		ID:     "HospitalTester2",
		Row:    hospitalRow,
		Col:    hospitalCol,
		Hurt:   false,
		Bullet: true,
	}
	gameFromHospital := &Game{
		Maze:                   m,
		Players:                []*Player{hospitalPlayer},
		current:                0,
		ShowVisibilityMessages: false,
		RiverMoveLength:        2,
	}

	if !CanReachUsingActions(gameFromHospital, hospitalRow, hospitalCol, exitRow, exitCol) {
		m.TreasureRow = treasureRow
		m.TreasureCol = treasureCol
		m.TreasureOnMap = true
		return false
	}

	return true
}
