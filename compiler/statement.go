package compiler

import (
	"fmt"
	"reflect"

	"github.com/Zac-Garby/lang/bytecode"

	"github.com/Zac-Garby/lang/ast"
)

// CompileStatement compiles a single ast.Statement
func (c *Compiler) CompileStatement(s ast.Statement) error {
	switch node := s.(type) {
	case *ast.ExpressionStatement:
		return c.CompileExpression(node.Expr)

	case *ast.ReturnStatement:
		return c.compileReturnStatement(node)

	case *ast.NextStatement:
		return c.compileNext(node)

	case *ast.BreakStatement:
		return c.compileBreak(node)

	case *ast.Loop:
		return c.compileLoop(node)

	case *ast.WhileLoop:
		return c.compileWhile(node)

	case *ast.ForLoop:
		return c.compileFor(node)

	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(s))
	}
}

func (c *Compiler) compileReturnStatement(node *ast.ReturnStatement) error {
	if node.Value != nil {
		if err := c.CompileExpression(node.Value); err != nil {
			return err
		}
	}

	c.push(bytecode.Return)

	return nil
}

func (c *Compiler) compileNext(node *ast.NextStatement) error {
	c.push(bytecode.Next)

	return nil
}

func (c *Compiler) compileBreak(node *ast.BreakStatement) error {
	c.push(bytecode.Break)

	return nil
}

func (c *Compiler) compileLoop(node *ast.Loop) error {
	c.push(bytecode.LoopStart)

	// Jump here to go to the next iteration
	start := len(c.Bytes) - 1

	// Compile the loop's body
	if err := c.CompileExpression(node.Body); err != nil {
		return err
	}

	// After the body, jump back to the beginning of the loop
	low, high := runeToBytes(rune(start))
	c.push(bytecode.Jump, high, low)

	c.push(bytecode.LoopEnd)

	return nil
}

func (c *Compiler) compileWhile(node *ast.WhileLoop) error {
	c.push(bytecode.LoopStart)

	// Jump here to go to the next iteration
	start := len(c.Bytes) - 1

	if err := c.CompileExpression(node.Condition); err != nil {
		return err
	}

	// An empty jump to the end of the loop
	c.push(bytecode.JumpIfFalse, 0, 0)
	skipJump := len(c.Bytes) - 3

	// Compile the loop's body
	if err := c.CompileExpression(node.Body); err != nil {
		return err
	}

	// After the body, jump back to the beginning of the loop
	low, high := runeToBytes(rune(start))
	c.push(bytecode.Jump, high, low)

	// If the condition isn't met, jump to the end of the loop
	skipIndex := rune(len(c.Bytes))
	low, high = runeToBytes(skipIndex)
	c.Bytes[skipJump+1] = high
	c.Bytes[skipJump+2] = low

	c.push(bytecode.LoopEnd)

	return nil
}

func (c *Compiler) compileFor(node *ast.ForLoop) error {
	if err := c.CompileExpression(node.Init); err != nil {
		return err
	}

	body := &ast.Block{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Expr: node.Body,
			},
			&ast.ExpressionStatement{
				Expr: node.Increment,
			},
		},
	}

	while := &ast.WhileLoop{
		Condition: node.Condition,
		Body:      body,
	}

	return c.CompileStatement(while)
}
