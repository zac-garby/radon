package vm

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/object"
)

func bytePrint(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()
	fmt.Fprint(f.vm.Out, top)
}

func bytePrintln(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()
	fmt.Fprintln(f.vm.Out, top)
}

func byteLength(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	if col, ok := top.(object.Collection); ok {
		f.stack.push(&object.Number{
			Value: float64(len(col.Elements())),
		})
	} else {
		f.vm.err = Errf("cannot get the length of type %s", ErrWrongType, top.Type())
	}
}

func byteTypeof(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	f.stack.push(&object.String{
		Value: fmt.Sprintf("%s", top.Type()),
	})
}

func byteModelof(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	if m, ok := top.(*object.Map); ok {
		if m.Model == nil {
			f.stack.push(object.NilObj)
		} else {
			f.stack.push(m.Model)
		}
	} else {
		f.vm.err = Err("can only call modelof() on maps", ErrWrongType)
	}
}

func byteToStr(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	f.stack.push(&object.String{
		Value: top.String(),
	})
}

func byteToNum(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	switch o := top.(type) {
	case *object.Number:
		f.stack.push(top)

	case *object.String:
		n, err := strconv.ParseFloat(o.Value, 64)
		if err != nil {
			f.vm.err = Errf("cannot convert %s to a number", ErrArgument, top.Debug())
		} else {
			f.stack.push(&object.Number{
				Value: n,
			})
		}

	case *object.Boolean:
		if o.Value == true {
			f.stack.push(&object.Number{
				Value: 1,
			})
		} else {
			f.stack.push(&object.Number{
				Value: 0,
			})
		}

	default:
		f.vm.err = Errf("cannot convert type %s to a number", ErrWrongType, top.Type())
	}
}

func byteToList(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	col, ok := top.(object.Collection)
	if !ok {
		f.vm.err = Errf("can only convert collection types to lists. got %s", ErrWrongType, top.Type())
		return
	}

	list := col.Elements()
	f.stack.push(&object.List{
		Value: list,
	})
}

func byteToTuple(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	col, ok := top.(object.Collection)
	if !ok {
		f.vm.err = Errf("can only convert collection types to tuples. got %s", ErrWrongType, top.Type())
		return
	}

	list := col.Elements()
	f.stack.push(&object.Tuple{
		Value: list,
	})
}

func byteRound(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	num, ok := top.(*object.Number)
	if !ok {
		f.vm.err = Errf("cannot round non-numeric type %s", ErrWrongType, top.Type())
		return
	}

	x := num.Value

	// https://gist.github.com/gdm85/44f648cc97bb3bf847f21c87e9d19b2d
	const (
		shift    = 64 - 11 - 1
		bias     = 1023
		mask     = 0x7FF
		signMask = 1 << 63
		fracMask = 1<<shift - 1
		half     = 1 << (shift - 1)
		one      = bias << shift
	)

	bits := math.Float64bits(x)
	e := uint(bits>>shift) & mask
	if e < bias {
		// Round abs(x) < 1 including denormals.
		bits &= signMask // +-0
		if e == bias-1 {
			bits |= one // +-1
		}
	} else if e < bias+shift {
		// Round any abs(x) >= 1 containing a fractional component [0,1).
		//
		// Numbers with larger exponents are returned unchanged since they
		// must be either an integer, infinity, or NaN.
		e -= bias
		bits += half >> e
		bits &^= fracMask >> e
	}

	f.stack.push(&object.Number{
		Value: math.Float64frombits(bits),
	})
}

func byteFloor(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	num, ok := top.(*object.Number)
	if !ok {
		f.vm.err = Errf("cannot floor non-numeric type %s", ErrWrongType, top.Type())
		return
	}

	f.stack.push(&object.Number{
		Value: math.Floor(num.Value),
	})
}

func byteCeil(f *Frame, i bytecode.Instruction) {
	top, _ := f.stack.pop()

	num, ok := top.(*object.Number)
	if !ok {
		f.vm.err = Errf("cannot ceil non-numeric type %s", ErrWrongType, top.Type())
		return
	}

	f.stack.push(&object.Number{
		Value: math.Ceil(num.Value),
	})
}
