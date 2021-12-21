package game

type Point struct {
	x int
	y int
}

func NewPoint(x int, y int) Point {
	return Point{
		x: x,
		y: y,
	}
}

type Player struct {
	body *Body
}

func NewPlayer(id int, body *Body) *Player {
	return &Player{
		body: body,
	}
}

func (p *Player) MoveLeft() {
}

func (p *Player) MoveRight() {
}

func (p *Player) Jump() {

}
