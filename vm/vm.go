package vm

import (
	"github.com/Zac-Garby/lang/object"
)

// A VM interprets and executes bytecode.
type VM struct {
	frames    []*Frame
	frame     *Frame
	returnVal object.Object
}

// New creates a new virtual machine.
func New() *VM {
	return &VM{
		frames:    make([]*Frame, 0),
		frame:     nil,
		returnVal: nil,
	}
}
