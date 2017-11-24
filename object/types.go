package object

// Type represents the type of an object
type Type string

const (
	_ Type = "" // Placeholder

	// NumberType is a numerical value, e.g. 5
	NumberType = "number"

	// BooleanType is either true or false
	BooleanType = "boolean"

	// StringType is a string of characters, e.g. "foo"
	StringType = "string"

	// ListType is a list of items, e.g. [1, 2, 3]
	ListType = "list"

	// TupleType is an immutable list of items, e.g. (1, 2, 3)
	TupleType = "tuple"

	// SetType is a hashed set of unique items, e.g. set[1, 2, 3]
	SetType = "set"

	// MapType is a hash map of key value mappings, e.g. map[foo: bar, x: y]
	MapType = "map"

	// NilType is the nil constant
	NilType = "nil"

	// FunctionType is a callable function, e.g. f(x, y) = x + y
	FunctionType = "function"

	// CollectionType is an abstract, non-concrete type, which represents
	// any type which can be thought of as a series of items.
	CollectionType = "collection"

	// ContainerType is another abstract, non-concrete type, which
	// represents any type which can be thought of a mapping of keys
	// to values.
	ContainerType = "container"

	// AnyType represents any single type.
	AnyType = "any"
)

func is(o Object, t Type) bool {
	switch t {
	case AnyType:
		return true

	case CollectionType:
		_, ok := o.(Collection)
		return ok

	case ContainerType:
		_, ok := o.(Container)
		return ok

	default:
		return o.Type() == t
	}
}
