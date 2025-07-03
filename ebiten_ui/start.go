package ebiten_ui

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type StartScreen struct {
	mouseWasDown bool
}

const (
	startButtonX = 200
	startButtonY = 200
	buttonWidth  = 160
	buttonHeight = 40

	exitButtonX = 200
	exitButtonY = 260
)

func NewStartScreen() *StartScreen {
	return &StartScreen{
		mouseWasDown: false,
	}
}

func (s *StartScreen) Update(u *UIManager) error {
	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	x, y := ebiten.CursorPosition()

	if mouseDown && !s.mouseWasDown {
		if x >= startButtonX && x <= startButtonX+buttonWidth &&
			y >= startButtonY && y <= startButtonY+buttonHeight {
			u.screen = ScreenConfig
		}
		if x >= exitButtonX && x <= exitButtonX+buttonWidth &&
			y >= exitButtonY && y <= exitButtonY+buttonHeight {
			os.Exit(0)
		}
	}

	s.mouseWasDown = mouseDown
	return nil
}

func (s *StartScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})

	drawButton(screen, startButtonX, startButtonY, "Start Game")
	drawButton(screen, exitButtonX, exitButtonY, "Exit Game")
}
