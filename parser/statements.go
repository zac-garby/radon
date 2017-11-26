package parser

import (
	"github.com/Zac-Garby/lang/ast"
	"github.com/Zac-Garby/lang/token"
)

func (p *Parser) parseStatement() ast.Statement {
	var node ast.Statement

	switch p.cur.Type {
	case token.Semi:
		return nil

	case token.Return:
		node = p.parseReturn()

	case token.Break:
		node = p.parseBreak()

	case token.Next:
		node = p.parseNext()

	case token.Loop:
		node = p.parseLoop()

	case token.While:
		node = p.parseWhile()

	case token.For:
		node = p.parseFor()

	default:
		node = p.parseExpressionStmt()
	}

	if !p.expect(token.Semi) {
		return nil
	}

	return node
}

func (p *Parser) parseExpressionStmt() ast.Statement {
	return &ast.ExpressionStatement{
		Tok:  p.cur,
		Expr: p.parseExpression(lowest),
	}
}

func (p *Parser) parseReturn() ast.Statement {
	if p.peekIs(token.Semi) {
		return &ast.ReturnStatement{
			Tok:   p.cur,
			Value: &ast.Tuple{},
		}
	}

	node := &ast.ReturnStatement{
		Tok: p.cur,
	}

	p.next()
	node.Value = p.parseExpression(lowest)

	return node
}

func (p *Parser) parseBreak() ast.Statement {
	return &ast.BreakStatement{
		Tok: p.cur,
	}
}

func (p *Parser) parseNext() ast.Statement {
	return &ast.NextStatement{
		Tok: p.cur,
	}
}

func (p *Parser) parseLoop() ast.Statement {
	node := &ast.Loop{
		Tok: p.cur,
	}

	p.next()

	node.Body = p.parseExpression(lowest)

	return node
}

func (p *Parser) parseWhile() ast.Statement {
	node := &ast.WhileLoop{
		Tok: p.cur,
	}

	if !p.expect(token.LeftParen) {
		return nil
	}

	p.next()
	node.Condition = p.parseExpression(lowest)

	if !p.expect(token.RightParen) {
		return nil
	}

	p.next()
	node.Body = p.parseExpression(lowest)

	return node
}

func (p *Parser) parseFor() ast.Statement {
	node := &ast.ForLoop{
		Tok: p.cur,
	}

	if !p.expect(token.LeftParen) {
		return nil
	}

	p.next()
	node.Init = p.parseExpression(lowest)

	if !p.expect(token.Semi) {
		return nil
	}

	p.next()
	node.Condition = p.parseExpression(lowest)

	if !p.expect(token.Semi) {
		return nil
	}

	p.next()
	node.Increment = p.parseExpression(lowest)

	if !p.expect(token.RightParen) {
		return nil
	}

	p.next()

	node.Body = p.parseExpression(lowest)

	return node
}
