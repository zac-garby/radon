package vm

import (
	"fmt"
	"strings"

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

func (s *stack) pop() object.Object {
	last := s.top()
	s.objects = s.objects[:len(s.objects)-1]
	return last
}

func (s *stack) top() object.Object {
	return s.objects[len(s.objects)-1]
}

func (s *stack) dup() {
	s.objects = append(s.objects, s.top())
}

func (s *stack) String() string {
	var strs []string

	for _, obj := range s.objects {
		strs = append(strs, obj.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}
