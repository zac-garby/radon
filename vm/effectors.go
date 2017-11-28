package vm

import (
	"math"

	"github.com/cnf/structhash"

	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/object"
)

type effector func(f *Frame, i bytecode.Instruction)

var effectors map[byte]effector

func init() {
	effectors = map[byte]effector{
		bytecode.Pop: bytePop,
		bytecode.Dup: byteDup,

		bytecode.LoadConst:  byteLoadConst,
		bytecode.LoadName:   byteLoadName,
		bytecode.StoreName:  byteStoreName,
		bytecode.LoadField:  byteLoadField,
		bytecode.StoreField: byteStoreField,

		bytecode.UnaryInvert: bytePrefix,
		bytecode.UnaryNegate: bytePrefix,

		bytecode.BinaryAdd:      byteInfix,
		bytecode.BinarySubtract: byteInfix,
		bytecode.BinaryMultiply: byteInfix,
		bytecode.BinaryDivide:   byteInfix,
		bytecode.BinaryExponent: byteInfix,
		bytecode.BinaryFloorDiv: byteInfix,
		bytecode.BinaryMod:      byteInfix,
		bytecode.BinaryBitOr:    byteInfix,
		bytecode.BinaryBitAnd:   byteInfix,
		bytecode.BinaryEquals:   byteEquals,
		bytecode.BinaryNotEqual: byteNotEqual,
		bytecode.BinaryLessThan: bincmp,
		bytecode.BinaryMoreThan: bincmp,
		bytecode.BinaryLessEq:   bincmp,
		bytecode.BinaryMoreEq:   bincmp,

		bytecode.CallFn: byteCall,
		bytecode.Return: byteReturn,

		bytecode.Print:   bytePrint,
		bytecode.Println: bytePrintln,
		bytecode.Length:  byteLength,

		bytecode.Jump:        byteJump,
		bytecode.JumpIfTrue:  byteJumpIfTrue,
		bytecode.JumpIfFalse: byteJumpIfFalse,
		bytecode.Break:       byteBreak,
		bytecode.Next:        byteNext,
		bytecode.LoopStart:   byteLoopStart,
		bytecode.LoopEnd:     byteLoopEnd,

		bytecode.MakeList:  byteMakeList,
		bytecode.MakeTuple: byteMakeTuple,
		bytecode.MakeMap:   byteMakeMap,
	}
}

func bytePop(f *Frame, i bytecode.Instruction) {
	f.stack.pop()
}

func byteDup(f *Frame, i bytecode.Instruction) {
	f.stack.dup()
}

func byteLoadConst(f *Frame, i bytecode.Instruction) {
	f.stack.push(f.constants[i.Arg])
}

func byteLoadName(f *Frame, i bytecode.Instruction) {
	name, ok := f.getName(i.Arg)
	if !ok {
		f.vm.err = Err("name not defined when loading a name", ErrInternal)
		return
	}

	val, ok := f.searchName(name)
	if !ok {
		f.vm.err = Errf("name %s not found in the current scope", ErrNotFound, name)
		return
	}

	f.stack.push(val)
}

func byteStoreName(f *Frame, i bytecode.Instruction) {
	name, ok := f.getName(i.Arg)
	if !ok {
		f.vm.err = Err("name not defined when storing a name", ErrInternal)
		return
	}

	f.store.Set(name, f.stack.top(), true)
}

func byteLoadField(f *Frame, i bytecode.Instruction) {
	field, obj := f.stack.pop(), f.stack.pop()

	var val object.Object

	if col, ok := obj.(object.Collection); ok {
		if index, ok := field.(*object.Number); ok {
			idx := int(index.Value)

			val = col.GetIndex(idx)
		} else {
			f.vm.err = Errf("non-numeric type %s used to index a collection", ErrWrongType, field.Type())
			return
		}
	} else if cont, ok := obj.(object.Container); ok {
		val = cont.GetKey(field)
	} else {
		f.vm.err = Errf("cannot index type %s", ErrNotFound, obj.Type())
	}

	f.stack.push(val)
}

func byteStoreField(f *Frame, i bytecode.Instruction) {
	field, obj, val := f.stack.pop(), f.stack.pop(), f.stack.top()

	if col, ok := obj.(object.Collection); ok {
		if index, ok := field.(*object.Number); ok {
			idx := int(index.Value)

			col.SetIndex(idx, val)
		} else {
			f.vm.err = Errf("non-numeric type %s used to index a collection", ErrWrongType, field.Type())
			return
		}
	} else if cont, ok := obj.(object.Container); ok {
		cont.SetKey(field, val)
	} else {
		f.vm.err = Errf("cannot index type %s", ErrWrongType, obj.Type())
	}
}

func bytePrefix(f *Frame, i bytecode.Instruction) {
	right := f.stack.pop()

	if i.Code == bytecode.UnaryInvert {
		f.stack.push(object.MakeObj(!object.IsTruthy(right)))
		return
	}

	if n, ok := right.(*object.Number); ok {
		val := n.Value
		f.stack.push(numPrefix(i.Code, val))
	} else {
		f.vm.err = Err("prefix r-value of invalid type", ErrWrongType)
	}
}

