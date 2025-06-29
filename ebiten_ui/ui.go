package ebiten_ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type UIScreen int

const (
	ScreenConfig UIScreen = iota
	ScreenDialog
	ScreenGame
	ScreenReveal
)

type UIManager struct {
	screen UIScreen
	config *ConfigScreen
	dialog *DialogScreen
	game   *GameScreen
	reveal *RevealScreen
}

func NewUIManager() *UIManager {
	return &UIManager{
		screen: ScreenConfig,
		config: NewConfigScreen(),
	}
}

func (u *UIManager) Update() error {
	switch u.screen {
	case ScreenConfig:
		if u.config.Done {
			u.dialog = NewDialogScreen()
			u.screen = ScreenDialog
		} else {
			u.config.Update()
		}
	case ScreenDialog:
		if u.dialog.Done {
			u.game = NewGameScreen()
			u.screen = ScreenGame
		} else {
			u.dialog.Update()
		}
	case ScreenGame:
		if u.game.Finished {
			u.reveal = NewRevealScreen(u.game.MazeStart, u.game.MazeEnd)
			u.screen = ScreenReveal
		} else {
			u.game.Update()
		}
	case ScreenReveal:
		u.reveal.Update()
	}
	return nil
}

func (u *UIManager) Draw(screen *ebiten.Image) {
	switch u.screen {
	case ScreenConfig:
		u.config.Draw(screen)
	case ScreenDialog:
		u.dialog.Draw(screen)
	case ScreenGame:
		u.game.Draw(screen)
	case ScreenReveal:
		u.reveal.Draw(screen)
	}
}

func (u *UIManager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
