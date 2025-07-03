package ebiten_ui

import (
	"image"
	"log"
	"math"
	"maze-game/maze"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	cellSize   = 64 // width and height of each cell image
	wallOffset = 8  // half-width for wall overlap
	spacing    = 20 // space between the two mazes
	buttonX    = 20
	buttonY    = 20
	buttonW    = 140
	buttonH    = 40
	offsetX    = 100
	offsetY    = 100
)

type RevealScreen struct {
	StartMaze    *maze.Maze
	FinalMaze    *maze.Maze
	Images       map[maze.CellType]*ebiten.Image
	WallH        *ebiten.Image
	WallV        *ebiten.Image
	Treasure     *ebiten.Image
	Treasure_big *ebiten.Image
	RiverCorner  *ebiten.Image
	ShowCurrent  bool
	mouseWasDown bool
}

func NewRevealScreen(start, final *maze.Maze) *RevealScreen {
	images := loadCellImages()
	wallH := loadImage("assets/wall_horizontal_long.png")
	wallV := loadImage("assets/wall_vertical_long.png")
	treasure := loadImage("assets/treasure.png")
	treasure_big := loadImage("assets/cell_treasure.png")
	river_corner := loadImage("assets/cell_river_corner.png")

	// for testing:
	// wallH := createWallImage(cellSize, 8, color.RGBA{0, 0, 0, 255}) // black horizontal wall
	// wallV := createWallImage(8, cellSize, color.RGBA{0, 0, 0, 255}) // black vertical wall
	// treasure := createColoredImage(color.RGBA{255, 215, 0, 128})    // semi-transparent gold

	return &RevealScreen{
		StartMaze:    start,
		FinalMaze:    final,
		Images:       images,
		WallH:        wallH,
		WallV:        wallV,
		Treasure:     treasure,
		Treasure_big: treasure_big,
		RiverCorner:  river_corner,
	}
}

func (r *RevealScreen) Update() error {
	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	x, y := ebiten.CursorPosition()

	if mouseDown && !r.mouseWasDown {
		// Mouse just pressed — check if it hit the button
		if x >= buttonX && x <= buttonX+buttonW && y >= buttonY && y <= buttonY+buttonH {
			r.ShowCurrent = !r.ShowCurrent
		}
	}

	r.mouseWasDown = mouseDown
	return nil
}

func (r *RevealScreen) Draw(screen *ebiten.Image) {

	if !r.ShowCurrent {
		drawButton(screen, buttonX, buttonY, "Show Current")
	} else {
		drawButton(screen, buttonX, buttonY, "Show Start")
	}

	// Draw maze based on toggle
	if r.ShowCurrent {
		r.drawMaze(screen, r.FinalMaze, offsetX, offsetY)
	} else {
		r.drawMaze(screen, r.StartMaze, offsetX, offsetY)
	}
}

