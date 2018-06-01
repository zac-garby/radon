package runtime

import (
	"fmt"

	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/object"
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

	Effectors[bytecode.LoadName] = func(v *VM, f *Frame, arg rune) error {
		val, ok := f.store().Get(f.names[arg])
		if !ok {
			return fmt.Errorf("variable %s doesn't exist in the current scope", f.names[arg])
		}

		return f.stack.Push(val.Value)
	}

	Effectors[bytecode.StoreName] = func(v *VM, f *Frame, arg rune) error {
		val, err := f.stack.Pop()
		if err != nil {
			return err
		}

		f.store().Set(f.names[arg], val, false)
		return nil
	}

	Effectors[bytecode.DeclareName] = func(v *VM, f *Frame, arg rune) error {
		val, err := f.stack.Pop()
		if err != nil {
			return err
		}

		f.store().Set(f.names[arg], val, false)
		return nil
	}

	Effectors[bytecode.LoadSubscript] = func(v *VM, f *Frame, arg rune) error {
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

	Effectors[bytecode.UnaryInvert] = unaryEffector("!")
	Effectors[bytecode.UnaryNegate] = unaryEffector("-")
	Effectors[bytecode.UnaryTuple] = unaryEffector(",")

	Effectors[bytecode.BinaryAdd] = binaryEffector("+")
	Effectors[bytecode.BinarySub] = binaryEffector("-")
	Effectors[bytecode.BinaryMul] = binaryEffector("*")
	Effectors[bytecode.BinaryDiv] = binaryEffector("/")
	Effectors[bytecode.BinaryExp] = binaryEffector("^")
	Effectors[bytecode.BinaryFloorDiv] = binaryEffector("//")
	Effectors[bytecode.BinaryMod] = binaryEffector("%")
	Effectors[bytecode.BinaryLogicOr] = binaryEffector("||")
	Effectors[bytecode.BinaryLogicAnd] = binaryEffector("&&")
	Effectors[bytecode.BinaryBitOr] = binaryEffector("|")
	Effectors[bytecode.BinaryBitAnd] = binaryEffector("&")
	Effectors[bytecode.BinaryEqual] = equalityEffector(true)
	Effectors[bytecode.BinaryNotEqual] = equalityEffector(false)
	Effectors[bytecode.BinaryLess] = binaryEffector("<")
	Effectors[bytecode.BinaryMore] = binaryEffector(">")
	Effectors[bytecode.BinaryLessEq] = binaryEffector("<=")
	Effectors[bytecode.BinaryMoreEq] = binaryEffector(">=")
	Effectors[bytecode.BinaryTuple] = binaryEffector(",")

	Effectors[bytecode.CallFunction] = func(v *VM, f *Frame, argCount rune) error {
		top, err := f.stack.Pop()
		if err != nil {
			return err
		}

		fn, ok := top.(*object.Function)
		if !ok {
			return fmt.Errorf("cannot call non-function type %s", top.Type())
		}

		if int(argCount) != len(fn.Parameters) {
			return fmt.Errorf("wrong amount of arguments passed a function. expected %d, got %d", len(fn.Parameters), argCount)
		}

		store := NewStore(f.store())

		for _, param := range fn.Parameters {
			arg, err := f.stack.Pop()
			if err != nil {
				return err
			}

			store.Set(param, arg, true)
		}

		if fn.Self != nil {
			store.Set("self", fn.Self, true)
		}

		frame := &Frame{
			prev:      f,
			code:      fn.Code,
			offset:    0,
			vm:        v,
			stores:    []*Store{store},
			stack:     NewStack(),
			constants: fn.Constants,
			names:     fn.Names,
			jumps:     fn.Jumps,
		}

		f.vm.PushFrame(frame)

		return nil
	}

	Effectors[bytecode.Return] = func(v *VM, f *Frame, arg rune) error {
		f.offset = len(f.code) - 1
		return nil
	}

	Effectors[bytecode.PushScope] = func(v *VM, f *Frame, arg rune) error {
		f.pushStore(v.storePool.Release(f.store()))
		return nil
	}

	Effectors[bytecode.PopScope] = func(v *VM, f *Frame, arg rune) error {
		f.popStore()
		return nil
	}
}

func equalityEffector(shouldEqual bool) Effector {
	return func(v *VM, f *Frame, arg rune) error {
		right, err := f.stack.Pop()
		if err != nil {
			return err
		}

		left, err := f.stack.Pop()
		if err != nil {
			return err
		}

		if left.Equals(right) == shouldEqual {
			return f.stack.Push(&object.Boolean{Value: true})
		} else {
			return f.stack.Push(&object.Boolean{Value: false})
		}
	}
}

func unaryEffector(op string) Effector {
	return func(v *VM, f *Frame, arg rune) error {
		obj, err := f.stack.Pop()
		if err != nil {
			return err
		}

		result, ok := obj.Prefix(op)
		if !ok {
			return fmt.Errorf("could not apply prefix operator %s to %s", op, obj.String())
		}

		return f.stack.Push(result)
	}
}

func binaryEffector(op string) Effector {
	return func(v *VM, f *Frame, arg rune) error {
		right, err := f.stack.Pop()
		if err != nil {
			return err
		}

		left, err := f.stack.Pop()
		if err != nil {
			return err
		}

		result, ok := left.Infix(op, right)
		if !ok {
			return fmt.Errorf("could not apply infix operator %s between %s and %s", op, left.String(), right.String())
		}

		return f.stack.Push(result)
	}
}
