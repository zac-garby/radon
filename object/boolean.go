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

func (b *Boolean) Type() Type {
	return BooleanType
}

func (n *Boolean) Equals(other Object) bool {
	switch o := other.(type) {
	case *Boolean:
		return n.Value == o.Value

	default:
		return false
	}
}

func (b *Boolean) Prefix(op string) (Object, bool) {
	if op == "!" {
		return &Boolean{Value: !b.Value}, true
	}
	return nil, false
}

func (b *Boolean) Infix(op string, right Object) (Object, bool) {
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

func (b *Boolean) Numeric() (float64, bool) {
	if b.Value {
		return 1, true
	}
	return 0, true
}
