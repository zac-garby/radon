package vm

import (
	"fmt"

	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/object"
)

func bytePrint(f *Frame, i bytecode.Instruction) {
	top, err := f.stack.pop()
	if err != nil {
		f.vm.err = err
		return
	}

	fmt.Print(top)
}

func bytePrintln(f *Frame, i bytecode.Instruction) {
	top, err := f.stack.pop()
	if err != nil {
		f.vm.err = err
		return
	}

	fmt.Println(top)
}

func byteLength(f *Frame, i bytecode.Instruction) {
	top, err := f.stack.pop()
	if err != nil {
		f.vm.err = err
		return
	}

	if col, ok := top.(object.Collection); ok {
		f.stack.push(&object.Number{
			Value: float64(len(col.Elements())),
		})
	} else {
		f.vm.err = Errf("cannot get the length of type %s", ErrWrongType, top.Type())
	}
}
