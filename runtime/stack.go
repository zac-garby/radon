package runtime

import (
	"github.com/Zac-Garby/radon/object"
)

var (
	// MaxDataStackSize is the maximum amount of objects which are allowed to
	// be inside the data stack.
	MaxDataStackSize = 100000

	// DefaultDataStackCapacity is the capacity initially allocated to the stack. A
	// higher value will possibly increase performance, but use more RAM.
	DefaultDataStackCapacity = 32

	// ErrDataStackOverflow tells the user that too many objects were pushed to
	// the data stack.
	ErrDataStackOverflow = makeError(InternalError, "too many objects on the data stack")

	// ErrDataStackUnderflow tells the user that a pop was attempted but no objects
	// were on the stack.
	ErrDataStackUnderflow = makeError(InternalError, "no objects to pop from the data stack")
)

// A Stack is created for each frame in the virtual machine, and stores a stack of
// objects which can be popped or pushed.
type Stack struct {
	Objects []object.Object
}

// NewStack makes a new empty Stack, with capacity equal to DefaultStackCapacity.
func NewStack() *Stack {
	return &Stack{
		Objects: make([]object.Object, 0, DefaultDataStackCapacity),
	}
}

// Len gets the length of the stack.
func (s *Stack) Len() int {
	return len(s.Objects)
}

// Push pushes an Object to the top of a Stack. If it returns an error, it will be
// ErrDataStackOverflow.
func (s *Stack) Push(obj object.Object) error {
	if s.Len() >= MaxDataStackSize {
		return ErrDataStackOverflow
	}
	s.Objects = append(s.Objects, obj)
	return nil
}

// Pop pops an Object from the Stack, returning it. If it returns an error, it will
// be ErrDataStackUnderflow.
func (s *Stack) Pop() (object.Object, error) {
	top, err := s.Top()
	if err != nil {
		return nil, err
	}

	s.Objects = s.Objects[:s.Len()-1]
	return top, nil
}

// Top returns the top item in the Stack. If there are no items, it will return
// an ErrDataStackUnderflow. It is equivalent to popping, but without removing
// the item from the stack.
func (s *Stack) Top() (object.Object, error) {
	if s.Len() == 0 {
		return nil, ErrDataStackUnderflow
	}

	return s.Objects[s.Len()-1], nil
}
