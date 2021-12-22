package game

import (
	"fmt"
	"time"
)

type Move struct {
	dir    Vec2
	player string
}

type Game struct {
	bodies  []*Body
	players map[string]*Player
	mapSize Vec2

	render RenderFunc

	stop  chan bool
	moves chan Move

	forces []ForceFunc
}

func NewGame(render RenderFunc, playerIds []string) *Game {
	g := &Game{
		bodies:  make([]*Body, 0),
		players: make(map[string]*Player),
		mapSize: Vec2{1000, 1000},
		render:  render,
		stop:    make(chan bool),
		moves:   make(chan Move),
		forces:  make([]ForceFunc, 0),
	}

	for i, pId := range playerIds {
		pBody := NewCircleBody(
			Vec2{50.0 * float64(i+1), 200},
			50,
			*NewMass(50), Material{density: 5, restitution: 0.5}, 0)
		g.players[pId] = &Player{
			body:      pBody,
			xVel:      10,
			jumpForce: 50,
		}
		g.bodies = append(g.bodies, pBody)
	}

	return g
}

func (g *Game) StartGame() {
	go g.Loop()
}

func (g *Game) StopGame() {
	g.stop <- true
}

func (g *Game) MakeMove(player string, dir Vec2) {
	g.moves <- Move{player: player, dir: dir}
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
	nextMoves := make(map[string]Vec2, 4)

	//Game Loop
	for {
	read_messages:
		for {
			select {
			case _, _ = <-g.stop:
				return
			case _, ok := <-frameTicker.C:
				if !ok {
					return
				}
				//continue with game loop
				break read_messages
			case m := <-g.moves:
				nextMoves[m.player] = m.dir
			}
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
			g.ApplyMoves(nextMoves)
			UpdatePhysics(dtf, g.bodies, g.forces)
			accumulator -= dt
		}

		alpha := accumulator / dt

		RenderGame(float64(alpha), objects, g.render)
	}
}

func (g *Game) ApplyMoves(moves map[string]Vec2) {
	for pId, dir := range moves {
		p, ok := g.players[pId]
		if !ok {
			fmt.Println("player", pId, "not in players")
		} else {
			p.ApplyMove(dir)
		}
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
