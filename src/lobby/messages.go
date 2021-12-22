package lobby

import "minigames-be/src/game"

type OutMsg struct {
	MsgType int
	Data    interface{}
}

type ClientMsg struct {
	clientId string
	msg      LobbyMsg
}

type LobbyMsg interface {
	Type() int
}

type ErrorMsg struct {
	errorType int
	msg       string
}

func (e *ErrorMsg) Type() int {
	return 1
}

type MembersMsg struct {
	members []string
}

func (e *MembersMsg) Type() int {
	return 2
}

const startMsgType = 3

type StartMsg struct{}

func (e *StartMsg) Type() int {
	return startMsgType
}

const (
	noDir       = 0
	left        = 1
	right       = 2
	moveMsgType = 4
)

type MoveMsg struct {
	Direction float64
}

func (e *MoveMsg) Type() int {
	return moveMsgType
}

func (e *MoveMsg) ToVec() game.Vec2 {
	v := game.Vec2{X: 0, Y: 0}

	switch e.Direction {
	case left:
		v.X = -1
	case right:
		v.X = 1
	}
	return v
}

type JumpMsg struct{}

func (e *JumpMsg) Type() int {
	return 5
}

type ObjState struct {
	Pos     game.Vec2
	ObjType string
}
type StateMsg struct {
	Objects []ObjState
}

func (e *StateMsg) Type() int {
	return 6
}
