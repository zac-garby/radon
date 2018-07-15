package object

import (
	"fmt"
	"strings"

	"github.com/cnf/structhash"
)

// A Map maps keys to values, where keys and values can be any type.
//
// Keys maps the hashes to the keys they were
// hashed from, and Values maps the hashes to
// the keys' corresponding values.
type Map struct {
	defaults
	Keys   map[string]Object
	Values map[string]Object
}

func (m *Map) String() string {
	if len(m.Keys) == 0 {
		return "{}"
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

	return fmt.Sprintf("{%s}", strings.Join(stringArr, ", "))
}

// Type returns the type of an Object.
func (m *Map) Type() Type {
	return MapType
}

// Equals checks whether or not two objects are equal to each other.
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

// Prefix applies a prefix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (m *Map) Prefix(op string) (Object, bool) {
	if op == "," {
		return &Tuple{Value: []Object{m}}, true
	}

	return nil, false
}

// Infix applies a infix operator to an object, returning the result. If the operation
// cannot be performed, (nil, false) is returned.
func (m *Map) Infix(op string, right Object) (Object, bool) {
	if op == "," {
		return &Tuple{
			Value: []Object{m, right},
		}, true
	}

	return nil, false
}

// Items returns a slice containing all objects in an Object, or false otherwise.
func (m *Map) Items() ([]Object, bool) {
	var pairs []Object

	for hash, key := range m.Keys {
		val := m.Values[hash]
		pairs = append(pairs, &Tuple{
			Value: []Object{
				key, val,
			},
		})
	}

	return pairs, true
}

// SetSubscript sets the value of a subscript of an Object, e.g. foo[bar] = baz.
// Returns false if it can't be done.
func (m *Map) SetSubscript(key Object, val Object) bool {
	if hash, err := structhash.Hash(key, 1); err == nil {
		m.Values[hash] = val
		m.Keys[hash] = key

		return true
	}

	return false
}
