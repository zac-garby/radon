package runtime

import (
	"math"

	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/object"
	"github.com/cnf/structhash"
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
			return makeError(NameError, "variable %s doesn't exist in the current scope", f.names[arg])
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
			return makeError(TypeError, "could not subscript a %s object with index %s", obj.Type(), index.String())
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
			return makeError(TypeError, "could not set subscript %s on a %s object", index.String(), obj.Type())
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

		builtin, ok := top.(*object.Builtin)
		if ok {
			args := make([]object.Object, 0, argCount)

			for i := 0; i < int(argCount); i++ {
				arg, err := f.stack.Pop()
				if err != nil {
					return err
				}

				args = append(args, arg)
			}

			result, errorType, errorMessage := builtin.Fn(args...)
			if err != nil {
				return err
			}

			if errorType != "" {
				return makeError(ErrorType(errorType), errorMessage)
			}

			return f.stack.Push(result)
		}

		if fn, ok := top.(*object.Function); ok {
			return callFunction(v, f, fn, argCount)
		} else if items, ok := top.Items(); ok {
			return indexCollection(v, f, items, argCount)
		} else {
			return makeError(TypeError, "cannot call an object of type %s", top.Type())
		}
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
		sto := f.popStore()
		v.storePool.Add(sto)
		return nil
	}

	Effectors[bytecode.Export] = func(v *VM, f *Frame, arg rune) error {
		name := f.names[int(arg)]

		variable, ok := f.store().Get(name)
		if !ok {
			return makeError(NameError, "can't export non-existent variable %s", name)
		}

		val := variable.Value

		enclosing := f.store().Enclosing

		if enclosing == nil {
			return makeError(StructureError, "can't export variable %s from a top-level scope", name)
		}

		enclosing.Set(name, val, true)
		return nil
	}

	Effectors[bytecode.Jump] = func(v *VM, f *Frame, arg rune) error {
		jump := f.offsetToInstructionIndex(f.jumps[int(arg)])
		f.offset = jump
		return nil
	}

	Effectors[bytecode.JumpIf] = func(v *VM, f *Frame, arg rune) error {
		top, err := f.stack.Pop()
		if err != nil {
			return err
		}

		if object.IsTruthy(top) {
			return Effectors[bytecode.Jump](v, f, arg)
		}

		return nil
	}

	Effectors[bytecode.JumpUnless] = func(v *VM, f *Frame, arg rune) error {
		top, err := f.stack.Pop()
		if err != nil {
			return err
		}

		if !object.IsTruthy(top) {
			return Effectors[bytecode.Jump](v, f, arg)
		}

		return nil
	}

	Effectors[bytecode.StartMatch] = func(v *VM, f *Frame, arg rune) error {
		top, err := f.stack.Pop()
		if err != nil {
			return err
		}

		f.matchInputs = append(f.matchInputs, top)

		hasEnd := false

		for i := f.offset; i < len(f.code); i++ {
			if f.code[i].Code == bytecode.EndMatch {
				f.breaks = append(f.breaks, i)
				hasEnd = true
				break
			}
		}

		if !hasEnd {
			return makeError(InternalError, "malformed bytecode -- match expression has no END_MATCH")
		}

		return nil
	}

	Effectors[bytecode.EndMatch] = func(v *VM, f *Frame, arg rune) error {
		if len(f.matchInputs) == 0 {
			return makeError(InternalError, "malformed bytecode -- END_MATCH found (likely) before START_MATCH")
		}

		f.matchInputs = f.matchInputs[:len(f.matchInputs)-1]

		return nil
	}

	Effectors[bytecode.StartBranch] = func(v *VM, f *Frame, arg rune) error {
		cond, err := f.stack.Pop()
		if err != nil {
			return err
		}

		if len(f.matchInputs) == 0 {
			return makeError(InternalError, "unexpected empty match input stack")
		}

		input := f.matchInputs[len(f.matchInputs)-1]

		if !cond.Equals(input) {
			hasEnd := false

			for i := f.offset; i < len(f.code); i++ {
				if f.code[i].Code == bytecode.EndMatch {
					return makeError(InternalError, "malformed bytecode -- END_BRANCH appears after END_MATCH")
				}

				if f.code[i].Code == bytecode.EndBranch {
					f.offset = i + 1
					hasEnd = true
					break
				}
			}

			if !hasEnd {
				return makeError(InternalError, "malformed bytecode -- match branch has no END_BRANCH")
			}
		}

		return nil
	}

	Effectors[bytecode.EndBranch] = func(v *VM, f *Frame, arg rune) error {
		if len(f.breaks) == 0 {
			return makeError(InternalError, "malformed bytecode -- END_BRANCH found outside match")
		}

		top := f.breaks[len(f.breaks)-1]
		f.offset = top

		return nil
	}

	Effectors[bytecode.Break] = func(v *VM, f *Frame, arg rune) error {
		if len(f.breaks) == 0 {
			return makeError(InternalError, "break statements are only valid inside loops and matches")
		}

		// Pop the scope
		sto := f.popStore()
		v.storePool.Add(sto)

		top := f.breaks[len(f.breaks)-1]
		f.offset = top

		return nil
	}

	Effectors[bytecode.Next] = func(v *VM, f *Frame, arg rune) error {
		if len(f.nexts) == 0 {
			return makeError(InternalError, "next statements are only valid inside loops")
		}

		// Pop the scope
		sto := f.popStore()
		v.storePool.Add(sto)

		top := f.nexts[len(f.nexts)-1]
		f.offset = top

		return nil
	}

	Effectors[bytecode.StartLoop] = func(v *VM, f *Frame, arg rune) error {
		f.nexts = append(f.nexts, f.offset)

		var o int

		for o = f.offset; f.code[o].Code != bytecode.EndLoop; o++ {
		}

		f.breaks = append(f.breaks, o)

		return nil
	}

	Effectors[bytecode.EndLoop] = func(v *VM, f *Frame, arg rune) error {
		f.breaks = f.breaks[:len(f.breaks)-1]
		f.nexts = f.nexts[:len(f.nexts)-1]

		return nil
	}

	Effectors[bytecode.PushIter] = func(v *VM, f *Frame, arg rune) error {
		top, err := f.stack.Pop()
		if err != nil {
			return err
		}

		iter, ok := top.Iter()
		if !ok {
			return makeError(TypeError, "cannot make an iterable from a %s", top.Type())
		}

		f.iterStack = append(f.iterStack, iter)

		return nil
	}

	Effectors[bytecode.PopIter] = func(v *VM, f *Frame, arg rune) error {
		if len(f.iterStack) == 0 {
			return makeError(InternalError, "no iters to pop from the iter stack")
		}

		f.iterStack = f.iterStack[:len(f.iterStack)-1]

		return nil
	}

	Effectors[bytecode.AdvIterFor] = func(v *VM, f *Frame, arg rune) error {
		if len(f.iterStack) == 0 {
			return makeError(InternalError, "no iters to advance from the iter stack")
		}

		iter := f.iterStack[len(f.iterStack)-1]

		val, ok := iter.Next()
		if !ok {
			return Effectors[bytecode.Break](v, f, arg)
		}

		name, ok := f.getName(arg)
		if !ok {
			return makeError(InternalError, "cannot get name: %d", arg)
		}

		f.store().Set(name, val, true)

		return nil
	}

	Effectors[bytecode.MakeList] = func(v *VM, f *Frame, arg rune) error {
		elems := make([]object.Object, arg)

		for n := int(arg) - 1; n >= 0; n-- {
			top, err := f.stack.Pop()
			if err != nil {
				return err
			}

			elems[n] = top
		}

		return f.stack.Push(&object.List{
			Value: elems,
		})
	}

	Effectors[bytecode.MakeMap] = func(v *VM, f *Frame, arg rune) error {
		keys := make(map[string]object.Object, arg)
		vals := make(map[string]object.Object, arg)

		for n := 0; n < int(arg); n++ {
			val, err := f.stack.Pop()
			if err != nil {
				return err
			}

			key, err := f.stack.Pop()
			if err != nil {
				return err
			}

			hash, err := structhash.Hash(key, 1)
			if err != nil {
				return err
			}

			keys[hash] = key
			vals[hash] = val
		}

		return f.stack.Push(&object.Map{
			Keys:   keys,
			Values: vals,
		})
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
		}

		return f.stack.Push(&object.Boolean{Value: false})
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
			return makeError(TypeError, "could not apply prefix operator %s to %s", op, obj.String())
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
			return makeError(TypeError, "could not apply infix operator %s between %s and %s", op, left.String(), right.String())
		}

		return f.stack.Push(result)
	}
}

func callFunction(v *VM, f *Frame, fn *object.Function, argCount rune) error {
	if int(argCount) != len(fn.Parameters) {
		return makeError(ArgumentError, "wrong amount of arguments passed to a function. expected %d, got %d", len(fn.Parameters), argCount)
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

func indexCollection(v *VM, f *Frame, items []object.Object, argCount rune) error {
	if argCount != 1 {
		return makeError(ArgumentError, "a list can only be called with one argument")
	}

	arg, err := f.stack.Pop()
	if err != nil {
		return err
	}

	var indexObj object.Object

	if argItems, ok := arg.Items(); ok && len(argItems) == 1 {
		indexObj = argItems[0]
	} else {
		indexObj = arg
	}

	idx, ok := indexObj.(*object.Number)
	if !ok {
		return makeError(ArgumentError, "a list can only be called with a number or a length-1 list/tuple containing a number")
	}

	index := int(math.Floor(idx.Value))

	if index < 0 || index >= len(items) {
		return makeError(IndexError, "%d is out of bounds", index)
	}

	return f.stack.Push(items[index])
}
