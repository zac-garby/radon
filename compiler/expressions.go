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
	case *ast.If:
		return c.compileIf(node)
	case *ast.List:
		return c.compileList(node)
	case *ast.Map:
		return c.compileMap(node)
	case *ast.Call:
		return c.compileCall(node)
	case *ast.Block:
		return c.compileBlock(node)
	case *ast.Match:
		return c.compileMatch(node)
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
	case ",":
		return c.compileCommaInfix(left, right)
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
	}[node.Operator]

	if !ok {
		return fmt.Errorf("compiler: operator %s not yet implemented", node.Operator)
	}

	c.push(op)

	return nil
}

func (c *Compiler) compileCommaInfix(left, right ast.Expression) error {
	if left == nil || right == nil {
		c.addAndLoad(&object.Tuple{})
		return nil
	}

	if err := c.CompileExpression(left); err != nil {
		return err
	}

	if err := c.CompileExpression(right); err != nil {
		return err
	}

	c.push(bytecode.BinaryTuple)

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
		",": bytecode.UnaryTuple,
	}[node.Operator]

	c.push(op)

	return nil
}

func (c *Compiler) compileIf(node *ast.If) error {
	if err := c.CompileExpression(node.Condition); err != nil {
		return err
	}

	// JumpIfFalse followed by two empty bytes for the argument
	c.push(bytecode.JumpUnless, 0, 0)
	condJump := len(c.Bytes) - 3

	c.pushScope()
	if err := c.CompileExpression(node.Consequence); err != nil {
		return err
	}

	var skipJump int

	if node.Alternative != nil {
		// Jump past the alternative
		c.push(bytecode.Jump, 0, 0)
		skipJump = len(c.Bytes) - 3
	}
	c.popScope()

	// Set the jump target after the conditional
	c.setJumpArg(condJump, len(c.Bytes)+1)

	if node.Alternative != nil {
		c.pushScope()
		if err := c.CompileExpression(node.Alternative); err != nil {
			return err
		}
		c.popScope()
	}

	// Set the jump target after a dummy byte
	c.push(bytecode.Nop)
	c.setJumpArg(skipJump, len(c.Bytes))

	return nil
}

func (c *Compiler) compileList(node *ast.List) error {
	for _, elem := range node.Value {
		if err := c.encloseExpression(elem); err != nil {
			return err
		}
	}

	low, high := runeToBytes(rune(len(node.Value)))
	c.push(bytecode.MakeList, high, low)

	return nil
}

func (c *Compiler) compileMap(node *ast.Map) error {
	for key, val := range node.Value {
		if err := c.encloseExpression(key); err != nil {
			return err
		}

		if err := c.encloseExpression(val); err != nil {
			return err
		}
	}

	low, high := runeToBytes(rune(len(node.Value)))
	c.push(bytecode.MakeMap, high, low)

	return nil
}

func (c *Compiler) compileCall(node *ast.Call) error {
	var args []ast.Expression

	if tupInf, ok := node.Argument.(*ast.Infix); ok && tupInf.Operator == "," {
		args = c.expandTuple(tupInf)
	} else {
		args = []ast.Expression{node.Argument}
	}

	// Iterate arguments in reverse order
	for i := len(args) - 1; i >= 0; i-- {
		if err := c.encloseExpression(args[i]); err != nil {
			return err
		}
	}

	if err := c.encloseExpression(node.Function); err != nil {
		return err
	}

	low, high := runeToBytes(rune(len(args)))
	c.push(bytecode.CallFunction, high, low)

	return nil
}

func (c *Compiler) compileBlock(node *ast.Block) error {
	c.pushScope()
	defer c.popScope()

	for _, stmt := range node.Value {
		if err := c.CompileStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) compileMatch(node *ast.Match) error {
	c.pushScope()
	defer c.popScope()

	if err := c.CompileExpression(node.Input); err != nil {
		return err
	}

	c.push(bytecode.StartMatch)

	var wildcard ast.Expression

	for _, branch := range node.Branches {
		if id, ok := branch.Condition.(*ast.Identifier); ok && id.Value == "_" {
			if wildcard != nil {
				return errors.New("compiler: only one wildcard branch is permitted per match-expression")
			}

			wildcard = branch.Body
			continue
		}

		if err := c.CompileExpression(branch.Condition); err != nil {
			return err
		}

		c.push(bytecode.StartBranch)

		if err := c.CompileExpression(branch.Body); err != nil {
			return err
		}

		c.push(bytecode.EndBranch)
	}

	if wildcard == nil {
		c.addAndLoad(&object.Nil{})
	} else {
		if err := c.CompileExpression(wildcard); err != nil {
			return err
		}
	}

	c.push(bytecode.EndMatch)

	return nil
}
