package parser

import (
	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/token"
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

	if p.peekIs(end) {
		p.next()
		return exprs
	}

	p.next()
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
	params := make([]ast.Expression, 0)

	if p.peekIs(token.RightParen) {
		p.next()
		return params
	}

	if !p.expect(token.ID) {
		return nil
	}

	id := p.parseID()
	params = append(params, id)

	for p.peekIs(token.Comma) {
		p.next()

		if p.peekIs(end) {
			break
		}

		if !p.expect(token.ID) {
			return nil
		}

		id := p.parseID()
		params = append(params, id)
	}

	if !p.expect(token.RightParen) {
		return nil
	}

	for i, ida := range params {
		for j, idb := range params {
			if i == j {
				continue
			}

			if ida.(*ast.Identifier).Value == idb.(*ast.Identifier).Value {
				p.defaultErr("identical parameters not allowed")
				return nil
			}
		}
	}

	return params
}
