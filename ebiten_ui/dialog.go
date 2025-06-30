package ebiten_ui

import (
	"fmt"
	"image/color"
	"strings"

	"maze-game/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type DialogScreen struct {
	Game        *game.Game
	Input       string
	Messages    []string
	MaxMessages int
	Done        bool
}

func NewDialogScreen(size, holes, riverPush int, names []string) *DialogScreen {
	g := game.NewGameWithConfig(size, holes, riverPush, names)
	return &DialogScreen{
		Game:        g,
		Messages:    []string{"Game started. Use commands like: UP, DOWN, LEFT, RIGHT, SHOOT <dir>, SHOW, EXIT"},
		MaxMessages: 20,
	}
}

func (d *DialogScreen) Update() {
	for _, key := range ebiten.InputChars() {
		if key >= 32 && key <= 126 {
			d.Input += string(key)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(d.Input) > 0 {
		d.Input = d.Input[:len(d.Input)-1]
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		d.processCommand(d.Input)
		d.Input = ""
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
	} else if len(parts) == 2 && strings.ToUpper(parts[0]) == "SHOOT" {
		dir := strings.ToUpper(parts[1])
		if dir == "UP" || dir == "DOWN" || dir == "LEFT" || dir == "RIGHT" {
			result, _ = d.Game.PerformAction(fmt.Sprintf("SHOOT %s", dir))
		} else {
			d.appendMessage("Invalid direction for SHOOT.")
			return
		}
	} else {
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
	if len(d.Messages) > d.MaxMessages {
		d.Messages = d.Messages[len(d.Messages)-d.MaxMessages:]
	}
}

func (d *DialogScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})

	height := screen.Bounds().Dy()
	lineHeight := 16
	margin := 10

	// Prepare input and player turn info
	inputLine := "> " + d.Input
	player := d.Game.CurrentPlayer()
	turnInfo := fmt.Sprintf("%s's turn", player.ID)

	// Reserve bottom lines for input and turn display
	availableHeight := height - 2*margin - 2*lineHeight
	messageLines := availableHeight / lineHeight

	// Clip messages to most recent that fit
	start := 0
	if len(d.Messages) > messageLines {
		start = len(d.Messages) - messageLines
	}
	visible := d.Messages[start:]

	// Calculate starting Y from bottom up
	y := height - margin - lineHeight*2 // start above turn line

	// Draw messages from bottom up
	for i := len(visible) - 1; i >= 0; i-- {
		text.Draw(screen, visible[i], basicfont.Face7x13, margin, y, color.White)
		y -= lineHeight
	}

	// Draw turn info
	text.Draw(screen, turnInfo, basicfont.Face7x13, margin, height-margin-lineHeight, color.RGBA{200, 200, 0, 255})

	// Draw input line
	text.Draw(screen, inputLine, basicfont.Face7x13, margin, height-margin, color.White)
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
