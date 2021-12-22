package game

import "math"

type Body struct {
	shape        Shape //either AABB or Circle
	force        Vec2
	vel          Vec2
	massData     MassData
	material     Material
	gravityScale float64
}

func NewBoxBody(pos Vec2, width float64, height float64, massData MassData, material Material, gravityScale float64) *Body {
	box := NewAABB(pos, width, height)

	return &Body{
		shape:        box,
		force:        Vec2{0, 0},
		vel:          Vec2{0, 0},
		massData:     massData,
		material:     material,
		gravityScale: gravityScale,
	}
}

func NewCircleBody(pos Vec2, radius float64, massData MassData, material Material, gravityScale float64) *Body {
	c := NewCircle(radius, pos)

	return &Body{
		shape:        c,
		force:        Vec2{0, 0},
		vel:          Vec2{0, 0},
		massData:     massData,
		material:     material,
		gravityScale: gravityScale,
	}
}

func (b *Body) UpdateKinematics(dt float64) {
	b.vel = b.vel.Plus(b.force.Times(b.massData.invMass * dt))
	b.shape.SetPosition(b.shape.Position().Plus(b.vel.Times(dt)))
}

func (b *Body) Clone() *Body {
	return &Body{
		shape:        b.shape.Clone(),
		force:        b.force,
		vel:          b.vel,
		massData:     b.massData,
		material:     b.material,
		gravityScale: b.gravityScale,
	}
}

type MassData struct {
	mass    float64
	invMass float64
}

type Material struct {
	density     float64
	restitution float64
}

func NewMass(mass float64) *MassData {
	var invMas float64 = 0
	if mass != 0 {
		invMas = 1 / mass
	}
	return &MassData{
		mass:    mass,
		invMass: invMas,
	}
}

type Pair struct {
	A Body
	B Body
}

func GeneratePairs(bodies []*Body) []*Pair {
	pairs := make([]*Pair, 0)

	for i, bodyA := range bodies {
		for j := i + 1; j < len(bodies); j++ {
			bodyB := bodies[j]
			if bodyA == bodyB {
				continue
			}

			a := bodyA.shape.ComputeAABB()
			b := bodyB.shape.ComputeAABB()

			if AABBCollidesWithAABB(a, b) {
				pairs = append(pairs, &Pair{*bodyA, *bodyB})
			}
		}
	}

	return pairs
}

type ForceFunc func(*Body)

func UpdatePhysics(dt float64, bodies []*Body, forces []ForceFunc) {
	//determine forces on each body
	for _, body := range bodies {
		ApplyGravity(body)
		for _, force := range forces {
			force(body)
		}
	}
	//collision forces
	pairs := GeneratePairs(bodies)
	for _, pair := range pairs {
		col, isCol := pair.A.shape.CollidesWith(pair.B.shape)
		if !isCol {
			continue
		}

		ResolveCollision(col)
		PositionalCorrection(col)
	}

	//update kinematics
	for _, body := range bodies {
		body.UpdateKinematics(dt)
	}
}

func ApplyGravity(b *Body) {
	b.force = b.force.Plus(Vec2{0, 9.8 * b.gravityScale})
}

func ResolveCollision(c *Collision) {
	// Calculate relative velocity
	a, b := c.a, c.b
	relV := b.vel.Minus(a.vel)

	// Calculate relative velocity in terms of the normal direction
	vecAlongNormal := relV.Dot(c.normal)

	// Do not resolve if velocities are separating
	if vecAlongNormal > 0 {
		return
	}

	// Calculate restitution
	e := math.Min(a.material.restitution, b.material.restitution)

	// Calculate impulse scalar
	j := -(1 + e) * vecAlongNormal
	j = j / (a.massData.invMass + b.massData.invMass)

	// Apply impulse
	// TODO: if this is wrong need to do fraction stuff
	a.vel = a.vel.Minus(c.normal.Times(j * a.massData.invMass))
	b.vel = b.vel.Plus(c.normal.Times(j * b.massData.invMass))
}

func PositionalCorrection(c *Collision) {
	const percent float64 = 0.2 // usually 20-80%
	a, b := c.a, c.b

	correction := c.normal.Times(c.penetration / (a.massData.invMass + b.massData.invMass) * percent)
	a.shape.SetPosition(a.shape.Position().Minus(correction.Times(a.massData.invMass)))
	b.shape.SetPosition(b.shape.Position().Minus(correction.Times(b.massData.invMass)))
}

type Collision struct {
	a           Body
	b           Body
	penetration float64
	normal      Vec2
}
