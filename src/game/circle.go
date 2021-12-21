package game

import "fmt"

type Circle struct {
	r   float64
	pos Vec2
}

func NewCircle(r float64, pos Vec2) *Circle {
	return &Circle{
		pos: pos,
		r:   r,
	}
}

func (a *Circle) ComputeAABB() *AABB {
	return NewAABB(a.pos, a.r, a.r)
}

func (a *Circle) SetPosition(newPosition Vec2) {
	a.pos = newPosition
}

func (a *Circle) Position() Vec2 {
	return a.pos
}

func (a *Circle) Clone() Shape {
	return &Circle{
		pos: a.pos,
		r:   a.r,
	}
}

func (a *Circle) CollidesWith(s Shape) (*Collision, bool) {
	switch t := s.(type) {
	case *AABB:
		return NewCvsAABBCollision(t, a)
	case *Circle:
		return NewCvsCCollision(a, t)
	default:
		fmt.Println("Unknown Collision")
		return nil, false
	}
}

func CircCollidesWithCirc(a Circle, b Circle) bool {
	r := a.r + b.r

	xs := a.pos.X + b.pos.X
	ys := a.pos.Y + b.pos.Y

	return r*r < xs*xs+ys*ys
}
