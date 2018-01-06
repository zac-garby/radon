package ast

import (
	"reflect"

	"github.com/Zac-Garby/lang/token"
)

type stmt struct{}

func (s stmt) Stmt() {}

func (s stmt) Token() token.Token {
	v := reflect.ValueOf(s).FieldByName("Tok")
	if v.IsValid() {
		return v.Interface().(token.Token)
	}

	panic("could not get token")
}

type (
	// ExpressionStatement is an expression which acts a statement
	ExpressionStatement struct {
		stmt
		Tok  token.Token
		Expr Expression
	}

	// ReturnStatement returns an expression from a BlockStatement
	ReturnStatement struct {
		stmt
		Tok   token.Token
		Value Expression
	}

	// NextStatement goes to the next iteration of a loop
	NextStatement struct {
		stmt
		Tok token.Token
	}

	// BreakStatement breaks a loop
	BreakStatement struct {
		stmt
		Tok token.Token
	}

	// Loop executes Body indefinitely
	Loop struct {
		stmt
		Tok  token.Token
		Body Expression
	}

	// WhileLoop executes Body while Condition holds true
	WhileLoop struct {
		stmt
		Tok       token.Token
		Condition Expression
		Body      Expression
	}

	// ForLoop executes Body for each Var in Collection.
	ForLoop struct {
		stmt
		Tok token.Token

		Init      Expression
		Condition Expression
		Increment Expression
		Body      Expression
	}

	// Import imports the file or directory specified into the
	// scope.
	Import struct {
		stmt
		Tok token.Token

		Path string
	}
)
