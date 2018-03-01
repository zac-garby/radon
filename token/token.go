package token

import "fmt"

// A Token is a lexical token representing
// a part of the source code.
type Token struct {
	Type       Type
	Literal    string
	Start, End Position
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"%s `%s` from %s → %s",
		t.Type,
		t.Literal,
		t.Start.String(),
		t.End.String(),
	)
}

// Keywords maps all possible keyword literals to their
// corresponding token types
var Keywords = map[string]Type{
	"return": Return,
	"true":   True,
	"false":  False,
	"nil":    Nil,
	"if":     If,
	"then":   Then,
	"else":   Else,
	"while":  While,
	"for":    For,
	"loop":   Loop,
	"next":   Next,
	"break":  Break,
	"match":  Match,
	"model":  Model,
	"map":    Map,
	"where":  Where,
	"import": Import,
	"do":     Do,
	"in":     In,
}

// IsKeyword checks if a token type is a keyword type.
func IsKeyword(t Type) bool {
	for _, k := range Keywords {
		if t == k {
			return true
		}
	}

	return false
}
