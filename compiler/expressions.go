package compiler

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/bytecode"
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
	case *ast.Infix:
		return c.compileInfix(node)
	case *ast.Prefix:
		return c.compilePrefix(node)
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

func (c *Compiler) compileInfix(node *ast.Infix) error {
	left, right := node.Left, node.Right

	// Some operators are handled differently
	switch node.Operator {
	case "=":
		return c.compileAssignOrDeclare(left, right, "assign")
	case ":=":
		return c.compileAssignOrDeclare(left, right, "declare")
	case ".":
		return c.compileDot(left, right)
	}

	if err := c.CompileExpression(left); err != nil {
		return err
	}

	if err := c.CompileExpression(right); err != nil {
		return err
	}

	op, ok := map[string]byte{
		"+":  bytecode.BinaryAdd,
		"-":  bytecode.BinarySub,
		"*":  bytecode.BinaryMul,
		"/":  bytecode.BinaryDiv,
		"^":  bytecode.BinaryExp,
		"//": bytecode.BinaryFloorDiv,
		"%":  bytecode.BinaryMod,
		"||": bytecode.BinaryLogicOr,
		"&&": bytecode.BinaryLogicAnd,
		"|":  bytecode.BinaryBitOr,
		"&":  bytecode.BinaryBitAnd,
		"==": bytecode.BinaryEqual,
		"!=": bytecode.BinaryNotEqual,
		"<":  bytecode.BinaryLess,
		">":  bytecode.BinaryMod,
		"<=": bytecode.BinaryLessEq,
		">=": bytecode.BinaryMoreEq,
		",":  bytecode.BinaryTuple,
	}[node.Operator]

	if !ok {
		return fmt.Errorf("compiler: operator %s not yet implemented", node.Operator)
	}

	c.push(op)

	return nil
}

func (c *Compiler) compileDot(left, right ast.Expression) error {
	if err := c.CompileExpression(left); err != nil {
		return err
	}

	if id, ok := right.(*ast.Identifier); ok {
		index, err := c.addConst(&object.String{Value: id.Value})
		if err != nil {
			return err
		}

		c.loadConst(index)
	} else {
		return errors.New("compiler: expected an identifier to the right of a dot (.)")
	}

	c.push(bytecode.LoadSubscript)

	return nil
}

func (c *Compiler) compileAssignOrDeclare(l, right ast.Expression, t string) error {
	switch left := l.(type) {
	case *ast.Identifier:
		return c.compileAssignToIdent(left, right, t)

	case *ast.Call:
		if list, ok := left.Argument.(*ast.List); ok {
			if len(list.Value) != 1 {
				return fmt.Errorf("compiler: exactly one element should be present in an index assignment: a[b] = c")
			}

			return c.compileAssignToIndex(left.Function, list.Value[0], right, t)
		}

		return c.compileAssignToFunction(left, right, t)
	}

	return nil
}

func (c *Compiler) compileAssignToIdent(ident *ast.Identifier, right ast.Expression, t string) error {
	if err := c.CompileExpression(right); err != nil {
		return err
	}

	index, err := c.addName(ident.Value)
	if err != nil {
		return err
	}

	low, high := runeToBytes(rune(index))
	if t == "assign" {
		c.push(bytecode.StoreName, high, low)
	} else {
		c.push(bytecode.DeclareName, high, low)
	}

	return nil
}

func (c *Compiler) compileAssignToIndex(obj, idx, val ast.Expression, t string) error {
	if t != "assign" {
		return errors.New("compiler: cannot declare to a subscript[expression], use an assignment instead: a[b] = c")
	}

	if err := c.CompileExpression(val); err != nil {
		return err
	}

	if err := c.CompileExpression(obj); err != nil {
		return err
	}

	if err := c.CompileExpression(idx); err != nil {
		return err
	}

	c.push(bytecode.StoreSubscript)

	return nil
}

func (c *Compiler) compileAssignToFunction(function *ast.Call, body ast.Expression, t string) error {
	fn := &object.Function{}

	// Set function parameters
	if argTuple, ok := function.Argument.(*ast.Infix); ok && argTuple.Operator == "," {
		for _, arg := range c.expandTuple(argTuple) {
			if id, ok := arg.(*ast.Identifier); ok {
				fn.Parameters = append(fn.Parameters, id.Value)
			} else {
				return errors.New("compiler: function parameters must be identifiers -- pattern matching is not supported (yet?)")
			}
		}
	} else if arg, ok := function.Argument.(*ast.Identifier); ok {
		fn.Parameters = append(fn.Parameters, arg.Value)
	} else {
		return errors.New("compiler: function parameters must be identifiers -- pattern matching is not supported (yet?)")
	}

	// Compile the function body in a new Compiler instance
	subCompiler := New()
	subCompiler.CompileExpression(body)

	code, err := bytecode.Read(bytes.NewReader(subCompiler.Bytes))
	if err != nil {
		return err
	}

	fn.Code = code
	fn.Constants = subCompiler.Constants
	fn.Names = subCompiler.Names
	fn.Jumps = subCompiler.Jumps

	c.addAndLoad(fn)

	switch name := function.Function.(type) {
	case *ast.Identifier:
		index, err := c.addName(name.Value)
		if err != nil {
			return err
		}

		low, high := runeToBytes(rune(index))
		if t == "assign" {
			c.push(bytecode.StoreName, high, low)
		} else {
			c.push(bytecode.DeclareName, high, low)
		}

	case *ast.Infix:
		if t != "assign" {
			return errors.New("compiler: cannot declare a function to a subscript expression, use an assignment instead: a.b <params> = <body>")
		}

		if name.Operator != "." {
			return errors.New("compiler: can only define functions as identifiers or model methods")
		}

		if err := c.CompileExpression(name.Left); err != nil {
			return err
		}

		if id, ok := name.Right.(*ast.Identifier); ok {
			c.addAndLoad(&object.String{Value: id.Value})
		} else {
			return errors.New("compiler: expected an identifier to the right of a dot (.)")
		}

		c.push(bytecode.StoreSubscript)

	default:
		return errors.New("compiler: can only define functions as identifiers or model methods")
	}

	return nil
}

func (c *Compiler) compilePrefix(node *ast.Prefix) error {
	if err := c.CompileExpression(node.Right); err != nil {
		return err
	}

	if node.Operator == "+" {
		return nil
	}

	op := map[string]byte{
		"-": bytecode.UnaryNegate,
		"!": bytecode.UnaryInvert,
	}[node.Operator]

	c.push(op)

	return nil
}
