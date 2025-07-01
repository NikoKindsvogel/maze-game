package ebiten_ui

import (
	"image"
	"image/color"
	"log"
	"maze-game/maze"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	cellSize   = 64 // width and height of each cell image
	wallOffset = 8  // half-width for wall overlap
	spacing    = 20 // space between the two mazes
)

type RevealScreen struct {
	StartMaze *maze.Maze
	FinalMaze *maze.Maze
	Images    map[maze.CellType]*ebiten.Image
	WallH     *ebiten.Image
	WallV     *ebiten.Image
	Treasure  *ebiten.Image
}

func NewRevealScreen(start, final *maze.Maze) *RevealScreen {
	images := loadCellImages()
	// wallH := loadImage("assets/wall_horizontal.png")
	// wallV := loadImage("assets/wall_vertical.png")
	// treasure := loadImage("assets/treasure.png")

	// for testing:
	wallH := createWallImage(cellSize, 8, color.RGBA{0, 0, 0, 255}) // black horizontal wall
	wallV := createWallImage(8, cellSize, color.RGBA{0, 0, 0, 255}) // black vertical wall
	treasure := createColoredImage(color.RGBA{255, 215, 0, 128})    // semi-transparent gold

	return &RevealScreen{
		StartMaze: start,
		FinalMaze: final,
		Images:    images,
		WallH:     wallH,
		WallV:     wallV,
		Treasure:  treasure,
	}
}

func (r *RevealScreen) Update() {
	// Optional: support switching views
}

func (r *RevealScreen) Draw(screen *ebiten.Image) {
	offsetX := 50
	offsetY := 50

	r.drawMaze(screen, r.StartMaze, offsetX, offsetY)
	r.drawMaze(screen, r.FinalMaze, offsetX+cellSize*r.StartMaze.Size+spacing, offsetY)
}

func (r *RevealScreen) drawMaze(screen *ebiten.Image, m *maze.Maze, ox, oy int) {
	for row := 0; row < m.Size; row++ {
		for col := 0; col < m.Size; col++ {
			cell := m.Grid[row][col]
			x := ox + col*cellSize
			y := oy + row*cellSize

			img := r.Images[cell.Type]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(img, op)

			// Draw walls
			if cell.Walls[maze.Right] && col < m.Size-1 {
				wallOp := &ebiten.DrawImageOptions{}
				wallOp.GeoM.Translate(float64(x+cellSize-wallOffset), float64(y))
				screen.DrawImage(r.WallV, wallOp)
			}
			if cell.Walls[maze.Down] && row < m.Size-1 {
				wallOp := &ebiten.DrawImageOptions{}
				wallOp.GeoM.Translate(float64(x), float64(y+cellSize-wallOffset))
				screen.DrawImage(r.WallH, wallOp)
			}

			// Treasure overlay
			if m.TreasureOnMap && m.TreasureRow == row && m.TreasureCol == col {
				tOp := &ebiten.DrawImageOptions{}
				tOp.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(r.Treasure, tOp)
			}
		}
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

// func loadCellImages() map[maze.CellType]*ebiten.Image {
// 	return map[maze.CellType]*ebiten.Image{
// 		maze.Empty:    loadImage("assets/cell_empty.png"),
// 		maze.Hospital: loadImage("assets/cell_hospital.png"),
// 		maze.Exit:     loadImage("assets/cell_exit.png"),
// 		maze.Hole:     loadImage("assets/cell_hole.png"),
// 		maze.Dragon:   loadImage("assets/cell_dragon.png"),
// 		maze.Armory:   loadImage("assets/cell_armory.png"),
// 		maze.River:    loadImage("assets/cell_river.png"),
// 		maze.Estuary:  loadImage("assets/cell_estuary.png"),
// 	}
// }

// for testing:

func loadCellImages() map[maze.CellType]*ebiten.Image {
	return map[maze.CellType]*ebiten.Image{
		maze.Empty:    createColoredImage(color.RGBA{200, 200, 200, 255}), // light gray
		maze.Hospital: createColoredImage(color.RGBA{255, 100, 100, 255}), // red
		maze.Exit:     createColoredImage(color.RGBA{0, 255, 0, 255}),     // green
		maze.Hole:     createColoredImage(color.RGBA{50, 50, 50, 255}),    // dark gray
		maze.Dragon:   createColoredImage(color.RGBA{150, 0, 0, 255}),     // dark red
		maze.Armory:   createColoredImage(color.RGBA{0, 0, 255, 255}),     // blue
		maze.River:    createColoredImage(color.RGBA{0, 150, 255, 255}),   // cyan
		maze.Estuary:  createColoredImage(color.RGBA{0, 255, 255, 255}),   // aqua
	}
}

func createColoredImage(color color.Color) *ebiten.Image {
	img := ebiten.NewImage(cellSize, cellSize)
	img.Fill(color)
	return img
}

func createWallImage(width, height int, col color.Color) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	img.Fill(col)
	return img
}
