package vm

import (
	"github.com/Zac-Garby/lang/object"
)

type stack struct {
	objects []object.Object
}

func newStack() *stack {
	return &stack{
		objects: []object.Object{},
	}
}
