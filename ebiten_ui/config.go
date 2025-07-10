package ebiten_ui

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
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

	Background *ebiten.Image
}

func NewConfigScreen() *ConfigScreen {
	bgImage := loadImageFromEmbed("backgrounds/background.png")
	return &ConfigScreen{
		Done: false,
		inputs: []string{
			"6", // Maze size
			"3", // Num holes
			"2", // Num players
			"2", // River push
			"8", // River lengt
		},
		fieldLabels: []string{
			"Maze Size:",
			"Number of Holes:",
			"Number of Players:",
			"River Push Distance:",
			"River Length:",
		},
		currentField: 0,
		Background:   bgImage,
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
	var _ = ebitenutil.DebugPrintAt
	if c.Background != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(c.Background, op)
	}

	title := "Maze Game Config - Press ENTER to confirm each field"
	text.Draw(screen, title, HeadlineFont, xMargin, yMargin, color.White)

	if !c.inputtingNames {
		for i, label := range c.fieldLabels {
			line := fmt.Sprintf("%s %s", label, c.inputs[i])
			if i == c.currentField {
				line += " <"
			}
			text.Draw(screen, line, MainFont, xMargin+HeadlineHeight, yMargin+HeadlineHeight+i*lineHeight, color.White)
		}
	} else {
		text.Draw(screen, "Enter player names:", MainFont, xMargin+HeadlineHeight, yMargin+HeadlineHeight, color.White)
		for i := 0; i < c.numPlayers; i++ {
			line := fmt.Sprintf("Player %d: %s", i+1, c.playerNames[i])
			if i == c.currentField {
				line += " <"
			}
			text.Draw(screen, line, MainFont, xMargin+2*HeadlineHeight, yMargin+HeadlineHeight+lineHeight+(i*lineHeight), color.White)
		}
	}
}

func (c *ConfigScreen) GetConfig() (size, holes, riverLength, riverPush int, names []string) {
	size, _ = strconv.Atoi(c.inputs[0])
	holes, _ = strconv.Atoi(c.inputs[1])
	riverPush, _ = strconv.Atoi(c.inputs[3])
	riverLength, _ = strconv.Atoi(c.inputs[4])

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
