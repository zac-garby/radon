package parser

import (
	"github.com/Zac-Garby/radon/ast"
	"github.com/Zac-Garby/radon/token"
)

func (p *Parser) parseStatement() ast.Statement {
	var node ast.Statement

	switch p.cur.Type {
	case token.Semi:
		return nil

	case token.Return:
		return p.parseReturn()

	case token.Break:
		return new(ast.Break)

	case token.Next:
		return new(ast.Next)

	case token.While:
		return p.parseWhile()

	case token.For:
		return p.parseFor()

	case token.Import:
		return p.parseImport()

	case token.Export:
		return p.parseExport()

	default:
		node = p.parseExpressionStatement()
	}

	if !p.expect(token.Semi) {
		return nil
	}

	return node
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	return &ast.ExpressionStatement{
		Expr: p.parseExpression(lowest),
	}
}

func (p *Parser) parseReturn() ast.Statement {
	if p.peekIs(token.Semi) {
		return &ast.Return{
			Value: &ast.Nil{},
		}
	}

	p.next()

	return &ast.Return{
		Value: p.parseExpression(lowest),
	}
}

func (p *Parser) parseWhile() ast.Statement {
	p.next()

	node := &ast.While{
		Condition: p.parseExpression(join),
	}

	if p.peekIs(token.Do) {
		p.next()
		node.Body = p.parseBlock()
	} else {
		if !p.expect(token.Comma) {
			return nil
		}

		p.next()
		node.Body = p.parseExpression(lowest)
	}

	return node
}

func (p *Parser) parseFor() ast.Statement {
	p.next()

	node := &ast.For{
		Var: p.parseExpression(lowest),
	}

	if !p.expect(token.In) {
		return nil
	}

	p.next()

	node.Collection = p.parseExpression(join)

	if p.peekIs(token.Do) {
		p.next()
		node.Body = p.parseBlock()
	} else {
		if !p.expect(token.Comma) {
			return nil
		}

		p.next()
		node.Body = p.parseExpression(lowest)
	}

	return node
}

func (p *Parser) parseImport() ast.Statement {
	if !p.expect(token.String) {
		return nil
	}

	str := p.parseExpression(lowest).(*ast.String)

	return &ast.Import{
		Path: str.Value,
	}
}

func (p *Parser) parseExport() ast.Statement {
	p.next()

	return &ast.Export{
		Names: p.parseExpression(lowest),
	}
}
