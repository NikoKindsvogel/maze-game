package maze

type CellType int

const (
	Empty CellType = iota
	Wall
	Hole
	River
	Estuary
	Exit
	Hospital
	Armory
	Dragon
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
	None
)

func Opposite(d Direction) Direction {
	return (d + 2) % 4
}
