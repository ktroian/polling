package action

import (
	"unsafe"
)

var (
	actions chan Action
)

type Action struct {
	Target      unsafe.Pointer
	Type 		int
}

func Init() {
	if actions == nil {
		actions = make(chan Action)
	}
}

func Add(actionType int, target unsafe.Pointer) {
	actions <- Action{target, actionType}
}

func Pop() Action {
	return <- actions 
}