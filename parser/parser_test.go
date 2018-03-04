package parser_test

import (
	"fmt"
	"testing"

	"github.com/Zac-Garby/radon/ast"
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
		"`hello`;;",

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

		"-5",
		"+do 1; 2; 3 end",
		"!true",

		"if true then 1",
		"if true then 1 else 2",
		`if true do
             1
             2
         end`,
		`if true do
             1
             2
         end else do
             3
             4
         end
        `,

		"match n where",
		`match n where
             | a -> b`,
		`match n where
             | a -> b,
             | b -> c`,
		`match n where
             | a -> b,
             | b -> c,
             | _ -> d`,

		"model ()",
		"model (a,)",
		"model (a, b)",
		"model (a, b) | parent",
		"model (a, b) | parent ('hello', 5, a)",

		"=> 10",

		"1 + 1",
		"1 + 2 * 3",
		"a, b",

		"a[b]",
		"a[b, c]",

		"return",
		"return 5",
		"return 1, 2, 3",

		"next",
		"break",

		"while true, x",
		"while true do a; b; c end",

		"for a in b, c",
		"for a in b do c; d; e end",

		"import 'foo'",
	}

	for i, test := range tests {
		_, err := parse(test, fmt.Sprintf("test %d: %s", i, test))

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

		"(":     "unexpected end of line",
		"[":     "unexpected end of line",
		"[1,":   "unexpected end of line",
		"{":     "unexpected end of line",
		"{1:":   "unexpected end of line",
		"{1:2,": "unexpected end of line",

		"if true": "unexpected end of line, wanted 'then'",

		"match x":           "unexpected end of line, wanted 'where'",
		"match x where | a": "unexpected end of line, wanted 'right-arrow'",

		"model":           "unexpected end of line, wanted 'left-paren'",
		"model (5)":       "expected 'identifier' but got 'number'",
		"model (a, true)": "expected 'identifier' but got 'true'",
		"model (x":        "unexpected end of line, wanted 'right-paren'",
		"model (a, b, a)": "identical parameter a not allowed",

		"=>": "unexpected end of line",

		"while true 5": "expected 'comma' but got 'number'",
		"for a do b":   "expected 'in' but got 'do'",
		"for a in b c": "expected 'comma' but got 'identifier'",

		"import":   "unexpected end of line, wanted 'string'",
		"import 5": "expected 'string' but got 'number'",

		"a[b": "unexpected end of line, wanted 'right-square'",
	}

	for test, expectedMessage := range tests {
		_, err := parse(test, fmt.Sprintf("test %s", test))

		if err == nil {
			t.Fail()
		}

		message := err.(*Error).Message

		if message != expectedMessage {
			fmt.Printf("expected `%s`, got `%s`\n", expectedMessage, message)
			t.Fail()
		}
	}
}

func TestPrecedence(t *testing.T) {
	tests := map[string]string{
		"1 + 2 * 3": "1 + (2 * 3)",
		"1 * 2 + 3": "(1 * 2) + 3",
		"1 + 2, 3":  "(1 + 2), 3",
		"1, 2 + 3":  "1, (2 + 3)",
	}

	for test, expected := range tests {
		testAST, err := parse(test, "test")
		if err != nil {
			t.Fail()
		}

		expectedAST, err := parse(expected, "test")
		if err != nil {
			t.Fail()
		}

		if testAST.Tree() != expectedAST.Tree() {
			fmt.Println(test, "doesn't parse the same as", expected)
			t.Fail()
		}
	}
}

func parse(str, file string) (*ast.Program, error) {
	var (
		l = lexer.Lexer(str, file)
		p = New(l)
	)

	return p.Parse()
}
