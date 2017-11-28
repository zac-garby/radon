package compiler

import (
	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/object"
)

type builtinFn struct {
	// The name of the function. Used to call it
	name string

	// The amount of parameters the function accepts
	parameters int

	// Compile compiles the builtin function. Compiler
	// will be called after all the parameters have been
	// added to the bytecode
	compile func(c *Compiler) error
}

var builtinFns = []*builtinFn{
	&builtinFn{
		name:       "print",
		parameters: 1,

		compile: func(c *Compiler) error {
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

		compile: func(c *Compiler) error {
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

		compile: func(c *Compiler) error {
			c.push(bytecode.Length)
			return nil
		},
	},

	&builtinFn{
		name:       "pop",
		parameters: 0,

		compile: func(c *Compiler) error {
			c.push(bytecode.PushTop)
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