func numPrefix(opcode byte, val float64) object.Object {
	switch opcode {
	case bytecode.UnaryNegate:
		val = -val
	}

	return &object.Number{Value: val}
}

func byteInfix(f *Frame, i bytecode.Instruction) {
	right, left := f.stack.pop(), f.stack.pop()

	if n, ok := left.(*object.Number); ok {
		if m, ok := right.(*object.Number); ok {
			f.stack.push(numInfix(f, i.Code, n.Value, m.Value))
		} else if m, ok := right.(object.Collection); ok {
			f.stack.push(numColInfix(f, i.Code, n.Value, m))
		} else {
			f.vm.err = Err("infix r-value of invalid type when l-value is <number>", ErrWrongType)
			return
		}
	} else if n, ok := left.(object.Collection); ok {
		if m, ok := right.(*object.Number); ok {
			f.stack.push(numColInfix(f, i.Code, m.Value, n))
		} else if m, ok := right.(object.Collection); ok {
			f.stack.push(colInfix(f, i.Code, n, m))
		} else {
			f.vm.err = Err("infix r-value of invalid type when l-value is a collection", ErrWrongType)
		}
	} else {
		f.vm.err = Err("infix l-value of invalid type", ErrWrongType)
		return
	}
}

func numInfix(f *Frame, opcode byte, left, right float64) object.Object {
	var val float64

	switch opcode {
	case bytecode.BinaryAdd:
		val = left + right
	case bytecode.BinarySubtract:
		val = left - right
	case bytecode.BinaryMultiply:
		val = left * right
	case bytecode.BinaryDivide:
		val = left / right
	case bytecode.BinaryExponent:
		val = math.Pow(left, right)
	case bytecode.BinaryFloorDiv:
		val = math.Floor(left / right)
	case bytecode.BinaryMod:
		val = math.Mod(left, right)
	case bytecode.BinaryBitOr:
		val = float64(int64(left) | int64(right))
	case bytecode.BinaryBitAnd:
		val = float64(int64(left) & int64(right))
	default:
		op := bytecode.Instructions[opcode].Name[7:]
		f.vm.err = Errf("operator %s not supported for two numbers", ErrNoOp, op)
	}

	return &object.Number{Value: val}
}

func numColInfix(f *Frame, opcode byte, left float64, right object.Collection) object.Object {
	var (
		result   []object.Object
		elements = right.Elements()
	)

	if opcode == bytecode.BinaryMultiply {
		for i := 0; i < int(left); i++ {
			result = append(result, elements...)
		}
	} else {
		op := bytecode.Instructions[opcode].Name[7:]
		f.vm.err = Errf("operator %s not supported for a collection and a number", ErrNoOp, op)
	}

	return object.MakeObj(result)
}

func colInfix(f *Frame, opcode byte, left, right object.Collection) object.Object {
	var (
		lefts  = left.Elements()
		rights = right.Elements()
		elems  []object.Object
	)

	switch opcode {
	case bytecode.BinaryAdd:
		elems = append(lefts, rights...)
	case bytecode.BinarySubtract:
		for _, el := range lefts {
			for _, rel := range rights {
				if el.Equals(rel) {
					goto next
				}
			}

			elems = append(elems, el)
		next:
		}
	case bytecode.BinaryBitOr:
		for _, el := range append(lefts, rights...) {
			unique := true

			for _, rel := range elems {
				if el.Equals(rel) {
					unique = false
					break
				}
			}

			if unique {
				elems = append(elems, el)
			}
		}
	case bytecode.BinaryBitAnd:
		for _, el := range lefts {
			both := false

			for _, rel := range rights {
				if el.Equals(rel) {
					both = true
					break
				}
			}

			if both {
				elems = append(elems, el)
			}
		}
	default:
		op := bytecode.Instructions[opcode].Name[7:]
		f.vm.err = Errf("operator %s not supported for two collections", ErrNoOp, op)
	}

	return object.MakeObj(elems)
}

func bincmp(f *Frame, i bytecode.Instruction) {
	f.offsetToInstructionIndex(int(i.Arg))

	b, a := f.stack.pop(), f.stack.pop()

	n, ok := a.(*object.Number)
	if !ok {
		f.vm.err = Err("non-numeric value in numeric binary expression", ErrWrongType)
		return
	}

	m, ok := b.(*object.Number)
	if !ok {
		f.vm.err = Err("non-numeric value in numeric binary expression", ErrWrongType)
		return
	}

	lval := n.Value
	rval := m.Value

	var result bool

	switch i.Code {
	case bytecode.BinaryLessThan:
		result = lval < rval
	case bytecode.BinaryMoreThan:
		result = lval > rval
	case bytecode.BinaryLessEq:
		result = lval <= rval
	case bytecode.BinaryMoreEq:
		result = lval >= rval
	}

	f.stack.push(&object.Boolean{Value: result})
}

