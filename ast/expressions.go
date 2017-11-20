package ast

import "github.com/Zac-Garby/pluto/token"

type (
	// Identifier is an identifier
	Identifier struct {
		Tok   token.Token
		Value string
	}

	// Number is a number literal
	Number struct {
		Tok   token.Token
		Value float64
	}

	// Boolean is a boolean literal
	Boolean struct {
		Tok   token.Token
		Value bool
	}

	// String is a string literal
	String struct {
		Tok   token.Token
		Value string
	}

	// Tuple is a tuple literal
	Tuple struct {
		Tok   token.Token
		Value []Expression
	}

	// Array is an array literal
	Array struct {
		Tok      token.Token
		Elements []Expression
	}

	// Map is a map literal
	Map struct {
		Tok   token.Token
		Pairs map[Expression]Expression
	}

	// A Block combines multiple statements into an expression
	Block struct {
		Tok        token.Token
		Statements []Statement
	}

	// Nil is the nil literal
	Nil struct {
		Tok token.Token
	}

	// PrefixExpression is a prefix operator expression
	PrefixExpression struct {
		Tok      token.Token
		Operator string
		Right    Expression
	}

	// An InfixExpression is an infix operator expression
	InfixExpression struct {
		Tok         token.Token
		Operator    string
		Left, Right Expression
	}

	// A DotExpression gets a value from a container
	DotExpression struct {
		Tok         token.Token
		Left, Right Expression
	}

	// An IndexExpression gets a value from a collection
	IndexExpression struct {
		Tok               token.Token
		Collection, Index Expression
	}

	// A FunctionCall calls a function
	FunctionCall struct {
		Tok       token.Token
		Function  Expression
		Arguments []Expression
	}

	// An IfExpression executes Consequence or Alternative based on Condition.
	// Alternative can be nil, in which case the expression will return nil.
	IfExpression struct {
		Tok                      token.Token
		Condition                Expression
		Consequence, Alternative Statement
	}
)

// Expr tells the compiler this node is an expression
func (n Identifier) Expr() {}

// Token returns the node's token
func (n Identifier) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Number) Expr() {}

// Token returns the node's token
func (n Number) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Boolean) Expr() {}

// Token returns the node's token
func (n Boolean) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n String) Expr() {}

// Token returns the node's token
func (n String) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Tuple) Expr() {}

// Token returns the node's token
func (n Tuple) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Array) Expr() {}

// Token returns the node's token
func (n Array) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Map) Expr() {}

// Token returns the node's token
func (n Map) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Nil) Expr() {}

// Token returns the node's token
func (n Nil) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n PrefixExpression) Expr() {}

// Token returns the node's token
func (n PrefixExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n InfixExpression) Expr() {}

// Token returns the node's token
func (n InfixExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n DotExpression) Expr() {}

// Token returns the node's token
func (n DotExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n IndexExpression) Expr() {}

// Token returns the node's token
func (n IndexExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n FunctionCall) Expr() {}

// Token returns the node's token
func (n FunctionCall) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n IfExpression) Expr() {}

// Token returns the node's token
func (n IfExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Block) Expr() {}

// Token returns the node's token
func (n Block) Token() token.Token { return n.Tok }
