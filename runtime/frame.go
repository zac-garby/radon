package runtime

import (
	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/object"
)

// A Frame is created for each function call, and also one for the main program. It contains
// the bytecode to execute, along with the frame's constants and names, and other data.
type Frame struct {
	prev          *Frame
	code          bytecode.Code
	offset        int
	vm            *VM
	store         *Store
	stack         *stack
	breaks, nexts []int
	constants     []object.Object
	names         []string
}
