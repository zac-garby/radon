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

func (p *Parser) parsePrefix() ast.Expression {
	node := &ast.Prefix{
		Operator: p.cur.Literal,
	}

	p.next()
	node.Right = p.parseExpression(prefix)

	return node
}

func (p *Parser) parseIf() ast.Expression {
	p.next()

	node := &ast.If{
		Condition: p.parseExpression(lowest),
	}

	if p.peekIs(token.Do) {
		p.next()
		node.Consequence = p.parseBlock()
	} else {
		if !p.expect(token.Then) {
			return nil
		}

		p.next()
		node.Consequence = p.parseExpression(lowest)
	}

	if p.peekIs(token.Else) {
		p.next()
		p.next()
		node.Alternative = p.parseExpression(lowest)
	} else {
		node.Alternative = &ast.Nil{}
	}

	return node
}

func (p *Parser) parseMatch() ast.Expression {
	p.next()
	node := &ast.Match{
		Input: p.parseExpression(lowest),
	}

	if !p.expect(token.Where) {
		return nil
	}

	for p.peekIs(token.BitOr) {
		pair := ast.MatchBranch{}

		p.next()
		p.next()

		pair.Condition = p.parseExpression(lowest)

		if !p.expect(token.RightArrow) {
			return nil
		}

		p.next()

		pair.Body = p.parseExpression(join)
		node.Branches = append(node.Branches, pair)

		if p.peekIs(token.Comma) {
			p.next()
		} else {
			break
		}
	}

	hasWildcard := false

	for _, branch := range node.Branches {
		if id, ok := branch.Condition.(*ast.Identifier); ok && id.Value == "_" {
			hasWildcard = true
			break
		}
	}

	if !hasWildcard {
		node.Branches = append(node.Branches, ast.MatchBranch{
			Condition: &ast.Identifier{Value: "_"},
			Body:      &ast.Nil{},
		})
	}

	return node
}

func (p *Parser) parseModel() ast.Expression {
	if !p.expect(token.LeftParen) {
		return nil
	}

	node := &ast.Model{
		Parameters: p.parseParams(token.RightParen, token.Comma),
	}

	if p.peekIs(token.BitOr) {
		p.next()
		p.next()

		node.Parent = p.parseExpression(lowest)

		if p.peekIs(token.LeftParen) {
			p.next()
			node.ParentParameters = p.parseExpressionList(token.RightParen, token.Comma)
		}
	}

	return node
}

func (p *Parser) parseNudLambda() ast.Expression {
	p.next()

	return &ast.Lambda{
		Body: p.parseExpression(lowest),
	}
}

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	node := &ast.Infix{
		Operator: p.cur.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.next()
	node.Right = p.parseExpression(precedence)

	return node
}

func (p *Parser) parseIndex(left ast.Expression) ast.Expression {
	p.next()

	node := &ast.Index{
		Left:  left,
		Right: p.parseExpression(lowest),
	}

	if !p.expect(token.RightSquare) {
		return nil
	}

	return node
}
