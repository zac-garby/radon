package parser

import (
	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/token"
)

type nud func() ast.Expression
type led func(ast.Expression) ast.Expression

// A Parser takes a token generator function and parses it into an AST.
type Parser struct {
	// Errors contains any errors encountered during parsing.
	Errors []error

	lex       func() token.Token
	cur, peek token.Token
	nuds      map[token.Type]nud
	leds      map[token.Type]led
}

// New creates a new parser for the given token generator function.
func New(lex func() token.Token) *Parser {
	p := &Parser{
		lex:    lex,
		Errors: make([]error, 0),
	}

	p.nuds = map[token.Type]nud{}
	p.leds = map[token.Type]led{}

	p.next()
	p.next()

	return p
}

// Parse parses an entire program into an `ast.Program`. Also, returns the
// first error encountered, if any.
func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{
		Statements: make([]ast.Statement, 0, 10),
	}

	for !p.curIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}

		p.next()
	}

	return prog
}
