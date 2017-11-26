package store

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
	Data map[string]*Item
}

// New creates an empty store.
func New() *Store {
	return &Store{
		Data: make(map[string]*Item),
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
	s.Data[name] = &Item{
		Name:    name,
		Value:   val,
		IsLocal: local,
	}
}

// Get gets the value of a name.
func (s *Store) Get(name string) (object.Object, bool) {
	val, ok := s.Data[name]
	return val.Value, ok
}

// Clone duplicates a store
func (s *Store) Clone() *Store {
	return &Store{
		Data: s.Data,
	}
}
