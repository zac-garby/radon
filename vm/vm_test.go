package vm_test

import (
	"io/ioutil"
	"testing"

	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/compiler"
	"github.com/Zac-Garby/radon/parser"
	. "github.com/Zac-Garby/radon/vm"
)

func TestEvaluation(t *testing.T) {
	tests := map[string]string{
		"5":      "5",
		"5 + 1":  "6",
		"-10":    "-10",
		"5 - 1":  "4",
		"5 * 5":  "25",
		"5 / 5":  "1",
		"5 ^ 2":  "25",
		"5 // 2": "2",
		"5 % 2":  "1",
		"5 | 3":  "7",
		"5 & 7":  "5",
		"5 == 5": "true",
		"5 != 5": "false",
		"5 > 5":  "false",
		"5 < 5":  "false",
		"5 >= 5": "true",
		"5 <= 5": "true",

		"x": "",

		`"Hello"`:                  "Hello",
		"true":                     "true",
		"nil":                      "nil",
		"a = 5; a":                 "5",
		"[1, 2, 3]":                "[1, 2, 3]",
		"(1, 2, 3)":                "(1, 2, 3)",
		"map['a': 5 + 0]":          `[a: 5]`,
		"if true then 5":           "5",
		"if false then 5":          "nil",
		"if true then 10 else 3":   "10",
		"if 5 > 10 then 10 else 3": "3",

		`if 1 > 2 then
			'foo'
		else if 2 > 2 then
			'bar'
		else
			'baz'`: "baz",

		"[1, 2, 3][1]":                           "2",
		"map['a': 2, 'b': 4].b":                  "4",
		"a = [1, 2, 3]; a[1] = 10; a[1]":         "10",
		"a = map['a': 2, 'b': 4]; a.a = 10; a.a": "10",

		"x := 5; x + 2":                      "7",
		"f() = 10; f()":                      "10",
		"f(); f() = 5":                       "5",
		"f() = { return 5; return 3; }; f()": "5",

		"for i = 0; i < 10; i = i + 1 do i":            "10",
		"loop { if i > 10 then { break }; i = i + 1 }": "",
		"a = 0; while (a < 10) { next; a }":            "",
		`for i = 0; i < 20; i = i + 1 {
			if i > 10 then {
				break
			} else i = i + 1
		}`: "12",

		`f(x, y, z) = x * y - z;
		f(1, 2, 3)`: "5",

		`f(x, y) = {
			z = x + y
			z - 1
		}
		f(1, 2)`: "2",

		`a = 10
		f() = a
		f()`: "10",

		`match 5 where
			| 0 -> "foo",
			| 1 -> "bar",
			| 5 -> "baz"
		`: "baz",

		`print("print :)")`: "()",
		`echo("echo'd")`:    "()",
		"len([1, 2, 3])":    "3",

		`vec = model(x, y)
		a = vec(2, 5)
		a.x`: "2",

		`vec = model(x, y)
		f(s) = {
			self.x = self.x * s
			self.y = self.y * s
			self
		}
		vec.scale = f
		a = vec(3, 2)
		a.scale(2)
		a.x + a.y`: "10",

		"!true":                 "false",
		"2 * [1, 2, 3]":         "[1, 2, 3, 1, 2, 3]",
		"[1, 2, 3] * 2":         "[1, 2, 3, 1, 2, 3]",
		"[1, 2, 3] | [2, 3, 4]": "[1, 2, 3, 4]",
		"[1, 2, 3] & [2, 3, 4]": "[2, 3]",
		"[1, 2, 3] + [2, 3, 4]": "[1, 2, 3, 2, 3, 4]",
		"[1, 2, 3] - [1, 3]":    "[2]",
	}

	for in, out := range tests {
		var (
			vm    = New()
			store = NewStore()
			cmp   = compiler.New()
			parse = parser.New(in, "test")
		)

		// Parse the input -- assume it parses correctly
		prog := parse.Parse()

		// Compile the program -- assume it compiles correctly
		if err := cmp.Compile(prog); err != nil {
			t.Error(err.Error())
			return
		}

		//
		code, err := bytecode.Read(cmp.Bytes)
		if err != nil {
			t.Error(err.Error())
			return
		}

		store.Names = cmp.Names
		vm.Out = ioutil.Discard
		vm.Run(code, store, cmp.Constants)

		result, err := vm.ExtractValue()
		if err != nil {
			t.Error(err.Error())
			return
		}

		if out == "" {
			if result != nil {
				t.Errorf("(%s) expected no output, got: %s", in, result.String())
			}
		} else if result != nil {
			resString := result.String()
			if resString != out {
				t.Errorf("(%s) expected: %s\ngot: %s\n", in, out, resString)
				return
			}
		} else {
			t.Errorf("(%s) expected: %s\ngot nothing", in, out)
		}
	}
}
