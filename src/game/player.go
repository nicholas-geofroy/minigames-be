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
	body      *Body
	xVel      float64
	jumpForce float64
	canJump   bool
}

func NewPlayer(id int, body *Body) *Player {
	return &Player{
		body: body,
	}
}

func (p *Player) ApplyMove(dir Vec2) {
	p.body.vel.X = dir.X * p.xVel
}

func (p *Player) OnCollision(c *Collision) {
	p.canJump = false
	if c.normal.Y < 0 {
		p.canJump = true
	}
}
