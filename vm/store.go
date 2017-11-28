package vm

import (
	"github.com/Zac-Garby/lang/object"
)

// An Item is created for each variable stored
// in a Store.
type Item struct {
	Name    string
	Value   object.Object
	IsLocal bool
}

// A Store contains the defined variables for
// a scope.
type Store struct {
	Data  map[string]*Item
	Names []string
}

// NewStore creates an empty store.
func NewStore() *Store {
	return &Store{
		Data:  make(map[string]*Item),
		Names: make([]string, 16),
	}
}

// Contains checks if the store contains a variable
// with the given name.
func (s *Store) Contains(name string) bool {
	_, ok := s.Data[name]
	return ok
}

// Set defines a name in the store.
func (s *Store) Set(name string, val object.Object, local bool) {
	if s.Contains(name) {
		s.Data[name].IsLocal = local
		s.Data[name].Value = val
	} else {
		s.Names = append(s.Names, name)

		s.Data[name] = &Item{
			Name:    name,
			Value:   val,
			IsLocal: local,
		}
	}
}

// Get gets the value of a name.
func (s *Store) Get(name string) (object.Object, bool) {
	val, ok := s.Data[name]
	if !ok {
		return nil, false
	}

	return val.Value, true
}

// Clone duplicates a store
func (s *Store) Clone() *Store {
	return &Store{
		Data: s.Data,
	}
}

// GetNameIndex gets the index in Names of the given name
func (s *Store) GetNameIndex(name string) rune {
	for i, n := range s.Names {
		if n == name {
			return rune(i)
		}
	}

	panic("name not found")
}
