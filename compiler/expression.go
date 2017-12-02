package compiler

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Zac-Garby/lang/ast"
	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/object"
)

// CompileExpression compiles an AST expression.
func (c *Compiler) CompileExpression(e ast.Expression) error {
	switch node := e.(type) {
	case *ast.InfixExpression:
		return c.compileInfix(node)
	case *ast.PrefixExpression:
		return c.compilePrefix(node)
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
	case *ast.List:
		return c.compileList(node)
	case *ast.Tuple:
		return c.compileTuple(node)
	case *ast.Map:
		return c.compileMap(node)
	case *ast.IfExpression:
		return c.compileIf(node)
	case *ast.FunctionCall:
		return c.compileFnCall(node)
	case *ast.IndexExpression:
		return c.compileIndex(node)
	case *ast.Block:
		return c.compileBlock(node)
	case *ast.Match:
		return c.compileMatch(node)
	case *ast.Model:
		return c.compileModel(node)
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(e))
	}
}

func (c *Compiler) compileNumber(node *ast.Number) error {
	var (
		obj        = &object.Number{Value: node.Value}
		index, err = c.addConst(obj)
	)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileString(node *ast.String) error {
	var (
		obj        = &object.String{Value: node.Value}
		index, err = c.addConst(obj)
	)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileBoolean(node *ast.Boolean) error {
	var (
		obj        = &object.Boolean{Value: node.Value}
		index, err = c.addConst(obj)
	)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileNil(node *ast.Nil) error {
	var (
		obj        = object.NilObj
		index, err = c.addConst(obj)
	)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileIdentifier(node *ast.Identifier) error {
	return c.compileName(node.Value)
}

func (c *Compiler) compileInfix(node *ast.InfixExpression) error {
	left, right := node.Left, node.Right

	// Some operators are handled differently
	switch node.Operator {
	case "=":
		return c.compileAssign(left, right)

	case ":=":
		return c.compileDeclare(left, right)

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
		"-":  bytecode.BinarySubtract,
		"*":  bytecode.BinaryMultiply,
		"/":  bytecode.BinaryDivide,
		"^":  bytecode.BinaryExponent,
		"//": bytecode.BinaryFloorDiv,
		"%":  bytecode.BinaryFloorDiv,
		"||": bytecode.BinaryOr,
		"&&": bytecode.BinaryAnd,
		"|":  bytecode.BinaryBitOr,
		"&":  bytecode.BinaryBitAnd,
		"==": bytecode.BinaryEquals,
		"!=": bytecode.BinaryNotEqual,
		"<":  bytecode.BinaryLessThan,
		">":  bytecode.BinaryMoreThan,
		"<=": bytecode.BinaryLessEq,
		">=": bytecode.BinaryMoreEq,
	}[node.Operator]

	if !ok {
		return fmt.Errorf("compiler: operator %s not yet implemented", node.Operator)
	}

	c.push(op)

	return nil
}

func (c *Compiler) compileAssign(l, right ast.Expression) error {
	switch left := l.(type) {
	case *ast.Identifier:
		if err := c.CompileExpression(right); err != nil {
			return err
		}

		index, err := c.addName(left.Value)
		if err != nil {
			return err
		}

		low, high := runeToBytes(rune(index))
		c.push(bytecode.StoreName, high, low)

	case *ast.IndexExpression:
		if err := c.CompileExpression(right); err != nil {
			return err
		}

		if err := c.CompileExpression(left.Collection); err != nil {
			return err
		}

		if err := c.CompileExpression(left.Index); err != nil {
			return err
		}

		c.push(bytecode.StoreField)

	case *ast.InfixExpression:
		if err := c.CompileExpression(right); err != nil {
			return err
		}

		if left.Operator != "." {
			return fmt.Errorf("compiler: cannot assign to an infix expression with operator: %s", left.Operator)
		}

		if err := c.CompileExpression(left.Left); err != nil {
			return err
		}

		if id, ok := left.Right.(*ast.Identifier); ok {
			obj := &object.String{Value: id.Value}
			index, err := c.addConst(obj)
			if err != nil {
				return err
			}

			low, high := runeToBytes(rune(index))
			c.push(bytecode.LoadConst, high, low)
		} else {
			return errors.New("compiler: expected an identifier to the right of a dot (.)")
		}

		c.push(bytecode.StoreField)

	case *ast.FunctionCall:
		if _, ok := left.Function.(*ast.Identifier); !ok {
			return errors.New("compiler: cannot define a function whose name isn't an identifier")
		}

		fn := &object.Function{}

		for _, arg := range left.Arguments {
			if id, ok := arg.(*ast.Identifier); ok {
				fn.Parameters = append(fn.Parameters, id.Value)
			} else {
				return errors.New("compiler: function parameters must be identifiers")
			}
		}

		fnComp := New()
		fnComp.CompileExpression(right)

		code, err := bytecode.Read(fnComp.Bytes)
		if err != nil {
			return err
		}

		fn.Code = code
		fn.Constants = fnComp.Constants
		fn.Names = fnComp.Names

		index, err := c.addConst(fn)
		if err != nil {
			return err
		}

		low, high := runeToBytes(rune(index))
		c.push(bytecode.LoadConst, high, low)

		index, err = c.addName(left.Function.(*ast.Identifier).Value)
		if err != nil {
			return err
		}

		low, high = runeToBytes(rune(index))
		c.push(bytecode.StoreName, high, low)

	default:
		return errors.New("compiler: can only assign to identifiers, functions and field access expressions")
	}

	return nil
}

func (c *Compiler) compileModel(node *ast.Model) error {
	model := &object.Model{}

	for _, param := range node.Parameters {
		if id, ok := param.(*ast.Identifier); ok {
			model.Parameters = append(model.Parameters, id.Value)
		} else {
			return errors.New("compiler: all model parameters should be identifiers")
		}
	}

	index, err := c.addConst(model)
	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileDeclare(l, right ast.Expression) error {
	switch left := l.(type) {
	case *ast.Identifier:
		if err := c.CompileExpression(right); err != nil {
			return err
		}

		index, err := c.addName(left.Value)
		if err != nil {
			return err
		}

		low, high := runeToBytes(rune(index))
		c.push(bytecode.DeclareName, high, low)

	case *ast.FunctionCall:
		if _, ok := left.Function.(*ast.Identifier); !ok {
			return errors.New("compiler: cannot define a function whose name isn't an identifier")
		}

		fn := &object.Function{}

		for _, arg := range left.Arguments {
			if id, ok := arg.(*ast.Identifier); ok {
				fn.Parameters = append(fn.Parameters, id.Value)
			} else {
				return errors.New("compiler: function parameters must be identifiers")
			}
		}

		fnComp := New()
		fnComp.CompileExpression(right)

		code, err := bytecode.Read(fnComp.Bytes)
		if err != nil {
			return err
		}

		fn.Code = code
		fn.Constants = fnComp.Constants
		fn.Names = fnComp.Names

		index, err := c.addConst(fn)
		if err != nil {
			return err
		}

		low, high := runeToBytes(rune(index))
		c.push(bytecode.LoadConst, high, low)

		index, err = c.addName(left.Function.(*ast.Identifier).Value)
		if err != nil {
			return err
		}

		low, high = runeToBytes(rune(index))
		c.push(bytecode.DeclareName, high, low)

	default:
		return errors.New("compiler: can only declare to functions and identifiers")
	}

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
		return errors.New("compiler: expected an identifier to the right of a dot")
	}

	c.push(bytecode.LoadField)

	return nil
}

func (c *Compiler) compilePrefix(node *ast.PrefixExpression) error {
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

func (c *Compiler) compileIf(node *ast.IfExpression) error {
	if err := c.CompileExpression(node.Condition); err != nil {
		return err
	}

	// JumpIfFalse (82) with 2 empty argument bytes
	c.push(bytecode.JumpIfFalse, 0, 0)
	condJump := len(c.Bytes) - 3

	if err := c.CompileExpression(node.Consequence); err != nil {
		return err
	}

	var skipJump int

	if node.Alternative != nil {
		// Jump past the alternative
		c.push(bytecode.Jump, 0, 0)
		skipJump = len(c.Bytes) - 3
	}

	// Set the jump target after the conditional
	condIndex := rune(len(c.Bytes))
	low, high := runeToBytes(condIndex)
	c.Bytes[condJump+1] = high
	c.Bytes[condJump+2] = low

	if node.Alternative != nil {
		if err := c.CompileExpression(node.Alternative); err != nil {
			return err
		}

		// Set the jump target after the conditional
		skipIndex := rune(len(c.Bytes))
		low, high = runeToBytes(skipIndex)
		c.Bytes[skipJump+1] = high
		c.Bytes[skipJump+2] = low
	}

	return nil
}

func (c *Compiler) compileList(node *ast.List) error {
	for _, elem := range node.Elements {
		if err := c.CompileExpression(elem); err != nil {
			return err
		}
	}

	low, high := runeToBytes(rune(len(node.Elements)))

	c.push(bytecode.MakeList, high, low)

	return nil
}

func (c *Compiler) compileTuple(node *ast.Tuple) error {
	for _, elem := range node.Value {
		if err := c.CompileExpression(elem); err != nil {
			return err
		}
	}

	low, high := runeToBytes(rune(len(node.Value)))

	c.push(bytecode.MakeTuple, high, low)

	return nil
}

func (c *Compiler) compileMap(node *ast.Map) error {
	for key, val := range node.Pairs {
		if err := c.CompileExpression(key); err != nil {
			return err
		}

		if err := c.CompileExpression(val); err != nil {
			return err
		}
	}

	low, high := runeToBytes(rune(len(node.Pairs)))

	c.push(bytecode.MakeMap, high, low)

	return nil
}

func (c *Compiler) compileFnCall(node *ast.FunctionCall) error {
	if id, ok := node.Function.(*ast.Identifier); ok {
		if builtin, ok := getBuiltin(id.Value); ok {
			if len(node.Arguments) == builtin.parameters {
				if builtin.autoPush {
					for _, arg := range node.Arguments {
						if err := c.CompileExpression(arg); err != nil {
							return err
						}
					}
				}
			} else {
				return fmt.Errorf("syntax: wrong amount of arguments to builtin function '%s'. expected %v", builtin.name, builtin.parameters)
			}

			return builtin.compile(c, node.Arguments)
		}
	}

	count := 0

	for _, arg := range node.Arguments {
		count++

		if err := c.CompileExpression(arg); err != nil {
			return err
		}
	}

	if err := c.CompileExpression(node.Function); err != nil {
		return err
	}

	index, err := c.addConst(object.MakeObj(count))
	if err != nil {
		return err
	}

	c.loadConst(index)
	c.push(bytecode.CallFn)

	return nil
}

func (c *Compiler) compileIndex(node *ast.IndexExpression) error {
	if err := c.CompileExpression(node.Collection); err != nil {
		return err
	}

	if err := c.CompileExpression(node.Index); err != nil {
		return err
	}

	c.push(bytecode.LoadField)

	return nil
}

func (c *Compiler) compileBlock(node *ast.Block) error {
	for _, stmt := range node.Statements {
		if err := c.CompileStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) compileMatch(node *ast.Match) error {
	var endJumps []int

	for _, branch := range node.Branches {
		var (
			branchJump int
			wildcard   = false
		)

		if id, ok := branch.Condition.(*ast.Identifier); ok && id.Value == "_" {
			wildcard = true
		}

		if !wildcard {
			if err := c.CompileExpression(node.Input); err != nil {
				return err
			}

			if err := c.CompileExpression(branch.Condition); err != nil {
				return err
			}

			c.push(bytecode.BinaryEquals, bytecode.JumpIfFalse, 0, 0)
			branchJump = len(c.Bytes) - 3
		}

		if err := c.CompileExpression(branch.Body); err != nil {
			return err
		}

		c.push(bytecode.Jump, 0, 0)
		endJumps = append(endJumps, len(c.Bytes)-3)

		if !wildcard {
			branchEnd := len(c.Bytes) - 1
			low, high := runeToBytes(rune(branchEnd))
			c.Bytes[branchJump+1] = high
			c.Bytes[branchJump+2] = low
		}
	}

	end := len(c.Bytes) - 1

	for _, jmp := range endJumps {
		low, high := runeToBytes(rune(end))
		c.Bytes[jmp+1] = high
		c.Bytes[jmp+2] = low
	}

	return nil
}
