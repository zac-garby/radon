package runtime

import (
	"fmt"
)

// ErrorType specifies the type of a runtime error.
type ErrorType string

const (
	_ ErrorType = ""

	// RuntimeError is used for an error of unknown/generic nature.
	RuntimeError = "Runtime"

	// TypeError is used for an error where types are incompatible or wrong.
	TypeError = "Type"

	// InternalError is used for internal errors, and when something is wrong with the bytecode.
	InternalError = "Internal"

	// NameError is used for undefined variables, functions, and other names.
	NameError = "Name"

	// ArgumentError is used for invalid arguments to a function or model.
	ArgumentError = "Argument"

	// StructureError is used for a structure error, such as a break outside a loop.
	StructureError = "Structure"

	// IndexError is used when an invalid index/key is used.
	IndexError = "Index"
)

// An Error represents any type of runtime error (not just RuntimeError), and implements
// the error interface.
type Error struct {
	Type    ErrorType
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("** %s error ~ %s", e.Type, e.Message)
}

func makeError(t ErrorType, format string, args ...interface{}) error {
	return &Error{
		Type:    t,
		Message: fmt.Sprintf(format, args...),
	}
}
