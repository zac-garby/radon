package ast

import (
	"reflect"

	"github.com/Zac-Garby/lang/token"
)

type expr struct{}

func (e expr) Expr() {}

func (e expr) Token() token.Token {
	v := reflect.ValueOf(e).FieldByName("Tok")
	if v.IsValid() {
		return v.Interface().(token.Token)
	}

	return token.Token{}
}

type (
	// Identifier is an identifier
	Identifier struct {
		expr
		Tok   token.Token
		Value string
	}

	// Number is a number literal
	Number struct {
		expr
		Tok   token.Token
		Value float64
	}

	// Boolean is a boolean literal
	Boolean struct {
		expr
		Tok   token.Token
		Value bool
	}

	// Nil is the nil literal
	Nil struct {
		expr
		Tok token.Token
	}

	// String is a string literal
	String struct {
		expr
		Tok   token.Token
		Value string
	}

	// Tuple is a tuple literal
	Tuple struct {
		expr
		Tok   token.Token
		Value []Expression
	}

	// List is an list literal
	List struct {
		expr
		Tok      token.Token
		Elements []Expression
	}

	// Map is a map literal
	Map struct {
		expr
		Tok   token.Token
		Pairs map[Expression]Expression
	}

	// Set is a set literal
	Set struct {
		expr
		Tok      token.Token
		Elements []Expression
	}

	// A Block combines multiple statements into an expression
	Block struct {
		expr
		Tok        token.Token
		Statements []Statement
	}

	// PrefixExpression is a prefix operator expression
	PrefixExpression struct {
		expr
		Tok      token.Token
		Operator string
		Right    Expression
	}

	// An InfixExpression is an infix operator expression
	InfixExpression struct {
		expr
		Tok         token.Token
		Operator    string
		Left, Right Expression
	}

	// An IndexExpression gets a value from a collection
	IndexExpression struct {
		expr
		Tok               token.Token
		Collection, Index Expression
	}

	// A FunctionCall calls a function
	FunctionCall struct {
		expr
		Tok       token.Token
		Function  Expression
		Arguments []Expression
	}

	// An IfExpression executes Consequence or Alternative based on Condition.
	// Alternative can be nil, in which case the expression will return nil.
	IfExpression struct {
		expr
		Tok                      token.Token
		Condition                Expression
		Consequence, Alternative Expression
	}

	// A MatchBranch is a single condition -> body branch in a match expression.
	MatchBranch struct {
		Condition, Body Expression
	}

	// A MatchExpression executes a different piece of code depending on the
	// input value. If a condition is an identifier who's value is a single
	// underscore, that condition always matches, so always put underscores
	// after everything else.
	MatchExpression struct {
		expr
		Tok      token.Token
		Input    Expression
		Branches []MatchBranch
	}

	// A Type expression defines a new type with the given parameters.
	Type struct {
		expr
		Tok        token.Token
		Parameters []Expression
	}
)
