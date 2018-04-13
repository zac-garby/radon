package object

// Nil represents the absense of any useful value.
type Nil struct {
	defaults
}

func (n *Nil) String() string {
	return "nil"
}

func (n *Nil) Type() Type {
	return NilType
}

func (n *Nil) Equals(other Object) bool {
	return other.Type() == NilType
}

func (n *Nil) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: []Object{n, right},
		}, true
	}

	return nil, false
}

func (n *Nil) Numeric() (float64, bool) {
	return 0, true
}
