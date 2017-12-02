package object

import "fmt"

// A Model is a user-defined type, similar to
// a struct in Go.
type Model struct {
	Parameters []string
}

// Type returns the type of this object
func (m *Model) Type() Type { return ModelType }

// Equals checks if two objects are equal to each other
func (m *Model) Equals(o Object) bool { return false }

// String returns a string representing an object
func (m *Model) String() string { return "<model>" }

// Debug returns a more verbose representation of an object
func (m *Model) Debug() string { return "<model>" }

// Instantiate creates a new map with the required fields
// for the model.
func (m *Model) Instantiate(args ...Object) (Object, error) {
	result := &Map{
		Model:  m,
		Keys:   make(map[string]Object),
		Values: make(map[string]Object),
	}

	if len(args) != len(m.Parameters) {
		return nil, fmt.Errorf(
			"instantiation: wrong amount of arguments. expected %v, got %v",
			len(m.Parameters),
			len(args),
		)
	}

	for i, arg := range args {
		name := &String{m.Parameters[i]}
		result.SetKey(name, arg)
	}

	return result, nil
}
