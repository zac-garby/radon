package vm

import (
	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/object"
)

// A Frame is a VM frame. One is created for every
// function call and for the main program.
type Frame struct {
	prev          *Frame
	code          bytecode.Code
	offset        int
	vm            *VM
	store         *Store
	stack         *stack
	breaks, nexts []int
	constants     []object.Object
}

func (f *Frame) execute() {
	for ; f.offset < len(f.code) && f.vm.err == nil; f.offset++ {
		instruction := f.code[f.offset]
		f.do(instruction)
	}
}

func (f *Frame) do(i bytecode.Instruction) {
	eff, ok := effectors[i.Code]
	if !ok {
		f.vm.err = Errf("execution of instruction %s not implemented", ErrNoInstruction, i.Name)
		return
	}

	eff(f, i)
}

func (f *Frame) offsetToInstructionIndex(o int) int {
	var index, counter int

	for _, instr := range f.code {
		if bytecode.Instructions[instr.Code].HasArg {
			counter += 3
		} else {
			counter++
		}

		if counter >= o {
			return index
		}

		index++
	}

	return index
}

func (f *Frame) getName(arg rune) (string, bool) {
	index := int(arg)

	if index < len(f.store.Names) {
		name := f.store.Names[index]
		return name, true
	} else if f.prev != nil {
		return f.prev.getName(arg)
	}

	return "", false
}

func (f *Frame) searchName(name string) (object.Object, bool) {
	if val, ok := f.store.Get(name); ok {
		return val, true
	} else if f.prev != nil {
		return f.prev.searchName(name)
	}

	return nil, false
}
