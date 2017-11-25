package compiler

import (
	"fmt"
	"reflect"

	"github.com/Zac-Garby/lang/ast"
)

// CompileExpression compiles an AST expression node
func (c *Compiler) CompileExpression(e ast.Expression) error {
	switch e.(type) {
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(e))
	}
}
