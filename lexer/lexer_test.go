package lexer_test

import (
	"testing"

	"github.com/Zac-Garby/radon/lexer"
	. "github.com/Zac-Garby/radon/token"
)

func TestLexing(t *testing.T) {
	input := `123.50 1 2#;
	"hello" "hello world" foo bar # where a semicolon should be inserted
	+-*^/ //%()<><=>={}[];==!=||&& # a comment after a line
	| & = := , -> : . ! += -=#this comment touches the token
	*= ^= /= //= %= ||= &&= |= &= #nospacesnospacesnospaces!!

	# keywords now :)
	return true false nil if then else
	while for loop next break match model in

	$
	`

	expected := []Type{
		Number, Number, Number, Semi,
		String, String, ID, ID, Semi,
		Plus, Minus, Star, Exp, Slash, FloorDiv, Mod,
		LeftParen, RightParen, LessThan, GreaterThan,
		LessThanEq, GreaterThanEq, LeftBrace, RightBrace,
		LeftSquare, RightSquare, Semi, Equal, NotEqual,
		Or, And, BitOr, BitAnd, Assign, Declare,
		Comma, RightArrow, Colon, Dot, Bang,
		PlusEquals, MinusEquals, StarEquals, ExpEquals,
		SlashEquals, FloorDivEquals, ModEquals, OrEquals,
		AndEquals, BitOrEquals, BitAndEquals,

		Return, True, False, Nil, If, Then, Else, While,
		For, Loop, Next, Break, Match, Model, In,

		Illegal,
	}

	next := lexer.Lexer(input, "test")

	i := 0
	for tok := next(); tok.Type != EOF; tok = next() {
		exp := expected[i]
		i++

		if tok.Type != exp {
			t.Errorf("(%v) expected %s, got %s\n", i, exp, tok.Type)
		}

		if tok.Start.Filename != "test" || tok.End.Filename != "test" {
			t.Errorf("(%v) reported wrong file name", i)
		}
	}
}
