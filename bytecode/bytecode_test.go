package bytecode_test

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/Zac-Garby/radon/bytecode"
)

func TestRead(t *testing.T) {
	// Bytes to parse, hard to read because I have to write in decimal
	b := []byte{
		2, 1, 56, // LOAD_CONST 312
		2, 0, 5, // LOAD_CONST 5
		10, // BINARY_ADD
	}

	var (
		buf       = bytes.NewBuffer(b)
		code, err = Read(buf)
	)

	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	exp := Code{
		Instruction{
			Code: 2,
			Name: "LOAD_CONST",
			Arg:  312,
		},
		Instruction{
			Code: 2,
			Name: "LOAD_CONST",
			Arg:  5,
		},
		Instruction{
			Code: 10,
			Name: "BINARY_ADD",
		},
	}

	if len(code) != len(exp) {
		fmt.Printf("got %d instructions, expected %d\n", len(code), len(exp))
		t.FailNow()
	}

	for i, instr := range code {
		e := exp[i]

		if instr.Code != e.Code || instr.Name != e.Name || instr.Arg != e.Arg {
			fmt.Printf("%v doesn't equal %v\n", instr, e)
			t.Fail()
		}
	}
}
