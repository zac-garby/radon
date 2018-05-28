package compiler

import (
	"errors"
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
	case *ast.Next:
		return c.compileNext(node)
	case *ast.Break:
		return c.compileBreak(node)
	case *ast.While:
		return c.compileWhile(node)
	case *ast.For:
		return c.compileFor(node)
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

func (c *Compiler) compileNext(node *ast.Next) error {
	c.push(bytecode.Next)
	return nil
}

func (c *Compiler) compileBreak(node *ast.Break) error {
	c.push(bytecode.Break)
	return nil
}

func (c *Compiler) compileWhile(node *ast.While) error {
	c.push(bytecode.StartLoop)

	// Jump here for the next iteration
	start := len(c.Bytes) - 1

	if err := c.CompileExpression(node.Condition); err != nil {
		return err
	}

	// An empty jump to the end of the loop
	c.push(bytecode.JumpUnless, 0, 0)
	skipJump := len(c.Bytes) - 3

	// Compile the loop's body
	if err := c.CompileExpression(node.Body); err != nil {
		return err
	}

	// After the body, jump back to the beginning
	index, err := c.addJump(start)
	if err != nil {
		return err
	}
	low, high := runeToBytes(index)
	c.push(bytecode.Jump, high, low)

	// If the condition ism't met, jump to the end of the loop
	c.setJumpArg(skipJump, len(c.Bytes))

	c.push(bytecode.EndLoop)

	return nil
}

func (c *Compiler) compileFor(node *ast.For) error {
	if err := c.CompileExpression(node.Collection); err != nil {
		return err
	}

	id, ok := node.Var.(*ast.Identifier)
	if !ok {
		return errors.New("compiler: a for-loop counter must be an identifier")
	}

	index, err := c.addName(id.Value)
	if err != nil {
		return err
	}

	low, high := runeToBytes(rune(index))
	c.push(bytecode.StartFor, high, low)

	if err := c.CompileExpression(node.Body); err != nil {
		return err
	}

	c.push(bytecode.EndFor)

	return nil
}
