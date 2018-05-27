package compiler

import (
	"fmt"
	"reflect"

	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/bytecode"
)

// CompileStatement takes an AST statement and generates some bytecode for it.
func (c *Compiler) CompileStatement(s ast.Statement) error {
	switch node := s.(type) {
	case *ast.ExpressionStatement:
		return c.CompileExpression(node.Expr)
	case *ast.Return:
		return c.compileReturn(node)
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(s))
	}
}

func (c *Compiler) compileReturn(node *ast.Return) error {
	if err := c.CompileExpression(node.Value); err != nil {
		return err
	}

	c.push(bytecode.Return)

	return nil
}
