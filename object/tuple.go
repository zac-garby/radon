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
	if op == "[]" || op == "." {
		num, ok := right.(*Number)
		if !ok {
			return nil, false
		}

		i := int(num.Value)
		if i < 0 || i >= len(t.Value) {
			return nil, false
		}
		return t.Value[i], true
	} else {
		return nil, false
	}
}
