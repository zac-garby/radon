package parser_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Zac-Garby/radon/parser"
)

// Test if errors are thrown where they should be. Assume that any error is the
// correct one.
func TestErrors(t *testing.T) {
	tests := []string{
		"a a",
		"(5 + 3",
		"a[b",

		"[a",
		"map[",
		"map",

		"if true",
		"match a",
		"match a where | b = c",

		"model",
		"model(5)",
		"model(a, 5)",
		"model(a",
		"model(a, a)",

		"while a",
		"while (a",

		"for a",
		"for a",
		"for a; b",
		"for a; b; c",

		"import",
		"import 5",
	}

	for _, test := range tests {
		p := parser.New(test, "test")
		p.Parse()

		if len(p.Errors) == 0 {
			fmt.Println("expected error on input:", test)
			p.PrintErrors(os.Stdout)

			t.Fail()
		}
	}
}

// Assume that, if a string parses with no errors, it has parsed successfully.
func TestNoErrors(t *testing.T) {
	tests := []string{
		"a",
		"hellΩ",

		"5",
		"10",
		"5.33",

		"true",
		"false",

		"nil",

		`'hello "world"'`,
		`"hello 'world'"`,
		"`'hello' \"world\"`",

		"()",
		"(1,)",
		"(1, 2)",
		"(1, 2,)",

		"[]",
		"[1]",
		"[1, 2]",
		"[1, 2,]",

		";",
		"break",
		"next",
		"loop a",
		"while a < b do c",
		"while a < b { c }",
		"for a; b; c do d",
		"for a; b; c { d }",
		"import 'path/to/file'",
		"return a",

		"map[]",
		"map[1:2]",
		"map[1:2, 'a': 'A']",
		"map[a: b,]",

		"{}",
		"{ return }",
		"{ a + b }",
		"{{}}",

		"-10",
		"+10",
		"!true",

		"a = b",
		"a := b",
		"a &&= b",
		"a &= b",
		"a |= b",
		"a ^= b",
		"a //= b",
		"a -= b",
		"a %= b",
		"a ||= b",
		"a += b",
		"a /= b",
		"a *= b",
		"a || b",
		"a && b",
		"a | b",
		"a & b",
		"a == b",
		"a != b",
		"a < b",
		"a > b",
		"a <= b",
		"a >= b",
		"a + b",
		"a - b",
		"a * b",
		"a / b",
		"a % b",
		"a ^ b",
		"a // b",
		"a.b",
		"a * b + c",
		"a + b * c",
		"a * (b + c)",

		"a[b]",

		"a()",
		"a(b)",
		"a(b, c)",
		"a(b, c,)",

		"if a then b",
		"if a then b else c",
		"if a then b else if c then d else e",

		"match a where",
		"match a where | b -> c",
		"match a where | b -> c, | d -> e",
		"match a where | _ -> b",

		"model()",
		"model(x)",
		"model(x, y)",
		"model(x, y,)",

		"import 'hello'",
	}

	for _, test := range tests {
		p := parser.New(test, "test")
		p.Parse()

		if len(p.Errors) > 0 {
			fmt.Println("failed on input:", test)
			p.PrintErrors(os.Stdout)

			t.Fail()
		}
	}
}
