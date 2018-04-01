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

func (t *Tuple) Type() Type {
	return TupleType
}

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

func (t *Tuple) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: append(t.Value, right),
		}, true
	}

	return nil, false
}

func (t *Tuple) Items() ([]Object, bool) {
	return t.Value, true
}

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
