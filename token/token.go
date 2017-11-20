package token

import "fmt"

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
	"type":   TypeK,
}

func IsKeyword(t Type) bool {
	for _, k := range Keywords {
		if t == k {
			return true
		}
	}

	return false
}
