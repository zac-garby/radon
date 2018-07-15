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

// Type returns the type of an Object.
func (s *String) Type() Type {
	return StringType
}

// Equals checks whether or not two objects are equal to each other.
func (s *String) Equals(other Object) bool {
	switch o := other.(type) {
	case *String:
		return s.Value == o.Value

	default:
		return false
	}
}

// Prefix applies a prefix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (s *String) Prefix(op string) (Object, bool) {
	if op == "," {
		return &Tuple{Value: []Object{s}}, true
	}

	return nil, false
}

// Items returns a slice containing all objects in an Object, or false otherwise.
func (s *String) Items() ([]Object, bool) {
	var (
		runes = []rune(s.Value)
		strs  = make([]Object, len(runes))
	)

	for i, r := range runes {
		strs[i] = &String{
			Value: string(r),
		}
	}

	return strs, true
}

// Infix applies a infix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (s *String) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: []Object{s, right},
		}, true
	}

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

// Iter creates an iterable from an Object.
func (s *String) Iter() (Iterable, bool) {
	items, _ := s.Items()

	return &ListIterable{
		List:  &List{Value: items},
		Index: 0,
	}, true
}

// SetSubscript sets the value of a subscript of an Object, e.g. foo[bar] = baz.
// Returns false if it can't be done.
func (s *String) SetSubscript(index Object, to Object) bool {
	num, ok := index.(*Number)
	if !ok {
		return false
	}

	i := int(num.Value)
	if i < 0 || i >= len(s.Value) {
		return false
	}

	toStr, ok := to.(*String)
	if !ok {
		return false
	}

	toRunes := []rune(toStr.Value)
	if len(toRunes) != 1 {
		return false
	}

	toRune := toRunes[0]

	runes := []rune(s.Value)
	if i < 0 || i >= len(runes) {
		return false
	}

	runes[i] = toRune
	s.Value = string(runes)

	return true
}
