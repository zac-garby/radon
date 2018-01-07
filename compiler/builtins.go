package compiler

import (
	"errors"

	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/object"
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
	&builtinFn{
		name:       "print",
		parameters: 1,
		autoPush:   true,

		compile: func(c *Compiler, args []ast.Expression) error {
			c.push(bytecode.Println)
			index, err := c.addConst(object.EmptyObj)
			if err != nil {
				return err
			}

			c.loadConst(index)
			return nil
		},
	},

	&builtinFn{
		name:       "echo",
		parameters: 1,
		autoPush:   true,

		compile: func(c *Compiler, args []ast.Expression) error {
			c.push(bytecode.Print)
			index, err := c.addConst(object.EmptyObj)
			if err != nil {
				return err
			}

			c.loadConst(index)
			return nil
		},
	},

	&builtinFn{
		name:       "len",
		parameters: 1,
		autoPush:   true,

		compile: func(c *Compiler, args []ast.Expression) error {
			c.push(bytecode.Length)
			return nil
		},
	},

	&builtinFn{
		name:       "pop",
		parameters: 0,

		compile: func(c *Compiler, args []ast.Expression) error {
			c.push(bytecode.PushTop)
			return nil
		},
	},

	&builtinFn{
		name:       "list",
		parameters: 1,

		compile: func(c *Compiler, args []ast.Expression) error {
			num, ok := args[0].(*ast.Number)
			if !ok {
				return errors.New("wrong_type: list(..) accepts a single numerical argument")
			}

			var (
				length    = num.Value
				low, high = runeToBytes(rune(length))
			)

			c.push(bytecode.MakeList, high, low)

			return nil
		},
	},
}

func getBuiltin(name string) (*builtinFn, bool) {
	for _, fn := range builtinFns {
		if fn.name == name {
			return fn, true
		}
	}

	return nil, false
}
