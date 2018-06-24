package parser

import (
	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/token"
)

// The precedences of different operators. The actual operators are described
// in the `precedences` variable.
const (
	lowest = iota
	assign
	lambda
	join
	or
	and
	bitOr
	bitAnd
	equals
	compare
	sum
	product
	exp
	prefix
	index
)

var precedences = map[token.Type]int{
	token.Assign:         assign,
	token.Declare:        assign,
	token.AndEquals:      assign,
	token.BitAndEquals:   assign,
	token.BitOrEquals:    assign,
	token.ExpEquals:      assign,
	token.FloorDivEquals: assign,
	token.MinusEquals:    assign,
	token.ModEquals:      assign,
	token.OrEquals:       assign,
	token.PlusEquals:     assign,
	token.SlashEquals:    assign,
	token.StarEquals:     assign,
	token.Or:             or,
	token.And:            and,
	token.BitOr:          bitOr,
	token.BitAnd:         bitAnd,
	token.Equal:          equals,
	token.NotEqual:       equals,
	token.LessThan:       compare,
	token.GreaterThan:    compare,
	token.LessThanEq:     compare,
	token.GreaterThanEq:  compare,
	token.Plus:           sum,
	token.Minus:          sum,
	token.Star:           product,
	token.Slash:          product,
	token.Mod:            product,
	token.Exp:            exp,
	token.FloorDiv:       exp,
	token.Bang:           prefix,
	token.Dot:            index,
	token.LambdaArrow:    lambda,
	token.Comma:          join,
}

// argTokens is the set of tokens which can appear as the first token to a function call
var argTokens = []token.Type{
	token.Number,
	token.String,
	token.ID,
	token.LeftParen,
	token.LeftSquare,
	token.LeftBrace,
	token.Bang,
	token.True,
	token.False,
	token.Nil,
	token.If,
	token.While,
	token.For,
	token.Match,
	token.Model,
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
		p.err(
			"illegal token encountered. literal: `%s`",
			p.peek.Start,
			p.peek.End,
			p.peek.Literal,
		)
	}
}

func (p *Parser) curIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.cur.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) peekIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.peek.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) expect(t token.Type) bool {
	if p.peekIs(t) {
		p.next()
		return true
	}

	p.peekErr(t)
	return false
}

func (p *Parser) parseExpressionList(end, sep token.Type) []ast.Expression {
	var exprs []ast.Expression

	if p.peekIs(end) {
		p.next()
		return exprs
	}

	p.next()
	exprs = append(exprs, p.parseExpression(join))

	for p.peekIs(sep) {
		p.next()

		if p.peekIs(end) {
			p.next()
			return exprs
		}

		p.next()
		exprs = append(exprs, p.parseExpression(join))
	}

	if !p.expect(end) {
		return nil
	}

	return exprs
}

func (p *Parser) parseExpressionPairs(end, sep token.Type) map[ast.Expression]ast.Expression {
	pairs := make(map[ast.Expression]ast.Expression)

	p.next()

	if p.curIs(end) {
		return pairs
	}

	key, val := p.parsePair()
	pairs[key] = val

	for p.peekIs(sep) {
		p.next()

		if p.peekIs(end) {
			p.next()
			return pairs
		}

		p.next()
		key, val = p.parsePair()
		pairs[key] = val
	}

	if !p.expect(end) {
		return nil
	}

	return pairs
}

func (p *Parser) parsePair() (ast.Expression, ast.Expression) {
	key := p.parseExpression(index)

	if !p.expect(token.Colon) {
		return nil, nil
	}

	p.next()

	value := p.parseExpression(join)

	return key, value
}
