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
	case *ast.String:
		return c.compileString(node)
	case *ast.Boolean:
		return c.compileBoolean(node)
	case *ast.Nil:
		return c.compileNil(node)
	case *ast.Identifier:
		return c.compileIdentifier(node)
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

func (c *Compiler) compileBoolean(node *ast.Boolean) error {
	_, err := c.addAndLoad(&object.Boolean{Value: node.Value})
	return err
}

func (c *Compiler) compileNil(node *ast.Nil) error {
	_, err := c.addAndLoad(&object.Nil{})
	return err
}

func (c *Compiler) compileIdentifier(node *ast.Identifier) error {
	return c.compileName(node.Value)
}
