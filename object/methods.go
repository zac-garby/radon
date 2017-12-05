package object

import "strings"

// Split splits a string by the separator.
func (s *String) Split(sep string) Object {
	strings := strings.Split(s.Value, sep)
	return MakeObj(strings)
}
