package object

import (
	"fmt"
)

type (
	// A Number is a 64-bit floating point
	Number struct {
		Value float64
	}

	// A Boolean is either true or false
	Boolean struct {
		Value bool
	}

	// A String is a string of characters
	String struct {
		Value string
	}

	// Nil is the absence of a value
	Nil struct{}
)

/* Type() implementations */

// Type returns the type of this object
func (n *Number) Type() Type { return NumberType }

// Type returns the type of this object
func (b *Boolean) Type() Type { return BooleanType }

// Type returns the type of this object
func (s *String) Type() Type { return StringType }

// Type returns the type of this object
func (n *Nil) Type() Type { return NilType }

/* String() implementations */

// String returns the nicely formatted string representation of this object
func (n *Number) String() string { return fmt.Sprintf("%v", n.Value) }

// String returns the nicely formatted string representation of this object
func (b *Boolean) String() string { return fmt.Sprintf("%v", b.Value) }

// String returns the nicely formatted string representation of this object
func (s *String) String() string { return s.Value }

// String returns the nicely formatted string representation of this object
func (n *Nil) String() string { return "nil" }

/* Debug() implementations */

// Debug returns the debug string representation of this object
func (n *Number) Debug() string { return fmt.Sprintf("%v", n.Value) }

// Debug returns the debug string representation of this object
func (b *Boolean) Debug() string { return fmt.Sprintf("%v", b.Value) }

// Debug returns the debug string representation of this object
func (s *String) Debug() string { return fmt.Sprintf(`"%s"`, s.Value) }

// Debug returns the debug string representation of this object
func (n *Nil) Debug() string { return "nil" }

/* Equals() implementations */

// Equals checks if two objects are equal to each other
func (n *Number) Equals(o Object) bool {
	if other, ok := o.(*Number); ok {
		return n.Value == other.Value
	}

	return false
}

// Equals checks if two objects are equal to each other
func (b *Boolean) Equals(o Object) bool {
	if other, ok := o.(*Boolean); ok {
		return b.Value == other.Value
	}

	return false
}

// Equals checks if two objects are equal to each other
func (s *String) Equals(o Object) bool {
	if other, ok := o.(*String); ok {
		return s.Value == other.Value
	}

	return false
}

// Equals checks if two objects are equal to each other
func (n *Nil) Equals(o Object) bool {
	_, ok := o.(*Nil)
	return ok
}

/* Collection implementations */

// Elements returns the elements in a collection
func (s *String) Elements() []Object {
	chars := make([]Object, 0, len(s.Value))

	for _, r := range []rune(s.Value) {
		chars = append(chars, &String{Value: string(r)})
	}

	return chars
}

// GetIndex returns the ith element in a collection
func (s *String) GetIndex(i int) Object {
	if i >= len(s.Value) || i < 0 {
		return &Nil{}
	}

	return &String{Value: string([]rune(s.Value)[i])}
}

// SetIndex sets the ith element in a collection to o
func (s *String) SetIndex(i int, o Object) {
	if i >= len(s.Value) || i < 0 {
		return
	}

	if ch, ok := o.(*String); ok && len(ch.Value) == 1 {
		bytes := []rune(s.Value)
		bytes[i] = rune([]rune(ch.Value)[0])
		s.Value = string(bytes)
	}
}
