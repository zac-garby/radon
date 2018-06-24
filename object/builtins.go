package object

import (
	"fmt"
)

// Builtins contains every builtin.
var Builtins = make(map[string]*Builtin)

func init() {
	Builtins["print"] = &Builtin{
		Name: "print",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			for i, arg := range args {
				if str, ok := arg.(*String); ok {
					fmt.Print(str.Value)
				} else {
					fmt.Print(arg)
				}

				if i+1 < len(args) {
					fmt.Print(" ")
				}
			}

			fmt.Print("\n")

			return &Nil{}, "", ""
		},
	}

	Builtins["put"] = &Builtin{
		Name: "put",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			for i, arg := range args {
				if str, ok := arg.(*String); ok {
					fmt.Print(str.Value)
				} else {
					fmt.Print(arg)
				}

				if i+1 < len(args) {
					fmt.Print(" ")
				}
			}

			return &Nil{}, "", ""
		},
	}

	Builtins["len"] = &Builtin{
		Name: "len",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			total := 0

			for _, arg := range args {
				items, ok := arg.Items()
				if !ok {
					return nil, "Type", fmt.Sprintf("cannot get the length of '%s'", arg.String())
				}

				total += len(items)
			}

			return &Number{Value: float64(total)}, "", ""
		},
	}

	Builtins["tup"] = &Builtin{
		Name: "tup",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			if len(args) == 1 {
				if items, ok := args[0].Items(); ok {
					return &Tuple{Value: items}, "", ""
				}
			}

			return &Tuple{Value: args}, "", ""
		},
	}

	Builtins["list"] = &Builtin{
		Name: "list",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			if len(args) == 1 {
				if items, ok := args[0].Items(); ok {
					return &List{Value: items}, "", ""
				}
			}

			return &List{Value: args}, "", ""
		},
	}

	Builtins["items"] = Builtins["list"]

	Builtins["str"] = &Builtin{
		Name: "str",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			if len(args) != 1 {
				return nil, "Argument", "expected exactly one argument to str(...)"
			}

			return &String{Value: args[0].String()}, "", ""
		},
	}

	Builtins["id"] = &Builtin{
		Name: "id",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			if len(args) != 1 {
				return nil, "Argument", "expected exactly one argument to id(...)"
			}

			return args[0], "", ""
		},
	}

	Builtins["type"] = &Builtin{
		Name: "type",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			if len(args) != 1 {
				return nil, "Argument", "expected exactly one argument to type(...)"
			}

			return &String{Value: string(args[0].Type())}, "", ""
		},
	}

	Builtins["prefix"] = &Builtin{
		Name: "prefix",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			if len(args) != 2 {
				return nil, "Argument", "expected exactly two arguments to prefix(...)"
			}

			arg := args[1]

			str, ok := args[0].(*String)
			if !ok {
				return nil, "Type", "the first argument to prefix(...) should be a string"
			}

			op := str.Value

			res, ok := arg.Prefix(op)
			if !ok {
				return nil, "Type", fmt.Sprintf("could not apply prefix operator %s to %s", op, arg)
			}

			return res, "", ""
		},
	}

	Builtins["infix"] = &Builtin{
		Name: "infix",
		Fn: func(args ...Object) (result Object, errorType string, errorMessage string) {
			if len(args) != 3 {
				return nil, "Argument", "expected exactly three arguments to infix(...)"
			}

			left := args[0]
			right := args[2]

			str, ok := args[1].(*String)
			if !ok {
				return nil, "Type", "the second argument to infix(...) should be a string"
			}

			op := str.Value

			res, ok := left.Infix(op, right)
			if !ok {
				return nil, "Type", fmt.Sprintf("could not apply infix operator %s between %s and %s", op, left, right)
			}

			return res, "", ""
		},
	}
}
