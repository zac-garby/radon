package parser

import (
	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/token"
)

// parseExpression parses an expression starting at the current token. It leaves
// cur on the last token of the expression.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	nud, ok := p.nuds[p.cur.Type]
	if !ok {
		p.unexpected(p.cur.Type)
		return nil
	}

	left := nud()

	for !p.peekIs(token.Semi) && precedence < p.peekPrecedence() {
		led, ok := p.leds[p.peek.Type]
		if !ok {
			return left
		}

		p.next()
		left = led(left)
	}

	return left
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Value: p.cur.Literal,
	}
}
