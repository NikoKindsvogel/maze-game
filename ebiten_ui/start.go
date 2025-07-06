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
	startButtonWidth  = 484
	startButtonHeight = 47

	startButtonY = 500
	exitButtonY  = 570
)

func NewStartScreen() *StartScreen {
	bgImage := loadImageFromEmbed("backgrounds/startscreen_background.png")
	startImage := loadImageFromEmbed("buttons/startscreen_button_new.png")
	exitImage := loadImageFromEmbed("buttons/startscreen_button_exit.png")
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
			u.config = NewConfigScreen()
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
