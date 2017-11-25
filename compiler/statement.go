package compiler

import (
	"fmt"
	"reflect"

	"github.com/Zac-Garby/lang/ast"
)

// CompileStatement compiles a single ast.Statement
func (c *Compiler) CompileStatement(s ast.Statement) error {
	switch stmt := s.(type) {
	case *ast.ExpressionStatement:
		return c.CompileExpression(stmt.Expr)
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(s))
	}
}
