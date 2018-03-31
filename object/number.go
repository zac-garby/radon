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

func (n *Number) Type() Type {
	return NumberType
}

func (n *Number) Equals(other Object) bool {
	switch o := other.(type) {
	case *Number:
		return n.Value == o.Value

	default:
		return false
	}
}

func (n *Number) Prefix(op string) (Object, bool) {
	var val float64

	switch op {
	case "+":
		val = n.Value

	case "-":
		val = -n.Value

	default:
		return nil, false
	}

	return &Number{Value: val}, true
}

func (n *Number) Infix(op string, right Object) (Object, bool) {
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
	case "%":
		val = float64(int64(l) % int64(r))
	case ",":
		panic("make a tuple")
	default:
		return nil, false
	}

	switch v := val.(type) {
	case float64:
		return &Number{Value: v}, true
	case bool:
		panic("bool")
	}

	return nil, false
}

func (n *Number) Numeric() (float64, bool) {
	return n.Value, true
}
