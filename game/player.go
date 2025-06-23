package game

type Player struct {
	ID          string
	Row, Col    int
	Hurt        bool
	HasTreasure bool
	Bullet      bool
}
