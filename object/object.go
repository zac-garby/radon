package object

// The Type of an Object indicates what type of object it is.
type Type string

// The set of available types in Radon.
const (
	_ Type = ""

	NumberType   = "number"
	BooleanType  = "boolean"
	StringType   = "string"
	ListType     = "list"
	TupleType    = "tuple"
	MapType      = "map"
	NilType      = "nil"
	FunctionType = "function"
	MethodType   = "method"
	BuiltinType  = "builtin"
	ModelType    = "model"
)

// An Object is the interface which every Radon object implements.
type Object interface {
	String() string
	Equals(Object) bool
	Type() Type

	// Prefix performs a prefix operation on an Object.
	// operator can be one of:
	// + - ! ,
	// If the 2nd return value is false, an error is raised.
	Prefix(operator string) (Object, bool)

	// Infix performs an infix operation on an Object.
	// operator can be one of:
	// + - * / == != < > || && | & ^ // % <= >= . ,
	// If the 2nd return value is false, an error is raised.
	Infix(operator string, right Object) (Object, bool)

	// Numeric returns the numeric value of an Object.
	// If the 2nd return value is false, an error is raised.
	Numeric() (float64, bool)

	// Items returns a slice of Objects representing an Object.
	// If the 2nd return value is false, an error is raised.
	Items() ([]Object, bool)

	// Subscript implements the [] operator, e.g. list[5]
	Subscript(Object) (Object, bool)

	// SetSubscript implements assigning to the [] operator, e.g. list[5] = "foo"
	SetSubscript(index Object, to Object) bool
}

// defaults supplies default implementations so other Object types automatically
// implement the methods.
type defaults struct{}

func (d *defaults) String() string                      { panic("not implemented") }
func (d *defaults) Type() Type                          { panic("not implemented") }
func (d *defaults) Equals(Object) bool                  { return false }
func (d *defaults) Prefix(string) (Object, bool)        { return nil, false }
func (d *defaults) Infix(string, Object) (Object, bool) { return nil, false }
func (d *defaults) Numeric() (float64, bool)            { return -1, false }
func (d *defaults) Items() ([]Object, bool)             { return nil, false }
func (d *defaults) Subscript(Object) (Object, bool)     { return nil, false }
func (d *defaults) SetSubscript(Object, Object) bool    { return false }

// IsTruthy checks whether or not an object is "truthy"
func IsTruthy(o Object) bool {
	return !(o.Equals(&Nil{}) ||
		o.Equals(&Boolean{Value: false}))
}
