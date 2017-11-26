package vm

import (
	"fmt"
)

// ErrType is a runtime error type
type ErrType string

const (
	// ErrInternal is thrown for any internal vm problems
	ErrInternal = "Internal"

	// ErrUnknown is thrown when there's an error, but the vm isn't
	// sure of what nature it is
	ErrUnknown = "Unknown"

	// ErrNoInstruction is thrown if an instruction in the bytecode
	// isn't yet implemented
	ErrNoInstruction = "NoInstruction"

	// ErrNotFound is thrown if a name, symbol, or index isn't found
	ErrNotFound = "NotFound"

	// ErrWrongType is thrown if an object is of the wrong type to be
	// operated on
	ErrWrongType = "WrongType"

	// ErrNoOp is thrown if an operator isn't defined for the given
	// operands
	ErrNoOp = "NoOp"

	// ErrSyntax is thrown for any syntax errors which couldn't be
	// found in the parsing stage
	ErrSyntax = "Syntax"
)

// Error is a runtime error thrown in the virtual machine
type Error struct {
	Type    ErrType
	Message string
}

// Err creates a new runtime error with the given message and type
func Err(msg string, t ErrType) *Error {
	return &Error{
		Type:    t,
		Message: msg,
	}
}

// Errf creates a new runtime error, formatting msg with format
func Errf(msg string, t ErrType, format ...interface{}) *Error {
	return &Error{
		Type:    t,
		Message: fmt.Sprintf(msg, format...),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
