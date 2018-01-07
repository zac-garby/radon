package object

import (
	"fmt"

	"github.com/Zac-Garby/radon/bytecode"
)

// A Function is a piece of code which can
// be called.
type Function struct {
	Parameters []string

	Code      bytecode.Code
	Constants []Object
	Names     []string
}

// Type returns the type of the object
func (f *Function) Type() Type { return FunctionType }

// Equals checks if the function is equal to another object
func (f *Function) Equals(o Object) bool { return false }

// String returns a string representing the function
func (f *Function) String() string { return "<function>" }

// Debug returns a string representing the function
func (f *Function) Debug() string { return "<function>" }

// A Method is a function bound to a
// particular map.
type Method struct {
	*Function
	*Map
}

// Type returns the type of the object
func (f *Method) Type() Type { return MethodType }

// Equals checks if the method is equal to another object
func (f *Method) Equals(o Object) bool { return false }

// String returns a string representing the method
func (f *Method) String() string { return "<method>" }

// Debug returns a string representing the method
func (f *Method) Debug() string { return "<method>" }

// A Builtin a function usable through the language
// but which is implemented in Go.
type Builtin struct {
	Fn   func(args ...Object) (Object, error)
	Name string
}

// Type returns the type of the object
func (b *Builtin) Type() Type { return BuiltinType }

// Equals checks if the builtin is equal to another object
func (b *Builtin) Equals(o Object) bool { return false }

// String returns a string representing the builtin
func (b *Builtin) String() string { return fmt.Sprintf("<builtin '%s'>", b.Name) }

// Debug returns a string representing the builtin
func (b *Builtin) Debug() string { return fmt.Sprintf("<builtin '%s'>", b.Name) }
