package ebiten_ui

import (
	"image/color"
	"math"
	"maze-game/game"
	"maze-game/maze"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	cellSize   = 64 // width and height of each cell image
	wallOffset = 8  // half-width for wall overlap
	spacing    = 20 // space between the two mazes
	buttonX    = 20
	buttonY    = 20
	buttonW    = 140
	buttonH    = 40
)

const (
	revealExitButtonX = screenWidth - 130
	revealExitButtonY = screenHeight - 150
	revealShowButtonX = screenWidth - 130
	revealShowButtonY = screenHeight - 250
)

type RevealScreen struct {
	StartGame        *game.Game
	FinalGame        *game.Game
	Images           map[maze.CellType]*ebiten.Image
	WallH            *ebiten.Image
	WallV            *ebiten.Image
	Treasure         *ebiten.Image
	Treasure_big     *ebiten.Image
	RiverCorner      *ebiten.Image
	ShowCurrent      bool
	PlayerImages     []*ebiten.Image
	PlayerImagePaths []string
	Background       *ebiten.Image
	ExitButton       *ebiten.Image
	ShowNowButton    *ebiten.Image
	ShowStartButton  *ebiten.Image
	PlayerBackground *ebiten.Image
}

func NewRevealScreen(start, final *game.Game) *RevealScreen {
	images := loadCellImages()
	playerImages, playerImagePaths := loadPlayerImages()
	wallH := loadImage("assets/walls/wall_horizontal_long.png")
	wallV := loadImage("assets/walls/wall_vertical_long.png")
	treasure := loadImage("assets/cells/treasure.png")
	treasure_big := loadImage("assets/cells/cell_treasure.png")
	river_corner := loadImage("assets/cells/cell_river_corner.png")
	bgImage := loadImage("assets/backgrounds/background.png")
	exitButton := loadImage("assets/buttons/reveal_button_exitgame.png")
	showNowButton := loadImage("assets/buttons/reveal_button_shownow.png")
	showStartButton := loadImage("assets/buttons/reveal_button_showstart.png")
	playerBackground := loadImage("assets/buttons/reveal_button_playercolor.png")

	return &RevealScreen{
		StartGame:        start,
		FinalGame:        final,
		Images:           images,
		WallH:            wallH,
		WallV:            wallV,
		Treasure:         treasure,
		Treasure_big:     treasure_big,
		RiverCorner:      river_corner,
		PlayerImages:     playerImages,
		PlayerImagePaths: playerImagePaths,
		Background:       bgImage,
		ExitButton:       exitButton,
		ShowNowButton:    showNowButton,
		ShowStartButton:  showStartButton,
		PlayerBackground: playerBackground,
	}
}

func (r *RevealScreen) Update(u *UIManager) error {
	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	x, y := ebiten.CursorPosition()

	if mouseDown && !u.mouseWasDown {
		if x >= revealExitButtonX && x <= revealExitButtonX+sideButtonWidth &&
			y >= revealExitButtonY && y <= revealExitButtonY+sideButtonHeight {
			u.screen = ScreenStart
		}
		if x >= revealShowButtonX && x <= revealShowButtonX+sideButtonWidth &&
			y >= revealShowButtonY && y <= revealShowButtonY+sideButtonHeight {
			r.ShowCurrent = !r.ShowCurrent
		}
	}

	u.mouseWasDown = mouseDown
	return nil
}

func (r *RevealScreen) Draw(screen *ebiten.Image) {

	if r.Background != nil {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(r.Background, op)
	}

	// Draw maze based on toggle
	if r.ShowCurrent {
		r.drawGame(screen, r.FinalGame)
	} else {
		r.drawGame(screen, r.StartGame)
	}

	if !r.ShowCurrent {
		drawButtonWithImage(screen, revealShowButtonX, revealShowButtonY, sideButtonWidth, sideButtonHeight, "", r.ShowNowButton)
	} else {
		drawButtonWithImage(screen, revealShowButtonX, revealShowButtonY, sideButtonWidth, sideButtonHeight, "", r.ShowStartButton)
	}

	drawButtonWithImage(screen, revealExitButtonX, revealExitButtonY, sideButtonWidth, sideButtonHeight, "", r.ExitButton)
}

func (r *RevealScreen) drawGame(screen *ebiten.Image, g *game.Game) {
	m := g.GetMaze()
	players := g.Players

	ox := (screenWidth - m.Size*cellSize) / 2
	oy := (screenHeight - m.Size*cellSize) / 2

	for row := 0; row < m.Size; row++ {
		for col := 0; col < m.Size; col++ {
			r.drawCell(screen, m, row, col, ox, oy)
		}
	}
	r.drawInnerWalls(screen, m, ox, oy)
	r.drawBorderWalls(screen, m, ox, oy)
	r.drawPlayers(screen, ox, oy, players)
}

