package object

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
