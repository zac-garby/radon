package compiler

import (
	"github.com/Zac-Garby/lang/ast"
	"github.com/Zac-Garby/lang/object"
)

// A Compiler translates an abstract syntax tree
// into some bytecode
type Compiler struct {
	// The generated code
	Bytes []byte

	Constants []object.Object
	Names     []string
}

// New instantiates a new Compiler
func New() *Compiler {
	return &Compiler{
		Bytes:     make([]byte, 0, 128),
		Constants: make([]object.Object, 0, 8),
		Names:     make([]string, 0, 8),
	}
}

// Compile compiles an entire ast.Program
func (c *Compiler) Compile(p ast.Program) error {
	p, err := PreprocessProgram(p)
	if err != nil {
		return err
	}

	for _, stmt := range p.Statements {
		if err := c.CompileStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}
