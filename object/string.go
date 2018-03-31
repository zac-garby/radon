package object

import (
	"fmt"
)

// A String is a UTF-8 encoded Unicode string.
type String struct {
	defaults
	Value string
}

func (s *String) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

func (s *String) Type() Type {
	return StringType
}

func (s *String) Equals(other Object) bool {
	switch o := other.(type) {
	case *String:
		return s.Value == o.Value

	default:
		return false
	}
}

func (s *String) Infix(op string, right Object) (Object, bool) {
	switch r := right.(type) {
	case *String:
		switch op {
		case "+":
			return &String{Value: s.Value + r.Value}, true
		case "<":
			return &Boolean{Value: s.Value < r.Value}, true
		case ">":
			return &Boolean{Value: s.Value > r.Value}, true
		case "<=":
			return &Boolean{Value: s.Value <= r.Value}, true
		case ">=":
			return &Boolean{Value: s.Value >= r.Value}, true
		default:
			return nil, false
		}
	default:
		return nil, false
	}
}
