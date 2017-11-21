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

	panic("could not get token")
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

	// Array is an array literal
	Array struct {
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

	// Nil is the nil literal
	Nil struct {
		expr
		Tok token.Token
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

	// A DotExpression gets a value from a container
	DotExpression struct {
		expr
		Tok         token.Token
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
		Consequence, Alternative Statement
	}

	// A MatchExpression executes a different piece of code depending on the
	// input value. If a condition is an identifier who's value is a single
	// underscore, that condition always matches, so always put underscores
	// after everything else.
	MatchExpression struct {
		expr
		Tok      token.Token
		Input    Expression
		Branches map[Expression]Expression
	}

	// A Type expression defines a new type with the given parameters.
	Type struct {
		expr
		Tok        token.Token
		Parameters []string
	}
)
