package object

import (
	"github.com/Zac-Garby/lang/bytecode"
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
