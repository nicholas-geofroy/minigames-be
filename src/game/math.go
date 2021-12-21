package game

import "math"

func Clamp(min, max, val float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

type Vec2 struct {
	X float64
	Y float64
}

func (v Vec2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v Vec2) Plus(o Vec2) Vec2 {
	return Vec2{
		v.X + o.X,
		v.Y + o.Y,
	}
}

func (a Vec2) Minus(b Vec2) Vec2 {
	return Vec2{
		a.X - b.X,
		a.Y - b.Y,
	}
}

func (a Vec2) Dot(b Vec2) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vec2) Times(b float64) Vec2 {
	return Vec2{
		a.X * b,
		a.Y * b,
	}
}

func (a Vec2) DividedBy(b float64) Vec2 {
	return Vec2{
		a.X / b,
		a.Y / b,
	}
}