func (r *RevealScreen) drawCell(screen *ebiten.Image, m *maze.Maze, row, col, ox, oy int) {
	cell := m.Grid[row][col]
	x := ox + col*cellSize
	y := oy + row*cellSize

	switch cell.Type {
	case maze.Exit:
		r.drawExit(screen, *cell, row, col, x, y, m.Size)
	case maze.River, maze.Estuary:
		r.drawRiverOrEstuary(screen, *cell, row, col, x, y, m)
	default:
		img := r.Images[cell.Type]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(img, op)
	}

	// Treasure overlay
	if m.TreasureOnMap && m.TreasureRow == row && m.TreasureCol == col {
		tOp := &ebiten.DrawImageOptions{}
		tOp.GeoM.Translate(float64(x), float64(y))
		if cell.Type == maze.Armory || cell.Type == maze.Hole || cell.Type == maze.Hospital {
			screen.DrawImage(r.Treasure, tOp)
		} else {
			screen.DrawImage(r.Treasure_big, tOp)
		}
	}
}

func (r *RevealScreen) drawExit(screen *ebiten.Image, cell maze.Cell, row, col, x, y, size int) {
	img := r.Images[cell.Type]
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(cellSize)/2, -float64(cellSize)/2)

	switch {
	case row == 0:
		op.GeoM.Rotate(-math.Pi / 2)
	case row == size-1:
		op.GeoM.Rotate(math.Pi / 2)
	case col == 0:
		op.GeoM.Rotate(math.Pi)
	}

	op.GeoM.Translate(float64(x)+float64(cellSize)/2, float64(y)+float64(cellSize)/2)
	screen.DrawImage(img, op)
}

func (r *RevealScreen) drawRiverOrEstuary(screen *ebiten.Image, cell maze.Cell, row, col, x, y int, m *maze.Maze) {
	dir := cell.RiverDir
	var img *ebiten.Image
	op := &ebiten.DrawImageOptions{}

	if cell.Type == maze.River {
		var prevDir maze.Direction
		foundPrev := false
		for _, d := range []maze.Direction{maze.Up, maze.Right, maze.Down, maze.Left} {
			dr, dc := maze.Delta(d)
			pr, pc := row+dr, col+dc
			if m.InBounds(pr, pc) && m.Grid[pr][pc].Type == maze.River {
				if m.Grid[pr][pc].RiverDir == maze.Opposite(d) {
					prevDir = d
					foundPrev = true
					break
				}
			}
		}
		nextDir := dir
		if foundPrev && prevDir != maze.Opposite(nextDir) {
			img = r.RiverCorner
			op.GeoM.Translate(-float64(cellSize)/2, -float64(cellSize)/2)
			switch {
			case prevDir == maze.Down && nextDir == maze.Right,
				prevDir == maze.Right && nextDir == maze.Down:
				op.GeoM.Rotate(0)
			case prevDir == maze.Right && nextDir == maze.Up,
				prevDir == maze.Up && nextDir == maze.Right:
				op.GeoM.Rotate(3 * math.Pi / 2)
			case prevDir == maze.Up && nextDir == maze.Left,
				prevDir == maze.Left && nextDir == maze.Up:
				op.GeoM.Rotate(math.Pi)
			case prevDir == maze.Left && nextDir == maze.Down,
				prevDir == maze.Down && nextDir == maze.Left:
				op.GeoM.Rotate(math.Pi / 2)
			}
			op.GeoM.Translate(float64(x)+float64(cellSize)/2, float64(y)+float64(cellSize)/2)
		} else {
			img = r.Images[cell.Type]
			if dir == maze.Up || dir == maze.Down {
				op.GeoM.Translate(-float64(cellSize)/2, -float64(cellSize)/2)
				op.GeoM.Rotate(math.Pi / 2)
				op.GeoM.Translate(float64(x)+float64(cellSize)/2, float64(y)+float64(cellSize)/2)
			} else {
				op.GeoM.Translate(float64(x), float64(y))
			}
		}
	} else if cell.Type == maze.Estuary {
		img = r.Images[cell.Type]
		op.GeoM.Translate(-float64(cellSize)/2, -float64(cellSize)/2)
		switch dir {
		case maze.Up:
			op.GeoM.Rotate(-math.Pi / 2)
		case maze.Down:
			op.GeoM.Rotate(math.Pi / 2)
		case maze.Left:
			op.GeoM.Rotate(math.Pi)
		case maze.Right:
		}
		op.GeoM.Translate(float64(x)+float64(cellSize)/2, float64(y)+float64(cellSize)/2)
	}
	screen.DrawImage(img, op)
}

