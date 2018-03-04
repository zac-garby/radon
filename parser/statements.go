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
