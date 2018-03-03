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

	p.nuds = map[token.Type]nud{
		token.ID:         p.parseIdentifier,
		token.Number:     p.parseNumber,
		token.True:       p.parseBoolean,
		token.False:      p.parseBoolean,
		token.Nil:        p.parseNil,
		token.String:     p.parseString,
		token.LeftParen:  p.parseGroupedExpression,
		token.LeftSquare: p.parseList,
		token.LeftBrace:  p.parseMap,
		token.Do:         p.parseBlock,
		token.Minus:      p.parsePrefix,
		token.Plus:       p.parsePrefix,
		token.Bang:       p.parsePrefix,
		token.If:         p.parseIf,
		token.Match:      p.parseMatch,
		token.Model:      p.parseModel,
	}

	p.leds = map[token.Type]led{}

	p.next()
	p.next()

	return p
}

// Parse parses an entire program into an `ast.Program`. Also, returns the
// first error encountered, if any.
func (p *Parser) Parse() (*ast.Program, error) {
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

	if len(p.Errors) > 0 {
		return nil, p.Errors[0]
	}

	return prog, nil
}
