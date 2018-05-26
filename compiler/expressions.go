package compiler

import (
	"fmt"
	"reflect"

	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/object"
)

// CompileExpression takes an AST expression and generates some bytecode
// for it.
func (c *Compiler) CompileExpression(e ast.Expression) error {
	switch node := e.(type) {
	case *ast.Number:
		return c.compileNumber(node)
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(e))
	}
}

func (c *Compiler) compileNumber(node *ast.Number) error {
	_, err := c.addAndLoad(&object.Number{Value: node.Value})
	return err
}

func (c *Compiler) compileString(node *ast.String) error {
	_, err := c.addAndLoad(&object.String{Value: node.Value})
	return err
}
