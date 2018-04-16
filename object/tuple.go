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

// Subscript subscrips an Object, e.g. foo[bar], or returns false if it can't be
// done.
func (t *Tuple) Subscript(index Object) (Object, bool) {
	num, ok := index.(*Number)
	if !ok {
		return nil, false
	}

	i := int(num.Value)
	if i < 0 || i >= len(t.Value) {
		return nil, false
	}

	return t.Value[i], true
}

// SetSubscript sets the value of a subscript of an Object, e.g. foo[bar] = baz.
// Returns false if it can't be done.
func (t *Tuple) SetSubscript(index Object, to Object) bool {
	num, ok := index.(*Number)
	if !ok {
		return false
	}

	i := int(num.Value)
	if i < 0 || i >= len(t.Value) {
		return false
	}

	t.Value[i] = to

	return true
}
