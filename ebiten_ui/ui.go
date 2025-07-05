package ebiten_ui

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type UIScreen int

const (
	ScreenStart UIScreen = iota
	ScreenConfig
	ScreenDialog
	ScreenGame
	ScreenReveal
)

type UIManager struct {
	screen  UIScreen
	config  *ConfigScreen
	dialog  *DialogScreen
	reveal  *RevealScreen
	start   *StartScreen
	endGame bool
}

func NewUIManager() *UIManager {
	return &UIManager{
		screen: ScreenStart,
		start:  NewStartScreen(),
		config: NewConfigScreen(),
	}
}

func (u *UIManager) Update() error {
	if u.endGame {
		os.Exit(0)
	}
	switch u.screen {
	case ScreenStart:
		u.start.Update(u)
	case ScreenConfig:
		if u.config.Done {
			size, holes, riverPush, names := u.config.GetConfig()
			u.dialog = NewDialogScreen(size, holes, riverPush, names)
			u.screen = ScreenDialog
		} else {
			u.config.Update()
		}
	case ScreenDialog:
		if u.dialog.Done {
			u.reveal = NewRevealScreen(&u.dialog.startGame, u.dialog.Game)
			u.screen = ScreenReveal
		} else {
			u.dialog.Update()
		}
	case ScreenReveal:
		u.reveal.Update()
	}
	return nil
}

func (u *UIManager) Draw(screen *ebiten.Image) {
	switch u.screen {
	case ScreenStart:
		u.start.Draw(screen)
	case ScreenConfig:
		u.config.Draw(screen)
	case ScreenDialog:
		u.dialog.Draw(screen)
	case ScreenReveal:
		u.reveal.Draw(screen)
	}
}

func (u *UIManager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func drawButton(screen *ebiten.Image, x, y, buttonWidth, buttonHeight int, label string) {
	// Draw button rectangle
	btn := ebiten.NewImage(buttonWidth, buttonHeight)
	btn.Fill(color.RGBA{100, 100, 255, 255}) // blue

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(btn, op)

	// Draw label text
	textX := x + 10
	textY := y + 25
	text.Draw(screen, label, basicfont.Face7x13, textX, textY, color.White)
}

// Draws a button using a provided image and overlays a centered label.
func drawButtonWithImage(screen *ebiten.Image, x, y, buttonWidth, buttonHeight int, label string, image *ebiten.Image) {
	// Draw the image as button background
	if image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(image, op)
	} else {
		// fallback: draw solid color if image is nil
		btn := ebiten.NewImage(buttonWidth, buttonHeight)
		btn.Fill(color.RGBA{100, 100, 255, 255}) // blue fallback
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(btn, op)
	}

	// Center the label text inside the button
	textBounds := text.BoundString(basicfont.Face7x13, label)
	textWidth := textBounds.Dx()
	textHeight := textBounds.Dy()

	textX := x + (buttonWidth-textWidth)/2
	textY := y + (buttonHeight+textHeight)/2 - 2 // adjust for baseline

	text.Draw(screen, label, basicfont.Face7x13, textX, textY, color.White)
}

func loadImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