func (r *RevealScreen) drawMaze(screen *ebiten.Image, m *maze.Maze, ox, oy int) {
	for row := 0; row < m.Size; row++ {
		for col := 0; col < m.Size; col++ {
			cell := m.Grid[row][col]
			x := ox + col*cellSize
			y := oy + row*cellSize

			// When drawing the exit tile inside drawMaze, after computing x, y:

			if cell.Type == maze.Exit {
				img := r.Images[cell.Type]
				op := &ebiten.DrawImageOptions{}

				// Move origin to center for rotation
				op.GeoM.Translate(-float64(cellSize)/2, -float64(cellSize)/2)

				// Determine rotation based on exit position
				switch {
				case row == 0: // Top edge → look up, rotate -90°
					op.GeoM.Rotate(-math.Pi / 2)
				case row == m.Size-1: // Bottom edge → look down, rotate 90°
					op.GeoM.Rotate(math.Pi / 2)
				case col == 0: // Left edge → look left, rotate 180°
					op.GeoM.Rotate(math.Pi)
				case col == m.Size-1: // Right edge → default (looking right), no rotation
					// no rotation needed
				}

				// Move origin back to top-left plus cell position
				op.GeoM.Translate(float64(x)+float64(cellSize)/2, float64(y)+float64(cellSize)/2)

				screen.DrawImage(img, op)
				continue // skip normal drawing for this cell since already drawn
			}

			if cell.Type == maze.River || cell.Type == maze.Estuary {
				dir := cell.RiverDir

				var img *ebiten.Image
				op := &ebiten.DrawImageOptions{}

				if cell.Type == maze.River {
					// Try to determine previous direction (incoming river)
					var prevDir maze.Direction
					foundPrev := false

					for _, d := range []maze.Direction{maze.Up, maze.Right, maze.Down, maze.Left} {
						dr, dc := maze.Delta(d) // unpack delta properly
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
					foundNext := true

					if foundPrev && foundNext && prevDir != maze.Opposite(nextDir) {
						// Corner tile
						img = r.RiverCorner

						// Rotate around center
						op.GeoM.Translate(-float64(cellSize)/2, -float64(cellSize)/2)

						switch {
						case prevDir == maze.Down && nextDir == maze.Right:
							op.GeoM.Rotate(0) // 0°
						case prevDir == maze.Right && nextDir == maze.Up:
							op.GeoM.Rotate(3 * math.Pi / 2) // 270°
						case prevDir == maze.Up && nextDir == maze.Left:
							op.GeoM.Rotate(math.Pi) // 180°
						case prevDir == maze.Left && nextDir == maze.Down:
							op.GeoM.Rotate(math.Pi / 2) // 90°
						case prevDir == maze.Right && nextDir == maze.Down:
							op.GeoM.Rotate(0) // 0°
						case prevDir == maze.Up && nextDir == maze.Right:
							op.GeoM.Rotate(3 * math.Pi / 2) // 270°
						case prevDir == maze.Left && nextDir == maze.Up:
							op.GeoM.Rotate(math.Pi) // 180°
						case prevDir == maze.Down && nextDir == maze.Left:
							op.GeoM.Rotate(math.Pi / 2) // 90°
						default:
							img = r.Images[cell.Type]
						}

						op.GeoM.Translate(float64(x)+float64(cellSize)/2, float64(y)+float64(cellSize)/2)
					} else {
						// Straight tile
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
					img = r.Images[cell.Type] // or estuary image if you have one

					op.GeoM.Translate(-float64(cellSize)/2, -float64(cellSize)/2) // move origin to center

					switch dir {
					case maze.Up:
						op.GeoM.Rotate(-math.Pi / 2) // rotate -90° (river coming from Up → face Down)
					case maze.Down:
						op.GeoM.Rotate(math.Pi / 2) // rotate 90° (river coming from Down → face Up)
					case maze.Left:
						op.GeoM.Rotate(math.Pi) // rotate 180° (river coming from Left → face Right)
					case maze.Right:
						op.GeoM.Rotate(0) // no rotation needed (river from Right → face Left)
					}

					op.GeoM.Translate(float64(x)+float64(cellSize)/2, float64(y)+float64(cellSize)/2) // move origin back to top-left + position
				}

				screen.DrawImage(img, op)
			} else {
				// Normal cell
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
	}

	halfWall := float64(wallOffset) / 2

	// 1. Draw inner walls first (between cells)
	for row := 0; row < m.Size; row++ {
		for col := 0; col < m.Size; col++ {
			x := ox + col*cellSize
			y := oy + row*cellSize
			cell := m.Grid[row][col]

			// vertical inner wall (right)
			if cell.Walls[maze.Right] && col < m.Size-1 {
				wallOp := &ebiten.DrawImageOptions{}
				wallOp.GeoM.Translate(float64(x+cellSize)-halfWall, float64(y)-halfWall)
				screen.DrawImage(r.WallV, wallOp)
			}
			// horizontal inner wall (down)
			if cell.Walls[maze.Down] && row < m.Size-1 {
				wallOp := &ebiten.DrawImageOptions{}
				wallOp.GeoM.Translate(float64(x)-halfWall, float64(y+cellSize)-halfWall)
				screen.DrawImage(r.WallH, wallOp)
			}
		}
	}

	// Left and right vertical borders
	for row := 0; row < m.Size; row++ {
		y := oy + row*cellSize

		// Left border wall, shifted half outside to the left (negative)
		wallOpLeft := &ebiten.DrawImageOptions{}
		wallOpLeft.GeoM.Translate(float64(ox)-halfWall, float64(y)-halfWall)
		screen.DrawImage(r.WallV, wallOpLeft)

		// Right border wall, shifted half inside the last cell
		wallOpRight := &ebiten.DrawImageOptions{}
		wallOpRight.GeoM.Translate(float64(ox+m.Size*cellSize)-halfWall, float64(y)-halfWall)
		screen.DrawImage(r.WallV, wallOpRight)
	}

	// Top and bottom horizontal borders
	for col := 0; col < m.Size; col++ {
		x := ox + col*cellSize

		// Top border wall, shifted half outside above the maze
		wallOpTop := &ebiten.DrawImageOptions{}
		wallOpTop.GeoM.Translate(float64(x)-halfWall, float64(oy)-halfWall)
		screen.DrawImage(r.WallH, wallOpTop)

		// Bottom border wall, shifted half inside the last row cell
		wallOpBottom := &ebiten.DrawImageOptions{}
		wallOpBottom.GeoM.Translate(float64(x)-halfWall, float64(oy+m.Size*cellSize)-halfWall)
		screen.DrawImage(r.WallH, wallOpBottom)
	}

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

func loadCellImages() map[maze.CellType]*ebiten.Image {
	return map[maze.CellType]*ebiten.Image{
		maze.Empty:    loadImage("assets/cell_empty.png"),
		maze.Hospital: loadImage("assets/cell_hospital.png"),
		maze.Exit:     loadImage("assets/cell_exit.png"),
		maze.Hole:     loadImage("assets/cell_hole.png"),
		maze.Dragon:   loadImage("assets/cell_dragon.png"),
		maze.Armory:   loadImage("assets/cell_armory.png"),
		maze.River:    loadImage("assets/cell_river.png"),
		maze.Estuary:  loadImage("assets/cell_estuary.png"),
	}
}
