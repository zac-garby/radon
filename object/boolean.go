package object

import "fmt"

// A Boolean is either true or false.
type Boolean struct {
	defaults
	Value bool
}

func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

// Type returns the type of an Object.
func (b *Boolean) Type() Type {
	return BooleanType
}

// Equals checks whether or not two objects are equal to each other.
func (b *Boolean) Equals(other Object) bool {
	switch o := other.(type) {
	case *Boolean:
		return b.Value == o.Value

	default:
		return false
	}
}

// Prefix applies a prefix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (b *Boolean) Prefix(op string) (Object, bool) {
	if op == "!" {
		return &Boolean{Value: !b.Value}, true
	} else if op == "," {
		return &Tuple{Value: []Object{b}}, true
	}
	return nil, false
}

// Infix applies a infix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (b *Boolean) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: []Object{b, right},
		}, true
	}

	switch r := right.(type) {
	case *Boolean:
		switch op {
		case "&&", "&":
			return &Boolean{Value: b.Value && r.Value}, true
		case "||", "|":
			return &Boolean{Value: b.Value || r.Value}, true
		default:
			return nil, false
		}
	default:
		return nil, false
	}
}

// Numeric returns the numeric value of an object, or false if it can't be a number.
func (b *Boolean) Numeric() (float64, bool) {
	if b.Value {
		return 1, true
	}
	return 0, true
}
