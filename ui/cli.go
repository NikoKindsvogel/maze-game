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

	fmt.Println("Game started. Enter commands like: UP, DOWN, LEFT, RIGHT, SHOOT <direction>, SHOW or EXIT")
	ShowMap(g)

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

		var res string
		var nextPlayer string

		if len(parts) == 1 {
			cmd := strings.ToUpper(parts[0])

			switch cmd {
			case "SAVE":
				fmt.Print("Enter filename to save: ")
				if !scanner.Scan() {
					break
				}
				saveFile := scanner.Text()
				if err := g.SaveToFile(saveFile); err != nil {
					fmt.Println("Error saving game:", err)
				} else {
					fmt.Println("Game saved.")
				}
				continue

			case "LOAD":
				fmt.Print("Enter filename to load: ")
				if !scanner.Scan() {
					break
				}
				loadFile := scanner.Text()
				newGame, err := game.LoadFromFile(loadFile)
				if err != nil {
					fmt.Println("Error loading game:", err)
				} else {
					*g = *newGame
					fmt.Println("Game loaded.")
					ShowMap(g)
				}
				continue
			case "REGEN":
				fmt.Print("Regenerating maze... ")
				g = game.NewGame()
				fmt.Println("Maze regenerated.")
				ShowMap(g)
				continue
			case "SHOW":
				ShowMap(g)
				continue
			case "EXIT":
				fmt.Println("Exiting game.")
				return
			case "SKIP":
				g.NextPlayer()
				continue
			case "UP", "DOWN", "LEFT", "RIGHT":
				res, nextPlayer = g.PerformAction(cmd)
			default:
				fmt.Println("Unknown command. Use UP, DOWN, LEFT, RIGHT, SHOW, SHOOT <direction> or EXIT")
				continue
			}
		} else if len(parts) == 2 {
			cmd := strings.ToUpper(parts[0])
			dir := strings.ToUpper(parts[1])

			if cmd == "SHOOT" && (dir == "UP" || dir == "DOWN" || dir == "LEFT" || dir == "RIGHT") {
				res, nextPlayer = g.PerformAction(fmt.Sprintf("SHOOT %s", dir))
			} else {
				fmt.Println("Invalid shoot command. Use SHOOT <UP|DOWN|LEFT|RIGHT>")
				continue
			}
		} else {
			fmt.Println("Invalid input. Use UP, DOWN, LEFT, RIGHT, SHOOT <direction>, SHOW or EXIT")
			continue
		}

		fmt.Println(res)

		// Check for game end conditions
		lowerRes := strings.ToLower(res)
		if strings.Contains(lowerRes, "game over") || strings.Contains(lowerRes, "you win") {
			fmt.Println("Game ended.")
			ShowMap(g)
			break
		}

		fmt.Printf("Next turn: %s\n", nextPlayer)
	}
}

func ShowMap(g *game.Game) {
	m := g.GetMaze()
	players := g.GetPlayers()
	size := m.Size

	// Print top boundary
	fmt.Print("+")
	for c := 0; c < size; c++ {
		fmt.Print("---+")
	}
	fmt.Println()

	for r := 0; r < size; r++ {
		line := "|"
		bottomLine := "+"

		for c := 0; c < size; c++ {
			cell := m.Grid[r][c]

			// Determine cell content: player or treasure or cell type
			cellChar := "   " // 3 spaces default

			// Check for player in cell
			for _, p := range players {
				if p.Row == r && p.Col == c {
					label := strings.ToUpper(p.ID)
					if p.Hurt {
						label = strings.ToLower(label)
					}
					// Make sure label is exactly 3 chars, padded or trimmed
					if len(label) > 3 {
						label = label[:3]
					} else {
						label = fmt.Sprintf("%-3s", label)
					}
					cellChar = label
					break
				}
			}

			// If no player and treasure is here
			if strings.TrimSpace(cellChar) == "" && m.TreasureOnMap && m.TreasureRow == r && m.TreasureCol == c {
				cellChar = " T "
			}

			// If still empty, use cell symbol (single char) centered in 3 spaces
			if strings.TrimSpace(cellChar) == "" {
				sym := cellSymbol(*cell)
				cellChar = fmt.Sprintf(" %s ", sym)
			}

			// Right wall
			rightWall := " "
			if cell.Walls[maze.Right] {
				rightWall = "|"
			}

			line += cellChar + rightWall

			// Bottom wall (3 dashes or spaces)
			bottomWall := "   "
			if cell.Walls[maze.Down] {
				bottomWall = "---"
			}
			bottomLine += bottomWall + "+"
		}

		fmt.Println(line)
		fmt.Println(bottomLine)
	}

	// Player info summary (unchanged)
	fmt.Println("\nPlayers:")
	for _, p := range players {
		id := p.ID
		status := ""
		if p.Hurt {
			status += " (hurt)"
		}
		if p.HasTreasure {
			status += " (has treasure)"
		}
		if p.Bullet {
			status += " (has bullet)"
		}
		fmt.Printf("- %s%s\n", id, status)
	}
}

// Helper for cell type
func cellSymbol(cell maze.Cell) string {
	switch cell.Type {
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
	case maze.River:
		switch cell.RiverDir {
		case maze.Up:
			return "↑"
		case maze.Down:
			return "↓"
		case maze.Left:
			return "←"
		case maze.Right:
			return "→"
		default:
			return "~" // fallback
		}
	case maze.Estuary:
		return "~"
	default:
		return "."
	}
}
