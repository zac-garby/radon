package object

// Nil represents the absence of any useful value.
type Nil struct {
	defaults
}

func (n *Nil) String() string {
	return "nil"
}

// Type returns the type of an Object.
func (n *Nil) Type() Type {
	return NilType
}

// Equals checks whether or not two objects are equal to each other.
func (n *Nil) Equals(other Object) bool {
	return other.Type() == NilType
}

// Infix applies a infix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (n *Nil) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: []Object{n, right},
		}, true
	}

	return nil, false
}

// Numeric returns the numeric value of an object, or false if it can't be a number.
func (n *Nil) Numeric() (float64, bool) {
	return 0, true
}
