package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"maze-game/game"
	"maze-game/maze"
)

func RunCLI(g *game.Game) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Game started. Enter commands like: UP, DOWN, LEFT, RIGHT, SHOOT <direction> or SHOW")

	for {
		p := g.CurrentPlayer()
		status := ""
		if p.Hurt {
			status = " (hurt)"
		}
		fmt.Printf("%s's turn%s > ", p.ID, status)

		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)

		if len(parts) == 1 {
			cmd := strings.ToUpper(parts[0])
			if cmd == "SHOW" {
				ShowMap(g)
				continue
			}
			res, nextPlayer := g.PerformAction(cmd)
			fmt.Println(res)
			fmt.Printf("Next turn: %s\n", nextPlayer)
			continue
		}

		if len(parts) == 2 {
			cmd := strings.ToUpper(parts[0])
			dir := strings.ToUpper(parts[1])

			if cmd == "SHOOT" && (dir == "UP" || dir == "DOWN" || dir == "LEFT" || dir == "RIGHT") {
				res, nextPlayer := g.PerformAction(fmt.Sprintf("SHOOT %s", dir))
				fmt.Println(res)
				fmt.Printf("Next turn: %s\n", nextPlayer)
				continue
			}
		}

		fmt.Println("Invalid input. Use UP, DOWN, LEFT, RIGHT, SHOOT <direction> or SHOW")
	}
}

func ShowMap(g *game.Game) {
	m := g.GetMaze()
	players := g.GetPlayers()
	size := m.Size

	// Print top boundary
	fmt.Print(" ")
	for c := 0; c < size; c++ {
		fmt.Print(" ----")
	}
	fmt.Println()

	for r := 0; r < size; r++ {
		line := "|"
		bottomLine := "+"

		for c := 0; c < size; c++ {
			cell := m.Grid[r][c]

			// Player or cell char
			playerStr := ""
			for _, p := range players {
				if p.Row == r && p.Col == c {
					if p.Hurt {
						playerStr = strings.ToLower(p.ID) // e.g. "p1"
					} else {
						playerStr = p.ID // e.g. "P1"
					}
					break
				}
			}
			if playerStr == "" {
				playerStr = cellChar(cell.Type) + "  " // pad to 3 chars width
			} else {
				// Pad player string to 3 chars width (e.g. "P1 ")
				playerStr = fmt.Sprintf("%-3s", playerStr)
			}

			rightWall := " "
			if cell.Walls[maze.Right] {
				rightWall = "|"
			}

			line += " " + playerStr + rightWall

			bottomWall := "    "
			if cell.Walls[maze.Down] {
				bottomWall = "----"
			}
			bottomLine += bottomWall + "+"
		}

		fmt.Println(line)
		fmt.Println(bottomLine)
	}

	// Print player info
	fmt.Println("\nPlayers:")
	for _, p := range players {
		hurtStatus := "Healthy"
		if p.Hurt {
			hurtStatus = "Hurt"
		}
		treasureStatus := "No Treasure"
		if p.HasTreasure {
			treasureStatus = "Has Treasure"
		}
		bulletStatus := "No Bullet"
		if p.Bullet {
			bulletStatus = "Has Bullet"
		}
		fmt.Printf("- %s: %s, %s, %s\n", p.ID, hurtStatus, treasureStatus, bulletStatus)
	}
}

func cellChar(t maze.CellType) string {
	switch t {
	case maze.Treasure:
		return "T"
	case maze.Hospital:
		return "H"
	case maze.Exit:
		return "E"
	case maze.Dragon:
		return "D"
	case maze.Hole:
		return "O"
	case maze.Armory:
		return "A"
	default:
		return "."
	}
}
