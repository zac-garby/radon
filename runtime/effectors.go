package runtime

import (
	"fmt"

	"github.com/Zac-Garby/radon/bytecode"
)

// An Effector is the function which is called for a particular instruction.
type Effector func(v *VM, f *Frame, arg rune) error

// Effectors stores the effector function for each instruction. The effector
// for an instruction n can be accessed via Effectors[n].
var Effectors [256]Effector

// Actually define the instruction effectors. These functions will assume that the necessary
// data is present, for example
func init() {
	Effectors[bytecode.Nop] = func(v *VM, f *Frame, arg rune) error { return nil }
	Effectors[bytecode.NopArg] = func(v *VM, f *Frame, arg rune) error { return nil }
	Effectors[bytecode.LoadConst] = func(v *VM, f *Frame, arg rune) error { return f.stack.Push(f.constants[arg]) }
	Effectors[bytecode.LoadName] = func(v *VM, f *Frame, arg rune) error { return f.store.Get(f.names[arg]) }

	Effectors[bytecode.StoreName] = func(v *VM, f *Frame, arg rune) error {
		val, err := f.stack.Pop()
		if err != nil {
			return err
		}

		return f.store.Set(f.names[arg], val, false)
	}

	Effectors[bytecode.DeclareName] = func(v *VM, f *Frame, arg rune) error {
		val, err := f.stack.Pop()
		if err != nil {
			return err
		}

		return f.store.Set(f.names[arg], val, false)
	}

	Effectors[bytecode.LoadSubcript] = func(v *VM, f *Frame, arg rune) error {
		index, err := f.stack.Pop()
		if err != nil {
			return err
		}

		obj, err := f.stack.Pop()
		if err != nil {
			return err
		}

		result, ok := obj.Subscript(index)
		if !ok {
			return fmt.Errorf("could not subscript a %s object with index %s", obj.Type(), index.String())
		}

		return f.stack.Push(result)
	}

	Effectors[bytecode.StoreSubscript] = func(v *VM, f *Frame, arg rune) error {
		index, err := f.stack.Pop()
		if err != nil {
			return err
		}

		obj, err := f.stack.Pop()
		if err != nil {
			return err
		}

		val, err := f.stack.Pop()
		if err != nil {
			return err
		}

		if !obj.SetSubscript(index, val) {
			return fmt.Errorf("could not set subscript %s on a %s object", index.String(), obj.Type())
		}

		return nil
	}
}
