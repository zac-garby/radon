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
	stores        []*Store
	stack         *Stack
	breaks, nexts []int
	constants     []object.Object
	names         []string
	jumps         []int
}

func (f *Frame) offsetToInstructionIndex(offset int) int {
	var index, counter int

	for _, instr := range f.code {
		if bytecode.Instructions[instr.Code].HasArg {
			counter += 3
		} else {
			counter++
		}

		if counter >= offset {
			return index
		}

		index++
	}

	return index
}

func (f *Frame) getName(arg rune) (string, bool) {
	index := int(arg)
	if index < len(f.names) {
		return f.names[index], true
	}

	return "", false
}

func (f *Frame) searchName(name string) (object.Object, bool) {
	if val, ok := f.store().Get(name); ok {
		return val.Value, true
	}

	return nil, false
}

func (f *Frame) store() *Store {
	return f.stores[0]
}

func (f *Frame) pushStore(s *Store) {
	s.Enclosing = f.store()
	f.stores = append([]*Store{s}, f.stores...)
}

func (f *Frame) popStore() {
	if len(f.stores) == 0 {
		panic("no stores in the store stack")
	}

	f.stores = f.stores[1:]
}
