package ebiten_ui

import (
	"fmt"
	"image/color"
	"maze-game/game"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	dialogExitButtonX = screenWidth - 130
	dialogExitButtonY = screenHeight - 150
)

type DialogScreen struct {
	Game      *game.Game
	Input     string
	Messages  []string
	Done      bool
	startGame game.Game

	KeyWasDown map[ebiten.Key]bool

	Background *ebiten.Image
	ExitButton *ebiten.Image
}

func NewDialogScreen(size, holes, riverLength, riverPush int, names []string) *DialogScreen {
	g := game.NewGameWithConfig(size, holes, riverLength, riverPush, names)
	bgImage := loadImageFromEmbed("backgrounds/background.png")
	exitImage := loadImageFromEmbed("buttons/dialog_button_exit.png")
	return &DialogScreen{
		Game:      g,
		Messages:  []string{"Game started. Use commands like: UP, DOWN, LEFT, RIGHT, SHOOT <dir>, EXIT"},
		startGame: *g.Copy(),
		KeyWasDown: map[ebiten.Key]bool{
			ebiten.KeyArrowUp:    false,
			ebiten.KeyArrowDown:  false,
			ebiten.KeyArrowLeft:  false,
			ebiten.KeyArrowRight: false,
			ebiten.KeyEnter:      false,
		},
		Background: bgImage,
		ExitButton: exitImage,
	}
}

func (d *DialogScreen) Update(u *UIManager) {
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

	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	x, y := ebiten.CursorPosition()

	if mouseDown && !u.mouseWasDown {
		if x >= dialogExitButtonX && x <= dialogExitButtonX+sideButtonWidth &&
			y >= dialogExitButtonY && y <= dialogExitButtonY+sideButtonHeight {
			d.Done = true
		}
	}

	u.mouseWasDown = mouseDown
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
				d.startGame = *newGame.Copy()
				d.appendMessage("Game loaded from " + arg)
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

	if strings.Contains(strings.ToLower(result), "win") {
		d.Done = true
	}
}

func (d *DialogScreen) appendMessage(msg string) {
	lines := strings.Split(msg, "\n")
	d.Messages = append(d.Messages, lines...)
}

func (d *DialogScreen) Draw(screen *ebiten.Image) {
	if d.Background != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(d.Background, op)
	}

	height := screen.Bounds().Dy()

	// Prepare input and turn display
	inputLine := "> " + d.Input
	player := d.Game.CurrentPlayer()
	turnInfo := fmt.Sprintf("%s's turn", player.ID)

	y := height - yMargin - 3*lineHeight

	// Draw messages bottom-up, but only if in-bounds
	for i := len(d.Messages) - 1; i >= 0; i-- {
		if y < yMargin {
			break // don’t draw past top margin
		}
		text.Draw(screen, d.Messages[i], MainFont, xMargin, y, color.White)
		y -= lineHeight
	}

	// Draw turn and input lines
	text.Draw(screen, turnInfo, MainFont, xMargin, height-yMargin-HeadlineHeight, color.RGBA{200, 200, 0, 255})
	text.Draw(screen, inputLine, MainFont, xMargin, height-yMargin, color.White)

	drawButtonWithImage(screen, dialogExitButtonX, dialogExitButtonY, sideButtonWidth, sideButtonHeight, "", d.ExitButton)
}
