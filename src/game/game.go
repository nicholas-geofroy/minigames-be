package game

import (
	"fmt"
	"time"
)

type Game struct {
	bodies  []*Body
	players map[string]*Player
	mapSize Vec2

	render RenderFunc
	stop   chan bool
}

func NewGame(render RenderFunc) *Game {
	g := &Game{
		bodies:  make([]*Body, 0),
		players: make(map[string]*Player),
		mapSize: Vec2{1000, 1000},
		render:  render,
	}

	g.bodies = append(g.bodies, NewCircleBody(
		Vec2{200, 200},
		50,
		*NewMass(50), Material{density: 5, restitution: 0.5}, 0.01))

	return g
}

func (g *Game) StartGame() {
	go g.Loop()
}

func (g *Game) StopGame() {
	g.stop <- true
}

func (g *Game) Loop() {
	fmt.Println("Loop start")
	const fps = 50
	const dt = time.Second / fps
	const maxAccumulator = dt * 2
	const dtf = float64(dt) / float64(time.Second)
	var accumulator time.Duration = 0

	var objects = make([]*InternalObj, len(g.bodies))
	for i, body := range g.bodies {
		objects[i] = &InternalObj{
			body:    body,
			lastPos: body.shape.Position(),
			intPos:  body.shape.Position(),
		}
	}

	frameStart := time.Now()
	frameTicker := time.NewTicker(dt)

	defer frameTicker.Stop()

	//Game Loop
	for {
		select {
		case _, ok := <-g.stop:
			if ok {
				return
			}
		case _, ok := <-frameTicker.C:
			if !ok {
				return
			}
			//continue with loop
		}
		curTime := time.Now()

		accumulator += curTime.Sub(frameStart)
		frameStart = curTime

		// Avoid spiral of death and clamp dt, thus clamping
		// how many times the UpdatePhysics can be called in
		// a single game loop.
		if accumulator > maxAccumulator {
			accumulator = maxAccumulator
		}

		for accumulator > dt {
			UpdatePhysics(dtf, g.bodies)
			accumulator -= dt
		}

		alpha := accumulator / dt

		RenderGame(float64(alpha), objects, g.render)
	}
}

type InternalObj struct {
	body    *Body
	lastPos Vec2
	intPos  Vec2 //position to use for rendering
}

func (o *InternalObj) GetShape() Shape {
	intShape := o.body.shape.Clone()
	intShape.SetPosition(o.intPos)
	return intShape
}

type Object interface {
	GetShape() Shape
}

type RenderFunc = func([]Object)

func RenderGame(alpha float64, intObjs []*InternalObj, render RenderFunc) {
	objs := make([]Object, len(intObjs))
	for i, obj := range intObjs {
		interpolatedPos := obj.lastPos.Times(alpha).Plus(obj.body.shape.Position().Times(1.0 - alpha))
		obj.lastPos = obj.body.shape.Position()
		obj.intPos = interpolatedPos

		objs[i] = obj
	}
	render(objs)
}