func byteEquals(f *Frame, i bytecode.Instruction) {
	right, left := f.stack.pop(), f.stack.pop()
	eq := left.Equals(right)

	f.stack.push(object.MakeObj(eq))
}

func byteNotEqual(f *Frame, i bytecode.Instruction) {
	right, left := f.stack.pop(), f.stack.pop()
	eq := left.Equals(right)

	f.stack.push(object.MakeObj(!eq))
}

func byteCall(f *Frame, i bytecode.Instruction) {
	argCount := f.stack.pop().(*object.Number).Value

	fn, ok := f.stack.pop().(*object.Function)
	if !ok {
		f.vm.err = Errf("cannot call non-function type: %s", ErrWrongType, fn.Type())
		return
	}

	if argCount != float64(len(fn.Parameters)) {
		f.vm.err = Errf("wrong amount of arguments supplied to the function. expected %v", ErrArgument, len(fn.Parameters))
		return
	}

	locals := f.store
	locals.Names = fn.Names

	for _, param := range fn.Parameters {
		locals.Set(param, f.stack.pop(), true)
	}

	data := map[string]object.Object{}

	for key, item := range locals.Data {
		data[key] = item.Value
	}

	result, err := fn.OnCall(fn, data)
	if err != nil {
		f.vm.err = err
		return
	}

	if result != nil {
		f.stack.push(result)
		return
	}

	// Create the function's frame
	fnFrame := &Frame{
		code:      fn.Code,
		constants: fn.Constants,
		store:     locals,
		offset:    0,
		prev:      f,
		stack:     newStack(),
		vm:        f.vm,
	}

	// Push and execute the function's frame
	f.vm.runFrame(fnFrame)

	if len(fnFrame.stack.objects) > 0 {
		ret := fnFrame.stack.pop()

		// Push the returned value
		f.stack.push(ret)
	}
}

func byteReturn(f *Frame, i bytecode.Instruction) {
	f.offset = len(f.code) - 1
}

func byteJump(f *Frame, i bytecode.Instruction) {
	f.offset = f.offsetToInstructionIndex(int(i.Arg))
}

func byteJumpIfTrue(f *Frame, i bytecode.Instruction) {
	obj := f.stack.pop()

	if object.IsTruthy(obj) {
		f.offset = f.offsetToInstructionIndex(int(i.Arg))
	}
}

func byteJumpIfFalse(f *Frame, i bytecode.Instruction) {
	obj := f.stack.pop()

	if !object.IsTruthy(obj) {
		f.offset = f.offsetToInstructionIndex(int(i.Arg))
	}
}

func byteBreak(f *Frame, i bytecode.Instruction) {
	if len(f.breaks) < 1 {
		f.vm.err = Err("break statement found outside loop", ErrSyntax)
		return
	}

	top := f.breaks[len(f.breaks)-1]
	f.offset = top
}

func byteNext(f *Frame, i bytecode.Instruction) {
	if len(f.nexts) < 1 {
		f.vm.err = Err("next statement found outside loop", ErrSyntax)
		return
	}

	top := f.nexts[len(f.nexts)-1]
	f.offset = top
}

func byteLoopStart(f *Frame, i bytecode.Instruction) {
	f.nexts = append(f.nexts, f.offset+1)

	var o int

	for o = f.offset; f.code[o].Code != bytecode.LoopEnd; o++ {
		// Nothing here
	}

	f.breaks = append(f.breaks, o)
}

func byteLoopEnd(f *Frame, i bytecode.Instruction) {
	f.breaks = f.breaks[:len(f.breaks)-1]
	f.nexts = f.nexts[:len(f.nexts)-1]
}

func byteMakeList(f *Frame, i bytecode.Instruction) {
	elems := make([]object.Object, i.Arg)

	for n := int(i.Arg) - 1; n >= 0; n-- {
		elems[n] = f.stack.pop()
	}

	f.stack.push(&object.List{
		Value: elems,
	})
}

func byteMakeTuple(f *Frame, i bytecode.Instruction) {
	elems := make([]object.Object, i.Arg)

	for n := int(i.Arg) - 1; n >= 0; n-- {
		elems[n] = f.stack.pop()
	}

	f.stack.push(&object.Tuple{
		Value: elems,
	})
}

func byteMakeMap(f *Frame, i bytecode.Instruction) {
	keys := make(map[string]object.Object, i.Arg)
	values := make(map[string]object.Object, i.Arg)

	for n := 0; n < int(i.Arg); n++ {
		val, key := f.stack.pop(), f.stack.pop()

		hash, err := structhash.Hash(key, 1)
		if err != nil {
			f.vm.err = err
		}

		keys[hash] = key
		values[hash] = val
	}

	obj := &object.Map{
		Keys:   keys,
		Values: values,
	}

	f.stack.push(obj)
}
