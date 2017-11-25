package object

import (
	"github.com/Zac-Garby/lang/bytecode"
)

// A Function is a piece of code which can
// be called.
type Function struct {
	Code      bytecode.Code
	Constants []Object
	Names     []string

	// A function to be called every time this Function
	// is called
	OnCall func(f *Function) Object
}

// Type returns the type of the object
func (f *Function) Type() Type { return FunctionType }

// Equals checks if the function is equal to another object
func (f *Function) Equals(o Object) bool { return false }

// String returns a string representing the function
func (f *Function) String() string { return "<function>" }

// Debug returns a string representing the function
func (f *Function) Debug() string { return "<function>" }
