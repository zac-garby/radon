package vm

import (
	"errors"

	"github.com/Zac-Garby/lang/object"
)

type stack struct {
	objects []object.Object
}

func newStack() *stack {
	return &stack{
		objects: []object.Object{},
	}
}

func (s *stack) push(obj object.Object) {
	s.objects = append(s.objects, obj)
}

func (s *stack) pop() (object.Object, error) {
	last, err := s.top()
	if err != nil {
		return nil, err
	}

	s.objects = s.objects[:len(s.objects)-1]
	return last, nil
}

func (s *stack) top() (object.Object, error) {
	if len(s.objects) == 0 {
		return nil, errors.New("stack: not enough items in the stack to pop")
	}

	return s.objects[len(s.objects)-1], nil
}

func (s *stack) dup() error {
	top, err := s.top()
	if err != nil {
		return err
	}

	s.objects = append(s.objects, top)
	return nil
}
