package object

import (
	"fmt"
)

// A Builtin is a function which has been written in Go but is callable from
// a Radon program.
type Builtin struct {
	defaults
	Name string
	Fn   func(args ...Object) Object
}

func (b *Builtin) String() string {
	return fmt.Sprintf("<builtin %s>", b.Name)
}

// Equals checks if two builtins are equal to each-other. Two builtins are equal
// if their names are the same.
func (b *Builtin) Equals(other Object) bool {
	switch o := other.(type) {
	case *Builtin:
		return b.Name == o.Name

	default:
		return false
	}
}
