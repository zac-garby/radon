package vm

import (
	"errors"
	"fmt"

	"github.com/Zac-Garby/lang/object"
)

// Builtins are the builtin functions, such as print
// and len.
var Builtins map[string]*object.Function

func init() {
	// In the builtin functions, the Parameters field contains
	// arbitrary strings just so the VM knows how many arguments
	// the function accepts.

	Builtins = map[string]*object.Function{
		"print": &object.Function{
			Parameters: []string{"arg"},

			OnCall: func(f *object.Function, args map[string]object.Object) (object.Object, error) {
				fmt.Println(args["arg"])

				return object.EmptyObj, nil
			},
		},

		"echo": &object.Function{
			Parameters: []string{"arg"},

			OnCall: func(f *object.Function, args map[string]object.Object) (object.Object, error) {
				fmt.Print(args["arg"])

				return object.EmptyObj, nil
			},
		},

		"len": &object.Function{
			Parameters: []string{"arg"},

			OnCall: func(f *object.Function, args map[string]object.Object) (object.Object, error) {
				total := 0

				for _, arg := range args {
					if col, ok := arg.(object.Collection); ok {
						total += len(col.Elements())
					} else {
						return nil, errors.New("wrong_type: argument %s isn't a collection and therefore doesn't have a length")
					}
				}

				return &object.Number{Value: float64(total)}, nil
			},
		},
	}
}
