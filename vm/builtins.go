package vm

import (
	"fmt"

	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/object"
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
