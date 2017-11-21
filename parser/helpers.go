package parser

import (
	"github.com/Zac-Garby/lang/ast"
	"github.com/Zac-Garby/lang/token"
)

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

func (p *Parser) expect(ts ...token.Type) bool {
	if p.peekIs(ts...) {
		p.next()
		return true
	}

	p.peekErr(ts...)
	return false
}

func (p *Parser) expectCur(ts ...token.Type) bool {
	if p.curIs(ts...) {
		p.next()
		return true
	}

	p.curErr(ts...)
	return false
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	exprs := []ast.Expression{}

	if p.curIs(end) {
		return exprs
	}

	exprs = append(exprs, p.parseExpression(lowest))

	for p.peekIs(token.Comma) {
		p.next()

		if p.peekIs(end) {
			p.next()
			return exprs
		}

		p.next()
		exprs = append(exprs, p.parseExpression(lowest))
	}

	if !p.expect(end) {
		return nil
	}

	return exprs
}

func (p *Parser) parseExpressionPairs(end token.Type) map[ast.Expression]ast.Expression {
	pairs := map[ast.Expression]ast.Expression{}

	if p.curIs(token.Colon) {
		p.next()
		return pairs
	}

	key, val := p.parsePair()
	pairs[key] = val

	for p.peekIs(token.Comma) {
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

	value := p.parseExpression(lowest)

	return key, value
}

func (p *Parser) parseParams(end token.Type) []ast.Expression {
	params := []ast.Expression{}

	if p.peekIs(end) {
		p.next()
		return params
	}

	p.next()
	params = append(params, p.parseID())

	for p.peekIs(token.Comma) {
		p.next()
		p.next()
		params = append(params, p.parseID())
	}

	if !p.expect(end) {
		return nil
	}

	return params
}
