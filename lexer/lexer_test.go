package lexer_test

import (
	"testing"

	"github.com/Zac-Garby/lang/lexer"
	. "github.com/Zac-Garby/lang/token"
)

func TestLexing(t *testing.T) {
	input := `123.50 1 2#;
	"hello" "hello world" foo bar # where a semicolon should be inserted
	+-*^/ //%()<><=>={}[];==!=||&& # a comment after a line
	| & = := , -> <- : . ! += -=#this comment touches the token
	*= ^= /= //= %= ||= &&= |= &= #nospacesnospacesnospaces!!
	
	# keywords now :)
	return true false nil if then else
	while for loop next break match type
	`

	expected := []Type{
		Number, Number, Number, Semi,
		String, String, ID, ID, Semi,
		Plus, Minus, Star, Exp, Slash, FloorDiv, Mod,
		LeftParen, RightParen, LessThan, GreaterThan,
		LessThanEq, GreaterThanEq, LeftBrace, RightBrace,
		LeftSquare, RightSquare, Semi, Equal, NotEqual,
		Or, And, BitOr, BitAnd, Assign, Declare,
		Comma, RightArrow, LeftArrow, Colon, Dot, Bang,
		PlusEquals, MinusEquals, StarEquals, ExpEquals,
		SlashEquals, FloorDivEquals, ModEquals, OrEquals,
		AndEquals, BitOrEquals, BitAndEquals,

		Return, True, False, Nil, If, Then, Else, While,
		For, Loop, Next, Break, Match, TypeK,
	}

	next := lexer.Lexer(input, "<repl>")

	i := 0
	for tok := next(); tok.Type != EOF; tok = next() {
		exp := expected[i]
		i++

		if tok.Type != exp {
			t.Logf("(%v) expected %s, got %s\n", i, exp, tok.Type)
			t.Fail()
		}
	}
}
