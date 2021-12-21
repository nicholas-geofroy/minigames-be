package lobby

import "minigames-be/src/game"

type OutMsg struct {
	MsgType int
	Data    interface{}
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

type StartMsg struct{}

func (e *StartMsg) Type() int {
	return 3
}

const (
	left  = 0
	right = 1
)

type MoveMsg struct {
	direction int
}

func (e *MoveMsg) Type() int {
	return 4
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
