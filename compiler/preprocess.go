package compiler

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Zac-Garby/lang/ast"
)

// Preprocess preprocesses an AST node.
func Preprocess(node ast.Node) (ast.Node, error) {
	if expr, ok := node.(ast.Expression); ok {
		return preprocessExpression(expr)
	} else if stmt, ok := node.(ast.Statement); ok {
		return preprocessStatement(stmt)
	} else {
		return nil, errors.New("preprocesser: somehow, a node was neither an expression or a statement")
	}
}

// PreprocessReduce preprocesses an AST node until it can't be processed any more.
func PreprocessReduce(node ast.Node) (ast.Node, error) {
	for {
		p, err := Preprocess(node)
		if err != nil {
			return nil, err
		}

		// A bit hacky -- compare Tree() output for the two nodes to check for
		// equality. Could definitely be improved if I ever make an easy way to
		// walk an AST.
		var (
			before = ast.Tree(node, 0, "")
			after  = ast.Tree(p, 0, "")
		)

		if before == after {
			break
		}
	}

	return node, nil
}

// preprocessExpression preprocesses an expression.
func preprocessExpression(n ast.Expression) (ast.Expression, error) {
	switch node := n.(type) {
	case *ast.Identifier,
		*ast.Number,
		*ast.Boolean,
		*ast.Nil,
		*ast.String:
		return node, nil

	case *ast.Tuple:
		for i, elem := range node.Value {
			p, err := preprocessExpression(elem)
			if err != nil {
				return nil, err
			}

			node.Value[i] = p
		}

		return node, nil

	case *ast.List:
		for i, elem := range node.Elements {
			p, err := preprocessExpression(elem)
			if err != nil {
				return nil, err
			}

			node.Elements[i] = p
		}

		return node, nil

	case *ast.Map:
		pairs := make(map[ast.Expression]ast.Expression)

		for k, v := range node.Pairs {
			pk, err := preprocessExpression(k)
			if err != nil {
				return nil, err
			}

			pv, err := preprocessExpression(v)
			if err != nil {
				return nil, err
			}

			pairs[pk] = pv
		}

		node.Pairs = pairs
		return node, nil

	case *ast.Block:
		for i, stmt := range node.Statements {
			p, err := preprocessStatement(stmt)
			if err != nil {
				return nil, err
			}

			node.Statements[i] = p
		}

		return node, nil

	case *ast.PrefixExpression:
		p, err := preprocessExpression(node.Right)
		if err != nil {
			return nil, err
		}

		node.Right = p

		return node, nil

	case *ast.InfixExpression:
		lp, err := preprocessExpression(node.Left)
		if err != nil {
			return nil, err
		}

		rp, err := preprocessExpression(node.Right)
		if err != nil {
			return nil, err
		}

		node.Left = lp
		node.Right = rp

		return node, nil

	case *ast.IndexExpression:
		cp, err := preprocessExpression(node.Collection)
		if err != nil {
			return nil, err
		}

		ip, err := preprocessExpression(node.Index)
		if err != nil {
			return nil, err
		}

		node.Collection = cp
		node.Index = ip

		return node, nil

	case *ast.FunctionCall:
		fp, err := preprocessExpression(node.Function)
		if err != nil {
			return nil, err
		}

		node.Function = fp

		for i, arg := range node.Arguments {
			p, err := preprocessExpression(arg)
			if err != nil {
				return nil, err
			}

			node.Arguments[i] = p
		}

		return node, nil

	case *ast.IfExpression:
		condp, err := preprocessExpression(node.Condition)
		if err != nil {
			return nil, err
		}

		conqp, err := preprocessExpression(node.Consequence)
		if err != nil {
			return nil, err
		}

		altp, err := preprocessExpression(node.Alternative)
		if err != nil {
			return nil, err
		}

		node.Condition = condp
		node.Consequence = conqp
		node.Alternative = altp

		return node, nil

	case *ast.Match:
		ip, err := preprocessExpression(node.Input)
		if err != nil {
			return nil, err
		}

		node.Input = ip

		for i, branch := range node.Branches {
			cp, err := preprocessExpression(branch.Condition)
			if err != nil {
				return nil, err
			}

			bp, err := preprocessExpression(branch.Body)
			if err != nil {
				return nil, err
			}

			branch.Condition = cp
			branch.Body = bp
			node.Branches[i] = branch
		}

		return node, nil

	case *ast.Model:
		for i, param := range node.Parameters {
			pp, err := preprocessExpression(param)
			if err != nil {
				return nil, err
			}

			node.Parameters[i] = pp
		}

		return node, nil

	default:
		return nil, fmt.Errorf("preprocessor: not implemented for node %s", reflect.TypeOf(node))
	}
}

// preprocessStatement preprocesses a statement.
func preprocessStatement(n ast.Statement) (ast.Statement, error) {
	switch node := n.(type) {
	case *ast.ExpressionStatement:
		p, err := preprocessExpression(node.Expr)
		if err != nil {
			return nil, err
		}

		node.Expr = p

		return node, nil

	case *ast.ReturnStatement:
		p, err := preprocessExpression(node.Value)
		if err != nil {
			return nil, err
		}

		node.Value = p

		return node, nil

	case *ast.NextStatement, *ast.BreakStatement:
		return node, nil

	case *ast.Loop:
		p, err := preprocessExpression(node.Body)
		if err != nil {
			return nil, err
		}

		node.Body = p

		return node, nil

	case *ast.WhileLoop:
		cp, err := preprocessExpression(node.Condition)
		if err != nil {
			return nil, err
		}

		bp, err := preprocessExpression(node.Body)
		if err != nil {
			return nil, err
		}

		node.Condition = cp
		node.Body = bp

		return node, nil

	case *ast.ForLoop:
		intp, err := preprocessExpression(node.Init)
		if err != nil {
			return nil, err
		}

		cp, err := preprocessExpression(node.Condition)
		if err != nil {
			return nil, err
		}

		incp, err := preprocessExpression(node.Increment)
		if err != nil {
			return nil, err
		}

		bp, err := preprocessExpression(node.Body)
		if err != nil {
			return nil, err
		}

		node.Init = intp
		node.Condition = cp
		node.Increment = incp
		node.Body = bp

		return node, nil

	default:
		return nil, fmt.Errorf("preprocessor: not implemented for node %s", reflect.TypeOf(node))
	}
}