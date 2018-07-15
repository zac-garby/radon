package object

import (
	"fmt"
	"strings"
)

// A List is a dynamic mutable linked list.
type List struct {
	defaults
	Value []Object
}

func (l *List) String() string {
	var strs []string

	for _, item := range l.Value {
		strs = append(strs, item.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

// Type returns the type of an Object.
func (l *List) Type() Type {
	return ListType
}

// Equals checks whether or not two objects are equal to each other.
func (l *List) Equals(other Object) bool {
	switch o := other.(type) {
	case *List:
		left, right := l.Value, o.Value
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
func (l *List) Prefix(op string) (Object, bool) {
	if op == "," {
		return &Tuple{Value: []Object{l}}, true
	}

	return nil, false
}

// Infix applies a infix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (l *List) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: []Object{l, right},
		}, true
	}

	if op == "+" {
		other, ok := right.(*List)
		if !ok {
			return nil, false
		}

		return &List{Value: append(l.Value, other.Value...)}, true
	}

	return nil, false
}

// Items returns a slice containing all objects in an Object, or false otherwise.
func (l *List) Items() ([]Object, bool) {
	return l.Value, true
}

// Subscript subscrips an Object, e.g. foo[bar], or returns false if it can't be
// done.
func (l *List) Subscript(index Object) (Object, bool) {
	num, ok := index.(*Number)
	if !ok {
		return nil, false
	}

	i := int(num.Value)
	if i < 0 || i >= len(l.Value) {
		return nil, false
	}

	return l.Value[i], true
}

// SetSubscript sets the value of a subscript of an Object, e.g. foo[bar] = baz.
// Returns false if it can't be done.
func (l *List) SetSubscript(index Object, to Object) bool {
	num, ok := index.(*Number)
	if !ok {
		return false
	}

	i := int(num.Value)
	if i < 0 || i >= len(l.Value) {
		return false
	}

	l.Value[i] = to

	return true
}

// Iter creates an iterable from an Object.
func (l *List) Iter() (Iterable, bool) {
	return &ListIterable{
		List:  l,
		Index: 0,
	}, true
}
