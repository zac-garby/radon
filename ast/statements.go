package ast

import "github.com/Zac-Garby/pluto/token"

type (
	// ExpressionStatement is an expression which acts a statement
	ExpressionStatement struct {
		Tok  token.Token
		Expr Expression
	}

	// ReturnStatement returns an expression from a BlockStatement
	ReturnStatement struct {
		Tok   token.Token
		Value Expression
	}

	// NextStatement goes to the next iteration of a loop
	NextStatement struct {
		Tok token.Token
	}

	// BreakStatement breaks a loop
	BreakStatement struct {
		Tok token.Token
	}

	// Loop executes Body indefinitely
	Loop struct {
		Tok  token.Token
		Body Expression
	}

	// WhileLoop executes Body while Condition holds true
	WhileLoop struct {
		Tok       token.Token
		Condition Expression
		Body      Expression
	}

	// ForLoop executes Body while Condition holds true, evaluating Increment
	// each iteration and evaluating Init at the start
	ForLoop struct {
		Tok token.Token

		// for (Init; Condition; Increment) { Body }
		Init      Expression
		Condition Expression
		Increment Expression
		Body      Expression
	}
)

// Stmt tells the compiler this node is a statement
func (n ExpressionStatement) Stmt() {}

//Token returns this node's token
func (n ExpressionStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n ReturnStatement) Stmt() {}

//Token returns this node's token
func (n ReturnStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n NextStatement) Stmt() {}

//Token returns this node's token
func (n NextStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n BreakStatement) Stmt() {}

//Token returns this node's token
func (n BreakStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n Loop) Stmt() {}

//Token returns this node's token
func (n Loop) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n WhileLoop) Stmt() {}

//Token returns this node's token
func (n WhileLoop) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n ForLoop) Stmt() {}

//Token returns this node's token
func (n ForLoop) Token() token.Token { return n.Tok }
