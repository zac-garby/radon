package ast

// A Node is the interface from which all AST nodes implement.
type Node interface{}

// A Statement is a node which doesn't evaluate to a value, for example
// a loop.
type Statement interface {
	Node
	Stmt()
}

// An Expression is a node which evaluates to a value, for example a
// number literal.
type Expression interface {
	Node
	Expr()
}

// A Program is a list of statements which, usually, represents an entire
// file.
type Program struct {
	Statements []Statement
}

// Tree returns a nicely formatted tree representation of the program.
func (p *Program) Tree() string {
	str := ""

	for _, stmt := range p.Statements {
		str += Tree(stmt, 0, "") + "\n"
	}

	return str
}
