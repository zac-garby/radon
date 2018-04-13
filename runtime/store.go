package runtime

import "github.com/Zac-Garby/radon/object"

// A Variable represents a Radon variable. A Variable includes a Name
// and Value.
type Variable struct {
	Name  string
	Value object.Object
}

// A Store contains all the variables defined in a particular scope. A
// Store also has a pointer to the enclosing scope.
type Store struct {
	Data      map[string]*Variable
	Enclosing *Store
}

// NewStore creates a new empty store with the given enclosing scope (can be nil).
func NewStore(enclosing *Store) *Store {
	return &Store{
		Data:      make(map[string]*Variable),
		Enclosing: enclosing,
	}
}

// Get gets a variable from the store. If it isn't found, it checks the enclosing scope,
// and so on.
func (s *Store) Get(name string) (*Variable, bool) {
	v, ok := s.Data[name]
	if !ok {
		if s.Enclosing != nil {
			return s.Enclosing.Get(name)
		}
		return nil, false
	}

	return v, true
}

// Set sets a variable in the store. If declare is false, enclosing scopes will be
// assigned to instead of this one, if the variable is already defined there.
func (s *Store) Set(name string, val object.Object, declare bool) {
	if !declare && s.Enclosing != nil {
		if _, ok := s.Enclosing.Get(name); ok {
			s.Enclosing.Set(name, val, false)
			return
		}
	}

	s.Data[name] = &Variable{
		Name:  name,
		Value: val,
	}
}
