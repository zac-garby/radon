package lexer

import (
	"strings"

	"github.com/Zac-Garby/radon/token"
)

type transformer func(token.Type, string, string) (token.Type, string, string)
type handler func([]string) (token.Type, string, string)

func lexemeHandler(t token.Type, group int, transformer transformer) handler {
	return func(m []string) (token.Type, string, string) {
		return transformer(t, m[group], m[0])
	}
}

func none(t token.Type, literal, whole string) (token.Type, string, string) {
	return t, literal, whole
}

func stringTransformer(t token.Type, literal, whole string) (token.Type, string, string) {
	escapes := map[string]string{
		`\n`: "\n",
		`"`:  "\"",
		`\a`: "\a",
		`\b`: "\b",
		`\f`: "\f",
		`\r`: "\r",
		`\t`: "\t",
		`\v`: "\v",
	}

	for k, v := range escapes {
		literal = strings.Replace(literal, k, v, -1)
	}

	return t, literal, whole
}

func idTransformer(t token.Type, literal, whole string) (token.Type, string, string) {
	if kwType, ok := token.Keywords[literal]; ok {
		return kwType, literal, whole
	}

	return t, literal, whole
}

type lexicalPair struct {
	regex   string
	handler handler
}

var lexicalDictionary = []lexicalPair{
	// Literals
	{regex: `^\d+(?:\.\d+)?`, handler: lexemeHandler(token.Number, 0, none)},
	{regex: `^"((\\"|[^"])*)"`, handler: lexemeHandler(token.String, 1, stringTransformer)},
	{regex: `^'((\\'|[^'])*)'`, handler: lexemeHandler(token.String, 1, stringTransformer)},
	{regex: "^`([^`]*)`", handler: lexemeHandler(token.String, 1, none)},
	{regex: `^[\p{L}\p{M}_][\p{L}\p{M}\d_!?]*`, handler: lexemeHandler(token.ID, 0, idTransformer)},

	// Punctuation
	{regex: `^->`, handler: lexemeHandler(token.RightArrow, 0, none)},
	{regex: `^=>`, handler: lexemeHandler(token.LambdaArrow, 0, none)},
	{regex: `^\+=`, handler: lexemeHandler(token.PlusEquals, 0, none)},
	{regex: `^\+`, handler: lexemeHandler(token.Plus, 0, none)},
	{regex: `^-=`, handler: lexemeHandler(token.MinusEquals, 0, none)},
	{regex: `^-`, handler: lexemeHandler(token.Minus, 0, none)},
	{regex: `^\^=`, handler: lexemeHandler(token.ExpEquals, 0, none)},
	{regex: `^\^`, handler: lexemeHandler(token.Exp, 0, none)},
	{regex: `^\*=`, handler: lexemeHandler(token.StarEquals, 0, none)},
	{regex: `^\*`, handler: lexemeHandler(token.Star, 0, none)},
	{regex: `^\/\/=`, handler: lexemeHandler(token.FloorDivEquals, 0, none)},
	{regex: `^\/\/`, handler: lexemeHandler(token.FloorDiv, 0, none)},
	{regex: `^\/=`, handler: lexemeHandler(token.SlashEquals, 0, none)},
	{regex: `^\/`, handler: lexemeHandler(token.Slash, 0, none)},
	{regex: `^\(`, handler: lexemeHandler(token.LeftParen, 0, none)},
	{regex: `^\)`, handler: lexemeHandler(token.RightParen, 0, none)},
	{regex: `^<=`, handler: lexemeHandler(token.LessThanEq, 0, none)},
	{regex: `^>=`, handler: lexemeHandler(token.GreaterThanEq, 0, none)},
	{regex: `^<`, handler: lexemeHandler(token.LessThan, 0, none)},
	{regex: `^>`, handler: lexemeHandler(token.GreaterThan, 0, none)},
	{regex: `^{`, handler: lexemeHandler(token.LeftBrace, 0, none)},
	{regex: `^}`, handler: lexemeHandler(token.RightBrace, 0, none)},
	{regex: `^\[`, handler: lexemeHandler(token.LeftSquare, 0, none)},
	{regex: `^]`, handler: lexemeHandler(token.RightSquare, 0, none)},
	{regex: `^;`, handler: lexemeHandler(token.Semi, 0, none)},
	{regex: `^==`, handler: lexemeHandler(token.Equal, 0, none)},
	{regex: `^!=`, handler: lexemeHandler(token.NotEqual, 0, none)},
	{regex: `^\|\|=`, handler: lexemeHandler(token.OrEquals, 0, none)},
	{regex: `^\|\|`, handler: lexemeHandler(token.Or, 0, none)},
	{regex: `^&&=`, handler: lexemeHandler(token.AndEquals, 0, none)},
	{regex: `^&&`, handler: lexemeHandler(token.And, 0, none)},
	{regex: `^\|=`, handler: lexemeHandler(token.BitOrEquals, 0, none)},
	{regex: `^\|`, handler: lexemeHandler(token.BitOr, 0, none)},
	{regex: `^&=`, handler: lexemeHandler(token.BitAndEquals, 0, none)},
	{regex: `^&`, handler: lexemeHandler(token.BitAnd, 0, none)},
	{regex: `^=`, handler: lexemeHandler(token.Assign, 0, none)},
	{regex: `^:=`, handler: lexemeHandler(token.Declare, 0, none)},
	{regex: `^\,`, handler: lexemeHandler(token.Comma, 0, none)},
	{regex: `^:`, handler: lexemeHandler(token.Colon, 0, none)},
	{regex: `^%=`, handler: lexemeHandler(token.ModEquals, 0, none)},
	{regex: `^%`, handler: lexemeHandler(token.Mod, 0, none)},
	{regex: `^\.`, handler: lexemeHandler(token.Dot, 0, none)},
	{regex: `^!`, handler: lexemeHandler(token.Bang, 0, none)},
}
