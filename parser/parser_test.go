package parser_test

import (
	"fmt"
	"testing"

	"github.com/Zac-Garby/radon/lexer"
	. "github.com/Zac-Garby/radon/parser"
)

func TestNoErrors(t *testing.T) {
	tests := []string{
		"hello",
		"100",
		"2.3",
		"true",
		"false",
		"nil",
		`"hello"`,
		"'hello'",
		"`hello`",

		"(5)",

		"[]",
		"[1, 2, 3]",
		"[1]",
		"[1,]",
		"[1, 2, 3,]",

		"{}",
		"{a: b}",
		"{a: b,}",
		"{a: b, c: d}",
		"{a: b, c: d,}",

		"do 1; 2; 3 end",
		"do end",
		"do do 5 end end",
	}

	for i, test := range tests {
		var (
			l      = lexer.Lexer(test, fmt.Sprintf("test %d: %s", i, test))
			p      = New(l)
			_, err = p.Parse()
		)

		if err != nil {
			fmt.Println(err.Error())
			t.Fail()
		}
	}
}

func TestErrors(t *testing.T) {
	tests := map[string]string{
		")": "unexpected token: right-paren",
		"$": "illegal token encountered. literal: `$`",

		"(":     "unexpected token: semi",
		"[":     "unexpected token: semi",
		"[1,":   "unexpected token: semi",
		"{":     "unexpected token: semi",
		"{1:":   "unexpected token: semi",
		"{1:2,": "unexpected token: semi",
	}

	for test, expectedMessage := range tests {
		var (
			l      = lexer.Lexer(test, fmt.Sprintf("test %s", test))
			p      = New(l)
			_, err = p.Parse()
		)

		if err == nil {
			t.Fail()
		}

		message := err.(*Error).Message

		if message != expectedMessage {
			fmt.Printf("expected %s, got %s\n", expectedMessage, message)
			t.Fail()
		}
	}
}
