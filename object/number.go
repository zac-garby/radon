package object

import (
	"fmt"
	"math"
)

// A Number is a 64-bit floating point decimal.
type Number struct {
	defaults
	Value float64
}

func (n *Number) String() string {
	return fmt.Sprintf("%v", n.Value)
}

// Type returns the type of an Object.
func (n *Number) Type() Type {
	return NumberType
}

// Equals checks whether or not two objects are equal to each other.
func (n *Number) Equals(other Object) bool {
	switch o := other.(type) {
	case *Number:
		return n.Value == o.Value

	default:
		return false
	}
}

// Prefix applies a prefix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (n *Number) Prefix(op string) (Object, bool) {
	var val float64

	switch op {
	case "+":
		val = n.Value

	case "-":
		val = -n.Value

	case ",":
		return &Tuple{Value: []Object{n}}, true

	default:
		return nil, false
	}

	return &Number{Value: val}, true
}

// Infix applies a infix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (n *Number) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: []Object{n, right},
		}, true
	}

	l := n.Value

	r, ok := right.Numeric()
	if !ok {
		return nil, false
	}

	var val interface{}
	switch op {
	case "+":
		val = l + r
	case "-":
		val = l - r
	case "*":
		val = l * r
	case "/":
		val = l / r
	case ">":
		val = l > r
	case "<":
		val = l < r
	case ">=":
		val = l >= r
	case "<=":
		val = l <= r
	case "^":
		val = math.Pow(l, r)
	case "//":
		val = math.Floor(l / r)
	case "|":
		val = float64(int64(l) | int64(r))
	case "&":
		val = float64(int64(l) & int64(r))
	case "%":
		val = float64(int64(l) % int64(r))
	default:
		return nil, false
	}

	switch v := val.(type) {
	case float64:
		return &Number{Value: v}, true
	case bool:
		return &Boolean{Value: v}, true
	}

	return nil, false
}

// Numeric returns the numeric value of an object, or false if it can't be a number.
func (n *Number) Numeric() (float64, bool) {
	return n.Value, true
}
