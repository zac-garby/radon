package vm_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/compiler"
	"github.com/Zac-Garby/lang/parser"
	. "github.com/Zac-Garby/lang/vm"
)

func TestErrors(t *testing.T) {
	tests := map[string]struct {
		t   ErrType
		msg string
	}{
		"len(5)": {
			t:   ErrWrongType,
			msg: "cannot get the length of type number",
		},

		"[1, 2, 3]['x']": {
			t:   ErrWrongType,
			msg: "non-numeric type string used to index a collection",
		},

		"5[1]": {
			t:   ErrWrongType,
			msg: "cannot index type number",
		},

		"[]['y'] = 'z'": {
			t:   ErrWrongType,
			msg: "non-numeric type string used to index a collection",
		},

		"5['y'] = 'z'": {
			t:   ErrWrongType,
			msg: "cannot index type number",
		},

		"-'a'": {
			t:   ErrWrongType,
			msg: "prefix r-value of invalid type",
		},

		"5 * true": {
			t:   ErrWrongType,
			msg: "infix r-value of invalid type when l-value is number",
		},

		"[1, 2, 3] + true": {
			t:   ErrWrongType,
			msg: "infix r-value of invalid type when l-value is a collection",
		},

		"true + false": {
			t:   ErrWrongType,
			msg: "infix l-value of invalid type",
		},

		"[] + 5": {
			t:   ErrNoOp,
			msg: "operator ADD not supported for a collection and a number",
		},

		"[] / []": {
			t:   ErrNoOp,
			msg: "operator DIVIDE not supported for two collections",
		},

		"1 > []": {
			t:   ErrWrongType,
			msg: "non-numeric value in numeric binary expression",
		},

		"[] > 1": {
			t:   ErrWrongType,
			msg: "non-numeric value in numeric binary expression",
		},

		"5()": {
			t:   ErrWrongType,
			msg: "can only call functions and models",
		},

		"f() = 5; f(1)": {
			t:   ErrArgument,
			msg: "wrong amount of arguments supplied to the function. expected 0",
		},

		"f() = { break }; f()": {
			t:   ErrSyntax,
			msg: "break statement found outside loop",
		},

		"m = model(x, y); m(1)": {
			t:   ErrArgument,
			msg: "wrong amount of arguments supplied to the function. expected 2",
		},

		"break": {
			t:   ErrSyntax,
			msg: "break statement found outside loop",
		},

		"next": {
			t:   ErrSyntax,
			msg: "next statement found outside loop",
		},
	}

	for in, expected := range tests {
		var (
			vm    = New()
			store = NewStore()
			cmp   = compiler.New()
			parse = parser.New(in, "test")
			code  bytecode.Code
			err   error
		)

		// Parse the input -- assume it parses correctly
		prog := parse.Parse()

		// Compile the program -- assume it compiles correctly
		if err = cmp.Compile(prog); err != nil {
			goto thrown
		}

		code, err = bytecode.Read(cmp.Bytes)
		if err != nil {
			goto thrown
		}

		store.Names = cmp.Names
		vm.Out = ioutil.Discard
		vm.Run(code, store, cmp.Constants)

		err = vm.Error()

	thrown:
		if err == nil {
			t.Errorf("no error thrown - expected '%s: %s'", expected.t, expected.msg)
			return
		}

		errString := err.Error()

		if !strings.HasPrefix(errString, string(expected.t)) {
			t.Errorf("expected prefix of '%s' on error '%s'", expected.t, errString)
		}
	}
}
