package object

import (
	"fmt"
	"strings"
)

// Split splits a string by the separator.
func (s *String) Split(args ...Object) (Object, error) {
	var sep string

	if len(args) == 0 {
		sep = " "
	} else if len(args) == 1 {
		sep = args[0].String()
	} else {
		return nil, fmt.Errorf("argument: wrong amount of arguments supplied to the function. expected 0 or 1, got %v", len(args))
	}

	var result []Object

	for _, substr := range strings.Split(s.Value, sep) {
		result = append(result, &String{
			Value: substr,
		})
	}

	return &List{result}, nil
}

// GetMethod gets the method of the given name from
// an object.
// TODO: This function, and the whole method system,
// needs to be refactored.
func (s *String) GetMethod(name string) (*Builtin, bool) {
	var builtin func(args ...Object) (Object, error)

	switch name {
	case "split":
		builtin = s.Split

	default:
		return nil, false
	}

	return &Builtin{
		Fn:   builtin,
		Name: name,
	}, true
}
