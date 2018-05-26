package compiler

import (
	"fmt"
	"reflect"

	"github.com/Zac-Garby/radon/ast"
)

// CompileStatement takes an AST statement and generates some bytecode for it.
func (c *Compiler) CompileStatement(s ast.Statement) error {
	switch node := s.(type) {
	case *ast.ExpressionStatement:
		return c.CompileExpression(node.Expr)
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(s))
	}
}
