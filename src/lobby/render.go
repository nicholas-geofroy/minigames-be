package lobby

import (
	"minigames-be/src/game"
)

func GetShapeType(s game.Shape) string {
	switch s.(type) {
	case *game.AABB:
		return "Rectangle"
	case *game.Circle:
		return "Circle"
	default:
		return ""
	}
}

func MessageRender(objects []game.Object) *StateMsg {
	outMsg := &StateMsg{
		Objects: make([]ObjState, len(objects)),
	}

	for i, obj := range objects {
		outMsg.Objects[i] = ObjState{
			Pos:     obj.GetShape().Position(),
			ObjType: GetShapeType(obj.GetShape()),
		}
	}

	return outMsg
}
