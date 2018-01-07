package ast

import "github.com/Zac-Garby/radon/token"

// A Node is the interface from which all syntax
// tree nodes extend from.
type Node interface {
	Token() token.Token
}

// A Statement is a piece of code which doesn't
// evaluate to a value, e.g. a return statement.
type Statement interface {
	Node
	Stmt()
}

// An Expression is a piece of code which *does*
// evaluate to a value, such as an infix expression.
type Expression interface {
	Node
	Expr()
}

// A Program is a sequence of statements (keep in
// mind, though, that an expression can be abstracted
// into a statement.)
type Program struct {
	Statements []Statement
}

// Tree returns the tree representation of a program.
func (p *Program) Tree() string {
	str := ""

	for _, stmt := range p.Statements {
		str += Tree(stmt, 0, "") + "\n"
	}

	return str
}
