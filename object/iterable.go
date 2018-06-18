package object

// An Iterable object can be iterated and looped over in for-loops. It is a superset
// of the Object interface.
type Iterable interface {
	Object
	Next() (Object, bool)
}

// A ListIterable is an iterable which operates over each element in a list.
type ListIterable struct {
	defaults
	List  *List
	Index int
}

func (i *ListIterable) String() string {
	return "<iterable>"
}

// Type returns the type of an Object.
func (i *ListIterable) Type() Type {
	return NilType
}

// Equals checks whether or not two objects are equal to each other.
func (i *ListIterable) Equals(_ Object) bool {
	return false
}

// Next returns the next object from the iterable. If false is returned
// as the second return value, the iterable has finished.
func (i *ListIterable) Next() (Object, bool) {
	if i.Index < 0 || i.Index >= len(i.List.Value) {
		return nil, false
	}

	val := i.List.Value[i.Index]
	i.Index++
	return val, true
}

// Iter turns an object into an iterable.
func (i *ListIterable) Iter() (Iterable, bool) {
	return i, true
}
