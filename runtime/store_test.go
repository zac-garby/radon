package runtime_test

import (
	"testing"

	"github.com/Zac-Garby/radon/object"
	. "github.com/Zac-Garby/radon/runtime"
)

func TestStoreCanHandleVariables(t *testing.T) {
	s := NewStore(nil)
	e := NewStore(s)

	s.Set("foo", &object.Number{Value: 5}, true)
	s.Set("bar", &object.String{Value: "hello"}, false)

	foo, ok := s.Get("foo")
	if !ok {
		t.Error("variable 'foo' wasn't set (or cannot be retrieved)")
	} else if !foo.Value.Equals(&object.Number{Value: 5}) {
		t.Error("variable 'foo' doesn't equal 5")
	}

	bar, ok := s.Get("bar")
	if !ok {
		t.Error("variable 'bar' wasn't set (or cannot be retrieved)")
	} else if !bar.Value.Equals(&object.String{Value: "hello"}) {
		t.Error("variable 'bar' doesn't equal 'hello'")
	}

	e.Set("foo", &object.Number{Value: 6}, false)
	eFoo, ok := e.Get("foo")
	if !ok {
		t.Error("variable 'foo' could not be retrieved from scope E after being assigned")
	}
	sFoo, ok := s.Get("foo")
	if !ok {
		t.Error("variable 'foo' could not be retrieved from scope S after being assigned in E")
	}

	if !eFoo.Value.Equals(sFoo.Value) {
		t.Error("variable 'foo' from scope E doesn't equal 'foo' from scope S")
	}
}
