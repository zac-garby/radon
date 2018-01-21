package vm

import (
	"fmt"

	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/object"
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
	f.forwardDeclare()

	for ; f.offset < len(f.code) && f.vm.err == nil && !f.vm.halted; f.offset++ {
		instruction := f.code[f.offset]
		f.handleInterrupts()
		f.do(instruction)
	}
}

func (f *Frame) handleInterrupts() {
	select {
	case i, ok := <-f.vm.Interrupts:
		if !ok {
			break
		}

		switch i {
		case InterruptHalt:
			fmt.Fprintln(f.vm.Out, "vm halted")
			f.vm.halted = true
		}

	default:
		break
	}
}

func (f *Frame) do(i bytecode.Instruction) {
	if f.vm.halted {
		return
	}

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
	}

	return "", false
}

func (f *Frame) searchName(name string) (object.Object, bool) {
	if val, ok := f.store.Get(name); ok {
		return val, true
	}

	return nil, false
}

// Goes through the frame's bytecode and finds pairs
// of LOAD_CONST followed by STORE_NAME where
// LOAD_CONST loads a function constant. It then
// replaces these instructions with dummy bytes and
// preloads the constants.
func (f *Frame) forwardDeclare() {
	for i, instr := range f.code {
		if instr.Code == bytecode.LoadConst {
			// The index of the next instruction
			nextIndex := i + 1

			if nextIndex >= len(f.code) {
				break
			}

			if next := f.code[nextIndex]; next.Code == bytecode.StoreName {
				var (
					name = f.store.Names[next.Arg]
					fn   = f.constants[instr.Arg]
				)

				if !(fn.Type() == object.FunctionType || fn.Type() == object.ModelType) {
					continue
				}

				f.code[i].Code = bytecode.Dummy
				f.code[i].Name = "DUMMY"

				f.code[nextIndex].Code = bytecode.Dummy
				f.code[nextIndex].Name = "DUMMY"

				f.store.Set(name, fn, true)
			}
		}
	}
}
