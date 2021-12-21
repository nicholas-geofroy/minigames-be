package game

import "fmt"

type AABB struct {
	min Vec2
	max Vec2

	pos Vec2

	//half of the width and height
	hDim Vec2
}

func NewAABB(pos Vec2, width float64, height float64) *AABB {
	hDim := Vec2{width / 2, height / 2}
	return &AABB{
		pos:  pos,
		min:  pos.Minus(hDim),
		max:  pos.Plus(hDim),
		hDim: hDim,
	}
}

func (a *AABB) SetPosition(newPosition Vec2) {
	a.pos = newPosition
	a.max = newPosition.Plus(a.hDim)
	a.min = newPosition.Minus(a.hDim)
}

func (a *AABB) Position() Vec2 {
	return a.pos
}

func (a *AABB) ComputeAABB() *AABB {
	return a
}

func (a *AABB) Clone() Shape {
	return &AABB{
		min:  a.min,
		max:  a.max,
		pos:  a.pos,
		hDim: a.hDim,
	}
}

func (a *AABB) CollidesWith(s Shape) (*Collision, bool) {
	switch t := s.(type) {
	case *AABB:
		return NewAABBvsAABBCollision(a, t)
	case *Circle:
		return NewCvsAABBCollision(a, t)
	default:
		fmt.Println("Unknown Collision")
		return nil, false
	}
}

func AABBCollidesWithAABB(a *AABB, b *AABB) bool {
	if a.max.X < b.min.X || b.max.X < a.min.X {
		return false
	}
	if a.max.Y < b.min.Y || b.max.Y < a.min.Y {
		return false
	}

	return true
}
