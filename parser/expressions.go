package parser

import (
	"strconv"

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

func (p *Parser) parseNumber() ast.Expression {
	val, _ := strconv.ParseFloat(p.cur.Literal, 64)

	return &ast.Number{
		Value: val,
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Value: p.cur.Type == token.True,
	}
}

func (p *Parser) parseNil() ast.Expression {
	return &ast.Nil{}
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{
		Value: p.cur.Literal,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.next()

	expr := p.parseExpression(lowest)

	if !p.expect(token.RightParen) {
		return nil
	}

	return expr
}

func (p *Parser) parseList() ast.Expression {
	return &ast.List{
		Value: p.parseExpressionList(token.RightSquare, token.Comma),
	}
}

func (p *Parser) parseMap() ast.Expression {
	return &ast.Map{
		Value: p.parseExpressionPairs(token.RightBrace, token.Comma),
	}
}

func (p *Parser) parseBlock() ast.Expression {
	node := &ast.Block{
		Value: make([]ast.Statement, 0, 8),
	}

	p.next()

	for !p.curIs(token.End) && !p.curIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			node.Value = append(node.Value, stmt)
		}

		p.next()
	}

	return node
}
