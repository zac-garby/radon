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
				fmt.Print(arg)
				if i+1 < len(args) {
					fmt.Print(" ")
				}
			}

			fmt.Print("\n")

			return &Nil{}, "", ""
		},
	}
}
