package ast

type stmt struct{}

func (s stmt) Stmt() {}

type (
	// An ExpressionStatement is an expression which can take the place of a
	// statement.
	ExpressionStatement struct {
		stmt
		Expr Expression
	}

	// A Return statement returns a value from a function.
	Return struct {
		stmt
		Value Expression
	}

	// A Next statement jumps to the next iteration of a loop.
	Next struct {
		stmt
	}

	// A Break statement exits out of a loop.
	Break struct {
		stmt
	}

	// A While loop executes Body while Condition is true.
	While struct {
		stmt
		Condition, Body Expression
	}

	// A For loop executes Body for each Var in Collection.
	For struct {
		stmt
		Var, Collection, Body Expression
	}

	// An Import statement imports the file or directory specified into the
	// scope.
	Import struct {
		stmt
		Path string
	}
)