func (r *RevealScreen) drawInnerWalls(screen *ebiten.Image, m *maze.Maze, ox, oy int) {
	halfWall := float64(wallOffset) / 2
	for row := 0; row < m.Size; row++ {
		for col := 0; col < m.Size; col++ {
			x := ox + col*cellSize
			y := oy + row*cellSize
			cell := m.Grid[row][col]

			if cell.Walls[maze.Right] && col < m.Size-1 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x+cellSize)-halfWall, float64(y)-halfWall)
				screen.DrawImage(r.WallV, op)
			}
			if cell.Walls[maze.Down] && row < m.Size-1 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)-halfWall, float64(y+cellSize)-halfWall)
				screen.DrawImage(r.WallH, op)
			}
		}
	}
}

func (r *RevealScreen) drawBorderWalls(screen *ebiten.Image, m *maze.Maze, ox, oy int) {
	halfWall := float64(wallOffset) / 2

	for row := 0; row < m.Size; row++ {
		y := oy + row*cellSize
		opL := &ebiten.DrawImageOptions{}
		opL.GeoM.Translate(float64(ox)-halfWall, float64(y)-halfWall)
		screen.DrawImage(r.WallV, opL)

		opR := &ebiten.DrawImageOptions{}
		opR.GeoM.Translate(float64(ox+m.Size*cellSize)-halfWall, float64(y)-halfWall)
		screen.DrawImage(r.WallV, opR)
	}

	for col := 0; col < m.Size; col++ {
		x := ox + col*cellSize
		opT := &ebiten.DrawImageOptions{}
		opT.GeoM.Translate(float64(x)-halfWall, float64(oy)-halfWall)
		screen.DrawImage(r.WallH, opT)

		opB := &ebiten.DrawImageOptions{}
		opB.GeoM.Translate(float64(x)-halfWall, float64(oy+m.Size*cellSize)-halfWall)
		screen.DrawImage(r.WallH, opB)
	}
}

func (r *RevealScreen) drawPlayers(screen *ebiten.Image, ox, oy int, players []*game.Player) {
	legendX := 20
	legendY := 20
	lineHeight := 50 // spacing between entries (including background)

	for i, player := range players {
		if i >= len(r.PlayerImages) {
			continue // ignore extra players
		}

		// --- Draw player on maze ---
		img := r.PlayerImages[i]
		x := ox + player.Col*cellSize
		y := oy + player.Row*cellSize
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(img, op)

		// --- Extract color from image filename ---
		color := "unknown"
		if i < len(r.PlayerImagePaths) {
			color = extractColorFromFilename(r.PlayerImagePaths[i])
		}

		// --- Draw background for legend item ---
		bgOp := &ebiten.DrawImageOptions{}
		bgOp.GeoM.Translate(float64(legendX), float64(legendY+i*lineHeight))
		screen.DrawImage(r.PlayerBackground, bgOp)

		// --- Draw player name + color on top of background ---
		textX := legendX + 10
		textY := legendY + i*lineHeight + 30 // vertical offset inside background
		text.Draw(screen, player.ID, MainFont, textX, textY, colorRGBA(color))
	}
}

func loadCellImages() map[maze.CellType]*ebiten.Image {
	return map[maze.CellType]*ebiten.Image{
		maze.Empty:    loadImage("assets/cells/cell_empty.png"),
		maze.Hospital: loadImage("assets/cells/cell_hospital.png"),
		maze.Exit:     loadImage("assets/cells/cell_exit.png"),
		maze.Hole:     loadImage("assets/cells/cell_hole.png"),
		maze.Dragon:   loadImage("assets/cells/cell_dragon.png"),
		maze.Armory:   loadImage("assets/cells/cell_armory.png"),
		maze.River:    loadImage("assets/cells/cell_river.png"),
		maze.Estuary:  loadImage("assets/cells/cell_estuary.png"),
	}
}

func colorRGBA(name string) color.Color {
	switch strings.ToLower(name) {
	case "red":
		return color.RGBA{R: 220, A: 255}
	case "blue":
		return color.RGBA{B: 220, A: 255}
	case "green":
		return color.RGBA{G: 200, A: 255}
	case "yellow":
		return color.RGBA{R: 240, G: 240, A: 255}
	case "white":
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	case "magenta":
		return color.RGBA{R: 235, G: 18, B: 234, A: 255}
	case "cyan":
		return color.RGBA{R: 19, G: 235, B: 233, A: 255}
	case "black":
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	default:
		return color.White
	}
}

func loadPlayerImages() ([]*ebiten.Image, []string) {
	paths := []string{
		"assets/players2/player_cyan.png",
		"assets/players2/player_magenta.png",
		"assets/players2/player_white.png",
		"assets/players2/player_yellow.png",
	}
	images := make([]*ebiten.Image, len(paths))
	for i, path := range paths {
		images[i] = loadImage(path)
	}
	return images, paths
}

func extractColorFromFilename(filename string) string {
	filename = strings.TrimSuffix(filename, ".png")
	parts := strings.Split(filename, "_")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown"
}
