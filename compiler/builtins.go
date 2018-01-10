package compiler

import (
	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/bytecode"
)

type builtinFn struct {
	// The name of the function. Used to call it
	name string

	// The amount of parameters the function accepts
	parameters int

	// Whether the compiler should automatically push
	// the arguments
	autoPush bool

	// Compile compiles the builtin function. Compiler
	// will be called after all the parameters have been
	// added to the bytecode
	compile func(c *Compiler, args []ast.Expression) error
}

var builtinFns = []*builtinFn{
	normalBuiltin("print", bytecode.Println),
	normalBuiltin("echo", bytecode.Print),
	normalBuiltin("len", bytecode.Length),
}

func normalBuiltin(name string, op byte) *builtinFn {
	return &builtinFn{
		name:       name,
		parameters: 1,
		autoPush:   true,

		compile: func(c *Compiler, args []ast.Expression) error {
			c.push(op)
			return nil
		},
	}
}

func getBuiltin(name string) (*builtinFn, bool) {
	for _, fn := range builtinFns {
		if fn.name == name {
			return fn, true
		}
	}

	return nil, false
}
