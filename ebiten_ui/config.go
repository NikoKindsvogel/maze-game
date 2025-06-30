package ebiten_ui

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ConfigScreen struct {
	Done           bool
	inputs         []string
	currentField   int
	fieldLabels    []string
	playerNames    []string
	inputtingNames bool
	numPlayers     int

	enterPressedLastFrame bool
}

func NewConfigScreen() *ConfigScreen {
	return &ConfigScreen{
		Done: false,
		inputs: []string{
			"7", // Maze size
			"3", // Num holes
			"2", // Num players
			"2", // River push
		},
		fieldLabels: []string{
			"Maze Size:",
			"Number of Holes:",
			"Number of Players:",
			"River Push Distance:",
		},
		currentField: 0,
	}
}

func (c *ConfigScreen) Update() {
	enterPressed := ebiten.IsKeyPressed(ebiten.KeyEnter)

	if enterPressed && !c.enterPressedLastFrame {
		// Only process Enter on the frame it was first pressed
		if c.inputtingNames {
			if c.currentField == len(c.playerNames)-1 {
				c.Done = true
			} else {
				c.currentField++
			}
		} else {
			if c.currentField == len(c.inputs)-1 {
				// Start name input after all numeric inputs
				n, err := strconv.Atoi(c.inputs[2])
				if err != nil || n <= 0 || n > 8 {
					n = 2
				}
				c.numPlayers = n
				c.playerNames = make([]string, n)
				c.currentField = 0
				c.inputtingNames = true
			} else {
				c.currentField++
			}
		}
	}
	c.enterPressedLastFrame = enterPressed

	// Handle text input
	for _, r := range ebiten.AppendInputChars(nil) {
		if unicode.IsPrint(r) {
			if c.inputtingNames {
				c.playerNames[c.currentField] += string(r)
			} else {
				c.inputs[c.currentField] += string(r)
			}
		}
	}

	// Handle backspace
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		if c.inputtingNames {
			s := c.playerNames[c.currentField]
			if len(s) > 0 {
				c.playerNames[c.currentField] = s[:len(s)-1]
			}
		} else {
			s := c.inputs[c.currentField]
			if len(s) > 0 {
				c.inputs[c.currentField] = s[:len(s)-1]
			}
		}
	}
}

func (c *ConfigScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})

	title := "Maze Game Config - Press ENTER to confirm each field"
	ebitenutil.DebugPrintAt(screen, title, 20, 20)

	if !c.inputtingNames {
		for i, label := range c.fieldLabels {
			text := fmt.Sprintf("%s %s", label, c.inputs[i])
			if i == c.currentField {
				text += " <"
			}
			ebitenutil.DebugPrintAt(screen, text, 40, 80+i*30)
		}
	} else {
		ebitenutil.DebugPrintAt(screen, "Enter player names:", 40, 80)
		for i := 0; i < c.numPlayers; i++ {
			text := fmt.Sprintf("Player %d: %s", i+1, c.playerNames[i])
			if i == c.currentField {
				text += " <"
			}
			ebitenutil.DebugPrintAt(screen, text, 60, 120+i*30)
		}
	}
}

func (c *ConfigScreen) GetConfig() (size, holes, riverPush int, names []string) {
	size, _ = strconv.Atoi(c.inputs[0])
	holes, _ = strconv.Atoi(c.inputs[1])
	riverPush, _ = strconv.Atoi(c.inputs[3])

	names = make([]string, len(c.playerNames))
	for i, name := range c.playerNames {
		name = strings.TrimSpace(name)
		if name == "" {
			name = fmt.Sprintf("P%d", i+1)
		}
		names[i] = name
	}

	return
}
