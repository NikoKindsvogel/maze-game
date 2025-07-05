package ebiten_ui

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type StartScreen struct {
	Background  *ebiten.Image
	StartButton *ebiten.Image
	ExitButton  *ebiten.Image
}

const (
	screenWidth  = 1200
	screenHeight = 900

	startButtonWidth  = 484
	startButtonHeight = 47

	startButtonY = 500
	exitButtonY  = 570
)

func NewStartScreen() *StartScreen {
	bgImage := loadImage("assets/backgrounds/startscreen_background.png")
	startImage := loadImage("assets/buttons/startscreen_button_new.png")
	exitImage := loadImage("assets/buttons/startscreen_button_exit.png")
	return &StartScreen{
		Background:  bgImage,
		StartButton: startImage,
		ExitButton:  exitImage,
	}
}

func (s *StartScreen) Update(u *UIManager) error {
	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	x, y := ebiten.CursorPosition()

	buttonX := (screenWidth - startButtonWidth) / 2

	if mouseDown && !u.mouseWasDown {
		if x >= buttonX && x <= buttonX+startButtonWidth &&
			y >= startButtonY && y <= startButtonY+startButtonHeight {
			u.screen = ScreenConfig
		}
		if x >= buttonX && x <= buttonX+startButtonWidth &&
			y >= exitButtonY && y <= exitButtonY+startButtonHeight {
			os.Exit(0)
		}
	}

	u.mouseWasDown = mouseDown
	return nil
}

func (s *StartScreen) Draw(screen *ebiten.Image) {
	if s.Background != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(s.Background, op)
	}

	buttonX := (screenWidth - startButtonWidth) / 2

	drawButtonWithImage(screen, buttonX, startButtonY, startButtonWidth, startButtonHeight, "", s.StartButton)
	drawButtonWithImage(screen, buttonX, exitButtonY, startButtonWidth, startButtonHeight, "", s.ExitButton)
}
