package object

import (
	"fmt"
	"strings"

	"github.com/cnf/structhash"
)

type (
	// A Tuple is an immutable set of values
	Tuple struct {
		Value []Object
	}

	// A List is a mutable set of values
	List struct {
		Value []Object
	}

	// A Map is a mapping of keys to values.
	//
	// Keys maps the hashes to the keys they were
	// hashed from, and Values maps the hashes to
	// the keys' corresponding values.
	//
	// The Model field, if not nil, specifies the
	// model from which the map derives. If it is
	// nil, it doesn't have a model and just acts
	// as a normal hash map.
	Map struct {
		Model  *Model
		Keys   map[string]Object
		Values map[string]Object
	}
)

/* Type() methods */

// Type returns the type of this object
func (t *Tuple) Type() Type { return TupleType }

// Type returns the type of this object
func (l *List) Type() Type { return ListType }

// Type returns the type of this object
func (m *Map) Type() Type { return MapType }

/* Equals() methods */

// Equals checks if two objects are equal to each other
func (t *Tuple) Equals(o Object) bool {
	if other, ok := o.(*Tuple); ok {
		if len(other.Value) != len(t.Value) {
			return false
		}

		for i, elem := range t.Value {
			if !elem.Equals(other.Value[i]) {
				return false
			}
		}

		return true
	}

	return false
}

// Equals checks if two objects are equal to each other
func (l *List) Equals(o Object) bool {
	if other, ok := o.(*List); ok {
		if len(other.Value) != len(l.Value) {
			return false
		}

		for i, elem := range l.Value {
			if !elem.Equals(other.Value[i]) {
				return false
			}
		}

		return true
	}

	return false
}

// Equals checks if two objects are equal to each other
func (m *Map) Equals(o Object) bool {
	if other, ok := o.(*Map); ok {
		if len(other.Values) != len(m.Values) {
			return false
		}

		for k, v := range m.Values {
			if _, ok := other.Values[k]; !ok {
				return false
			}

			if !v.Equals(other.Values[k]) {
				return false
			}
		}

		return true
	}

	return false
}

func join(arr []Object, sep string) string {
	sarr := make([]string, len(arr))

	for i, elem := range arr {
		sarr[i] = elem.String()
	}

	return strings.Join(sarr, sep)
}

/* String() methods */

// String returns the type of this object
func (t *Tuple) String() string { return "(" + join(t.Value, ", ") + ")" }

// String returns the type of this object
func (l *List) String() string { return "[" + join(l.Value, ", ") + "]" }

// String returns the type of this object
func (m *Map) String() string {
	if len(m.Keys) == 0 {
		return "[:]"
	}

	stringArr := make([]string, len(m.Values))
	i := 0

	for k, v := range m.Values {
		stringArr[i] = fmt.Sprintf(
			"%s: %s",
			m.Keys[k].String(),
			v.String(),
		)

		i++
	}

	return fmt.Sprintf("[%s]", strings.Join(stringArr, ", "))
}

// Debug returns the debug string for the object
func (t *Tuple) Debug() string { return t.String() }

// Debug returns the debug string for the object
func (l *List) Debug() string { return l.String() }

// Debug returns the debug string for the object
func (m *Map) Debug() string { return m.String() }

/* Collection implementations */

// Elements returns the elements in a collection
func (t *Tuple) Elements() []Object {
	return t.Value
}

// GetIndex returns the ith element in a collection
func (t *Tuple) GetIndex(i int) Object {
	if i >= len(t.Value) || i < 0 {
		return NilObj
	}

	return t.Value[i]
}

// SetIndex sets the ith element in a collection to o
func (t *Tuple) SetIndex(i int, o Object) {
	if i >= len(t.Value) || i < 0 {
		return
	}

	t.Value[i] = o
}

// Elements returns the elements in a collection
func (l *List) Elements() []Object {
	return l.Value
}

// GetIndex returns the ith element in a collection
func (l *List) GetIndex(i int) Object {
	if i >= len(l.Value) || i < 0 {
		return NilObj
	}

	return l.Value[i]
}

// SetIndex sets the ith element in a collection to o
func (l *List) SetIndex(i int, o Object) {
	if i >= len(l.Value) || i < 0 {
		return
	}

	l.Value[i] = o
}

/* Container implementations */

// GetKey gets an object at the given key
func (m *Map) GetKey(key Object) Object {
	if hash, err := structhash.Hash(key, 1); err == nil {
		if val, ok := m.Values[hash]; ok {
			return val
		}
	}

	super := m.Model.GetKey(key)
	fn, ok := super.(*Function)
	if !ok {
		return super
	}

	method := &Method{
		Function: fn,
		Map:      m,
	}

	return method
}

// SetKey sets an object at the given key
func (m *Map) SetKey(key, value Object) {
	if hash, err := structhash.Hash(key, 1); err == nil {
		m.Values[hash] = value
		m.Keys[hash] = key
	}
}
