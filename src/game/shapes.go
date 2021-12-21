package game

import "math"

type Shape interface {
	ComputeAABB() *AABB
	CollidesWith(Shape) (*Collision, bool)
	Position() Vec2
	SetPosition(Vec2)
	Clone() Shape
}

func NewCvsCCollision(a *Circle, b *Circle) (*Collision, bool) {
	c := &Collision{
		// a: a, //todo make circle implement object
		// b: b,
	}

	n := b.pos.Minus(a.pos)

	r := a.r + b.r
	r *= r

	if n.LengthSquared() > r {
		return nil, false
	}

	d := n.Length()

	if d != 0 {
		c.penetration = r - d
		c.normal = n.DividedBy(d)
	} else {
		c.penetration = a.r
		c.normal = Vec2{1, 0}
	}
	return c, true
}

func NewAABBvsAABBCollision(a, b *AABB) (*Collision, bool) {
	c := &Collision{
		//todo set a and b
	}

	// Vector from A to B
	n := b.pos.Minus(a.pos)

	// Calculate half extents along x axis for each object
	a_extent := (a.max.X - a.min.X) / 2
	b_extent := (b.max.X - b.min.X) / 2

	// Calculate overlap on x axis
	x_overlap := a_extent + b_extent - math.Abs(n.X)

	// SAT test on x axis
	if x_overlap > 0 {
		// Calculate half extents along y axis for each object
		a_extent = (a.max.Y - a.min.Y) / 2
		b_extent = (b.max.Y - b.min.Y) / 2

		// Calculate overlap on y axis
		y_overlap := a_extent + b_extent - math.Abs(n.Y)

		// SAT test on y axis
		if y_overlap > 0 {
			if x_overlap > y_overlap {
				// Point towards B knowing that n points from A to B
				if n.X < 0 {
					c.normal = Vec2{-1, 0}
				} else {
					c.normal = Vec2{1, 0}
				}
				c.penetration = x_overlap
			} else {
				if n.Y < 0 {
					c.normal = Vec2{0, -1}
				} else {
					c.normal = Vec2{0, 1}
				}
				c.penetration = y_overlap
			}
		}
	}

	return c, true
}

func NewCvsAABBCollision(a *AABB, b *Circle) (*Collision, bool) {
	c := &Collision{
		//todo set a and b
	}

	n := b.pos.Minus(a.pos)
	closest := n

	x_extent := (a.max.X - a.min.X) / 2
	y_extent := (a.max.Y - a.min.Y) / 2

	closest.X = Clamp(-x_extent, x_extent, closest.X)
	closest.Y = Clamp(-y_extent, y_extent, closest.Y)

	inside := false
	// Circle is inside the AABB, so we need to clamp the circle's center
	// to the closest edge
	if n == closest {
		inside = true

		//find the closest  axis
		if math.Abs(n.X) > math.Abs(n.Y) {
			// clamp to closest extent
			if closest.X > 0 {
				closest.X = x_extent
			} else {
				closest.X = -x_extent
			}
		} else {
			if closest.Y > 0 {
				closest.Y = y_extent
			} else {
				closest.Y = -y_extent
			}
		}
	}

	normal := n.Minus(closest)
	d := normal.LengthSquared()

	if d > b.r*b.r && !inside {
		return nil, false
	}

	d = normal.Length()

	// Collision normal needs to be flipped to point outside if circle was
	// inside the AABB
	if inside {
		c.normal = n.Times(-1)
		c.penetration = b.r - d
	} else {
		c.normal = n
		c.penetration = b.r - d
	}
	return c, true
}
