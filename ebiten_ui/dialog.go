package ebiten_ui

import (
	"fmt"
	"image/color"
	"maze-game/game"
	"maze-game/maze"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type DialogScreen struct {
	Game      *game.Game
	Input     string
	Messages  []string
	Done      bool
	startMaze maze.Maze

	KeyWasDown map[ebiten.Key]bool
}

func NewDialogScreen(size, holes, riverPush int, names []string) *DialogScreen {
	g := game.NewGameWithConfig(size, holes, riverPush, names)
	return &DialogScreen{
		Game:      g,
		Messages:  []string{"Game started. Use commands like: UP, DOWN, LEFT, RIGHT, SHOOT <dir>, SHOW, EXIT"},
		startMaze: *maze.CopyMaze(g.GetMaze()),
		KeyWasDown: map[ebiten.Key]bool{
			ebiten.KeyArrowUp:    false,
			ebiten.KeyArrowDown:  false,
			ebiten.KeyArrowLeft:  false,
			ebiten.KeyArrowRight: false,
			ebiten.KeyEnter:      false,
		},
	}
}

func (d *DialogScreen) Update() {
	// Handle text input
	for _, key := range ebiten.InputChars() {
		if key >= 32 && key <= 126 {
			d.Input += string(key)
		}
	}

	// Backspace
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(d.Input) > 0 {
		d.Input = d.Input[:len(d.Input)-1]
	}

	// Enter
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if !d.KeyWasDown[ebiten.KeyEnter] {
			d.processCommand(d.Input)
			d.Input = ""
		}
		d.KeyWasDown[ebiten.KeyEnter] = true
	} else {
		d.KeyWasDown[ebiten.KeyEnter] = false
	}

	// Arrow keys as movement commands
	d.checkArrowKey(ebiten.KeyArrowUp, "UP")
	d.checkArrowKey(ebiten.KeyArrowDown, "DOWN")
	d.checkArrowKey(ebiten.KeyArrowLeft, "LEFT")
	d.checkArrowKey(ebiten.KeyArrowRight, "RIGHT")
}

func (d *DialogScreen) checkArrowKey(key ebiten.Key, command string) {
	if ebiten.IsKeyPressed(key) {
		if !d.KeyWasDown[key] {
			d.processCommand(command)
		}
		d.KeyWasDown[key] = true
	} else {
		d.KeyWasDown[key] = false
	}
}

func (d *DialogScreen) processCommand(input string) {
	p := d.Game.CurrentPlayer()
	input = strings.TrimSpace(input)
	parts := strings.Fields(input)

	var result string

	if len(parts) == 1 {
		cmd := strings.ToUpper(parts[0])
		switch cmd {
		case "SHOW":
			d.appendMessage(d.renderMap())
			return
		case "EXIT":
			d.appendMessage("Game ended.")
			d.Done = true
			return
		case "SKIP":
			d.Game.NextPlayer()
			d.appendMessage("Turn skipped.")
			return
		case "UP", "DOWN", "LEFT", "RIGHT":
			result, _ = d.Game.PerformAction(cmd)
		default:
			d.appendMessage("Unknown command.")
			return
		}
	} else if len(parts) == 2 {
		cmd := strings.ToUpper(parts[0])
		arg := parts[1]

		switch cmd {
		case "SHOOT":
			dir := strings.ToUpper(arg)
			if dir == "UP" || dir == "DOWN" || dir == "LEFT" || dir == "RIGHT" {
				result, _ = d.Game.PerformAction(fmt.Sprintf("SHOOT %s", dir))
			} else {
				d.appendMessage("Invalid direction for SHOOT.")
				return
			}
		case "SAVE":
			err := d.Game.SaveToFile(arg)
			if err != nil {
				d.appendMessage("Error saving game: " + err.Error())
			} else {
				d.appendMessage("Game saved to " + arg)
			}
			return
		case "LOAD":
			newGame, err := game.LoadFromFile(arg)
			if err != nil {
				d.appendMessage("Error loading game: " + err.Error())
			} else {
				*d.Game = *newGame
				d.startMaze = *maze.CopyMaze(newGame.GetMaze())
				d.appendMessage("Game loaded from " + arg)
				d.appendMessage(d.renderMap())
			}
			return

		default:
			d.appendMessage("Unknown command.")
			return
		}
	} else {
		d.appendMessage("Invalid command.")
		return
	}

	d.appendMessage(fmt.Sprintf("%s's turn: %s", p.ID, input))
	d.appendMessage(result)

	if strings.Contains(strings.ToLower(result), "game over") {
		d.Done = true
	}
}

func (d *DialogScreen) appendMessage(msg string) {
	lines := strings.Split(msg, "\n")
	d.Messages = append(d.Messages, lines...)
}

func (d *DialogScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})

	height := screen.Bounds().Dy()
	lineHeight := 18
	margin := 40

	// Prepare input and turn display
	inputLine := "> " + d.Input
	player := d.Game.CurrentPlayer()
	turnInfo := fmt.Sprintf("%s's turn", player.ID)

	y := height - margin - 3*lineHeight

	// Draw messages bottom-up, but only if in-bounds
	for i := len(d.Messages) - 1; i >= 0; i-- {
		if y < margin {
			break // don’t draw past top margin
		}
		text.Draw(screen, d.Messages[i], basicfont.Face7x13, margin, y, color.White)
		y -= lineHeight
	}

	// Draw turn and input lines
	text.Draw(screen, turnInfo, basicfont.Face7x13, margin, height-margin-2*lineHeight, color.RGBA{200, 200, 0, 255})
	text.Draw(screen, inputLine, basicfont.Face7x13, margin, height-margin-lineHeight, color.White)
}

func (d *DialogScreen) renderMap() string {
	m := d.Game.GetMaze()
	players := d.Game.GetPlayers()
	size := m.Size
	var out strings.Builder

	// Top boundary
	out.WriteString("+")
	for c := 0; c < size; c++ {
		out.WriteString("---+")
	}
	out.WriteString("\n")

	for r := 0; r < size; r++ {
		line := "|"
		bottom := "+"

		for c := 0; c < size; c++ {
			cell := m.Grid[r][c]
			content := "   "

			for _, p := range players {
				if p.Row == r && p.Col == c {
					label := p.ID
					if p.Hurt {
						label = strings.ToLower(label)
					}
					if len(label) > 3 {
						label = label[:3]
					}
					content = fmt.Sprintf("%-3s", label)
					break
				}
			}

			if strings.TrimSpace(content) == "" && m.TreasureOnMap && m.TreasureRow == r && m.TreasureCol == c {
				content = " T "
			}

			if strings.TrimSpace(content) == "" {
				content = " . "
			}

			rightWall := " "
			if cell.Walls[1] {
				rightWall = "|"
			}
			line += content + rightWall

			bottomWall := "   "
			if cell.Walls[2] {
				bottomWall = "---"
			}
			bottom += bottomWall + "+"
		}

		out.WriteString(line + "\n")
		out.WriteString(bottom + "\n")
	}

	return out.String()
}
