package ast_test

import (
	"testing"

	. "github.com/Zac-Garby/radon/ast"
)

// Since the AST package only defines data structures,
// and doesn't really do anything, it really only makes
// sense to check that the things defined in the package
// don't break.
func TestNoErrors(t *testing.T) {
	p := Program{
		Statements: []Statement{
			ExpressionStatement{
				Expr: InfixExpression{
					Left: Number{
						Value: 5.3,
					},
					Right: InfixExpression{
						Left: Number{
							Value: 3,
						},
						Right: Number{
							Value: 100,
						},
						Operator: "+",
					},
					Operator: "*",
				},
			},

			ReturnStatement{
				Value: Nil{},
			},
		},
	}

	p.Tree()
}
