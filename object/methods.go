package object

// Split splits a string by the separator.
func (s *String) Split(args ...Object) (Object, error) {
	return NilObj, nil
}

// GetMethod gets the method of the given name from
// an object.
func (s *String) GetMethod(name string) (*Builtin, bool) {
	builtins := map[string]func(...Object) (Object, error){
		"split": s.Split,
	}

	builtin, ok := builtins[name]
	if !ok {
		return nil, false
	}

	return &Builtin{
		Fn:   builtin,
		Name: name,
	}, true
}
