package bytecode

import (
	"fmt"
	"io"
)

// An Instruction is a parsed bytecode instruction.
type Instruction struct {
	Code byte
	Name string
	Arg  rune
}

// Code is a bytecode program. Usually obtained by parsing a
// series of bytes.
type Code []Instruction

// Read takes an io.Reader and parses it into a series of instructions
// in the form of a Code instance. This allows it to be executed more
// efficiently by the virtual machine.
//
// An error can occur if an argument is expected but EOF is found, or if
// the reader encounters any other error.
func Read(r io.Reader) (Code, error) {
	var code Code

	for {
		var (
			p      = make([]byte, 1)
			_, err = r.Read(p)
		)
		if err == io.EOF {
			break
		}
		if err != nil {
			return code, err
		}

		var (
			bc    = p[0]
			data  = Instructions[bc]
			instr = Instruction{
				Code: bc,
				Name: data.Name,
			}
		)

		if data.HasArg {
			var (
				argBytes = make([]byte, 2)
				_, err   = r.Read(argBytes)
			)
			if err == io.EOF {
				return code, fmt.Errorf("read: not enough arguments to %s", data.Name)
			}
			if err != nil {
				return code, err
			}

			// Combine the next two bytes into a single rune.
			// e.g. if the next two bytes are 01100010 and 10010011, the
			// argument will be 0110001010010011
			instr.Arg = (rune(argBytes[0]) << 8) + rune(argBytes[1])
		}

		code = append(code, instr)
	}

	return code, nil
}
