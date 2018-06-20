package object

import (
	"fmt"
)

// A Builtin is a function which has been written in Go but is callable from
// a Radon program.
type Builtin struct {
	defaults
	Name string
	Fn   func(args ...Object) (result Object, errorType string, errorMessage string)
}

func (b *Builtin) String() string {
	return fmt.Sprintf("<builtin %s>", b.Name)
}

// Type returns the type of an Object.
func (b *Builtin) Type() Type {
	return BuiltinType
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

// Prefix applies a prefix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (b *Builtin) Prefix(op string) (Object, bool) {
	if op == "," {
		return &Tuple{Value: []Object{b}}, true
	}

	return nil, false
}

// Infix applies a infix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (b *Builtin) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: []Object{b, right},
		}, true
	}

	return nil, false
}
