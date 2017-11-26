package vm

import (
	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/object"
	"github.com/Zac-Garby/lang/store"
)

// A Frame is a VM frame. One is created for every
// function call and for the main program.
type Frame struct {
	prev          *Frame
	code          bytecode.Code
	offset        uint
	vm            *VM
	store         *store.Store
	stack         stack
	breaks, nexts []uint
	constants     []object.Object
}
