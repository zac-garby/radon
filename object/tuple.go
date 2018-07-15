package object

import (
	"fmt"
	"strings"
)

// A Tuple is an statically-sized collection of items.
type Tuple struct {
	defaults
	Value []Object
}

func (t *Tuple) String() string {
	var strs []string

	for _, item := range t.Value {
		strs = append(strs, item.String())
	}

	return fmt.Sprintf("(%s)", strings.Join(strs, ", "))
}

// Type returns the type of an Object.
func (t *Tuple) Type() Type {
	return TupleType
}

// Equals checks whether or not two objects are equal to each other.
func (t *Tuple) Equals(other Object) bool {
	switch o := other.(type) {
	case *Tuple:
		left, right := t.Value, o.Value
		if len(left) != len(right) {
			return false
		}

		for i, item := range left {
			if !item.Equals(right[i]) {
				return false
			}
		}

		return true

	default:
		return false
	}
}

// Prefix applies a prefix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (t *Tuple) Prefix(op string) (Object, bool) {
	if op == "," {
		return &Tuple{Value: []Object{t}}, true
	}

	return nil, false
}

// Infix applies a infix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (t *Tuple) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: append(t.Value, right),
		}, true
	}

	return nil, false
}

// Items returns a slice containing all objects in an Object, or false otherwise.
func (t *Tuple) Items() ([]Object, bool) {
	return t.Value, true
}
