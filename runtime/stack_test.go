package runtime_test

import (
	"testing"

	"github.com/Zac-Garby/radon/object"
	. "github.com/Zac-Garby/radon/runtime"
)

func TestStack(t *testing.T) {
	s := NewStack()

	if err := s.Push(&object.Number{Value: 5}); err != nil {
		t.Error("could not push a number to the stack")
	}

	obj, err := s.Pop()
	if err != nil {
		t.Error("could not pop from the stack")
	}

	if !obj.Equals(&object.Number{Value: 5}) {
		t.Errorf("popped object %s doesn't equal 5\n", obj.String())
	}
}
