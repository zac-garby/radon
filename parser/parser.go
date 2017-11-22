package parser

import (
	"fmt"

	"github.com/Zac-Garby/lang/ast"
	"github.com/Zac-Garby/lang/lexer"
	"github.com/Zac-Garby/lang/token"
)

type prefixParser func() ast.Expression
type infixParser func(ast.Expression) ast.Expression

// Parser parses a string into an
// abstract syntax tree
type Parser struct {
	Errors []error

	lex       func() token.Token
	text      string
	cur, peek token.Token
	prefixes  map[token.Type]prefixParser
	infixes   map[token.Type]infixParser
}

// New returns a new parser for the
// given string
func New(text, file string) *Parser {
	p := &Parser{
		lex:    lexer.Lexer(text, file),
		text:   text,
		Errors: []error{},
	}

	p.prefixes = map[token.Type]prefixParser{
		token.ID:         p.parseID,
		token.Number:     p.parseNum,
		token.True:       p.parseBool,
		token.False:      p.parseBool,
		token.Nil:        p.parseNil,
		token.String:     p.parseString,
		token.LeftSquare: p.parseList,
		token.Map:        p.parseMap,
		token.Set:        p.parseSet,
		token.Minus:      p.parsePrefix,
		token.Plus:       p.parsePrefix,
		token.Bang:       p.parsePrefix,
		token.LeftParen:  p.parseGroupedExpression,
		token.If:         p.parseIfExpression,
		token.Match:      p.parseMatchExpression,
		token.LeftBrace:  p.parseBlock,
		token.TypeK:      p.parseType,
	}

	p.infixes = map[token.Type]infixParser{
		token.Plus:           p.parseInfix,
		token.Minus:          p.parseInfix,
		token.Star:           p.parseInfix,
		token.Slash:          p.parseInfix,
		token.Equal:          p.parseInfix,
		token.NotEqual:       p.parseInfix,
		token.LessThan:       p.parseInfix,
		token.GreaterThan:    p.parseInfix,
		token.Or:             p.parseInfix,
		token.And:            p.parseInfix,
		token.BitOr:          p.parseInfix,
		token.BitAnd:         p.parseInfix,
		token.Exp:            p.parseInfix,
		token.FloorDiv:       p.parseInfix,
		token.Mod:            p.parseInfix,
		token.LessThanEq:     p.parseInfix,
		token.GreaterThanEq:  p.parseInfix,
		token.AndEquals:      p.parseInfix,
		token.BitAndEquals:   p.parseInfix,
		token.BitOrEquals:    p.parseInfix,
		token.ExpEquals:      p.parseInfix,
		token.FloorDivEquals: p.parseInfix,
		token.MinusEquals:    p.parseInfix,
		token.ModEquals:      p.parseInfix,
		token.OrEquals:       p.parseInfix,
		token.PlusEquals:     p.parseInfix,
		token.SlashEquals:    p.parseInfix,
		token.StarEquals:     p.parseInfix,
		token.Assign:         p.parseInfix,
		token.Declare:        p.parseInfix,
		token.LeftSquare:     p.parseIndex,
		token.LeftParen:      p.parseFunctionCall,
	}

	p.next()
	p.next()

	return p
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peek.Type]; ok {
		return precedence
	}

	return lowest
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.cur.Type]; ok {
		return precedence
	}

	return lowest
}

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.lex()

	if p.peek.Type == token.Illegal {
		p.Err(
			fmt.Sprintf("illegal token found: `%s`", p.peek.Literal),
			p.peek.Start,
			p.peek.End,
		)
	}
}

// Parse parses an entire program
func (p *Parser) Parse() ast.Program {
	prog := ast.Program{
		Statements: []ast.Statement{},
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
