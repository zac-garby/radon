package bytecode

import (
	"errors"
)

// Raw is the raw bytecode, i.e. a list of bytes.
// If a byte has an argument, the next two bytes
// are converted into a single 16-bit integer:
//
//		(a << 8) + b
//
// Where a is the first argument byte, and b is the
// second.
type Raw []byte

// Code is the "parsed" bytecode, i.e. a list of
// instructions, with their arguments.
type Code []Instruction

// Instruction is a parsed bytecode
// instruction.
type Instruction struct {
	Code byte
	Arg  rune
	Name string
}

// ErrOutOfBytes is thrown by Read when a byte
// which takes arguments isn't followed by at
// least two more bytes.
var ErrOutOfBytes = errors.New("bytecode: not enough bytes remaining")

// Read takes some raw bytecode and outputs
// the "parsed" bytecode as a Code struct.
//
// If there is an error, it is ErrOutOfBytes,
// signifying there aren't enough bytes left
// after an instruction with arity > 0.
func Read(raw Raw) (Code, error) {
	var (
		code  Code
		index int
	)

	for index < len(raw) {
		var (
			cur  = raw[index]
			data = Instructions[cur]

			instr = Instruction{
				Code: cur,
				Name: data.Name,
			}
		)

		if data.HasArg {
			if index+2 >= len(raw) {
				return code, ErrOutOfBytes
			}

			var (
				a   = raw[index+1]
				b   = raw[index+2]
				arg = (rune(a) << 8) + rune(b)
			)

			index += 2

			instr.Arg = arg
		}

		code = append(code, instr)

		index++
	}

	return code, nil
}
