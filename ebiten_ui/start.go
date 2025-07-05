package ebiten_ui

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type StartScreen struct {
	mouseWasDown bool

	Background  *ebiten.Image
	StartButton *ebiten.Image
	ExitButton  *ebiten.Image
}

const (
	screenWidth  = 1200
	screenHeight = 900

	buttonWidth  = 484
	buttonHeight = 47

	startButtonY = 500
	exitButtonY  = 570
)

func NewStartScreen() *StartScreen {
	bgImage := loadImage("assets/backgrounds/startscreen_background.png")
	startImage := loadImage("assets/buttons/startscreen_button_new.png")
	exitImage := loadImage("assets/buttons/startscreen_button_exit.png")
	return &StartScreen{
		mouseWasDown: false,
		Background:   bgImage,
		StartButton:  startImage,
		ExitButton:   exitImage,
	}
}

func (s *StartScreen) Update(u *UIManager) error {
	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	x, y := ebiten.CursorPosition()

	buttonX := (screenWidth - buttonWidth) / 2

	if mouseDown && !s.mouseWasDown {
		if x >= buttonX && x <= buttonX+buttonWidth &&
			y >= startButtonY && y <= startButtonY+buttonHeight {
			u.screen = ScreenConfig
		}
		if x >= buttonX && x <= buttonX+buttonWidth &&
			y >= exitButtonY && y <= exitButtonY+buttonHeight {
			os.Exit(0)
		}
	}

	s.mouseWasDown = mouseDown
	return nil
}

func (s *StartScreen) Draw(screen *ebiten.Image) {
	if s.Background != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(s.Background, op)
	}

	buttonX := (screenWidth - buttonWidth) / 2

	drawButtonWithImage(screen, buttonX, startButtonY, buttonWidth, buttonHeight, "", s.StartButton)
	drawButtonWithImage(screen, buttonX, exitButtonY, buttonWidth, buttonHeight, "", s.ExitButton)
}
