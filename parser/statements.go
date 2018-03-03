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
