package object_test

import (
	"testing"

	. "github.com/Zac-Garby/lang/object"
)

func TestEquality(t *testing.T) {
	var (
		n1 = &Number{Value: 0}
		n2 = &Number{Value: 1}

		s1 = &String{Value: "foo"}
		s2 = &String{Value: "bar"}

		a0 = &List{Value: []Object{}}
		a1 = &List{Value: []Object{n1}}
		a2 = &List{Value: []Object{n1, n2}}

		t0 = &Tuple{Value: []Object{}}
		t1 = &Tuple{Value: []Object{n1}}
		t2 = &Tuple{Value: []Object{n1, n2}}

		m0 = &Map{Keys: make(map[string]Object), Values: make(map[string]Object)}
		m1 = &Map{Keys: make(map[string]Object), Values: make(map[string]Object)}
		m2 = &Map{Keys: make(map[string]Object), Values: make(map[string]Object)}
	)

	m1.SetKey(s1, n1)

	m2.SetKey(s1, n1)
	m2.SetKey(s2, n2)

	cases := []struct {
		left, right Object
		shouldEqual bool
	}{
		// Numbers
		{n1, n1, true},
		{n1, n2, false},
		{n1, s1, false},

		// Booleans
		{TrueObj, TrueObj, true},
		{TrueObj, FalseObj, false},
		{TrueObj, n1, false},

		// Strings
		{s1, s1, true},
		{s1, s2, false},
		{s1, n1, false},

		// Null
		{NilObj, NilObj, true},
		{NilObj, n1, false},

		// Arrays
		{a0, a0, true},
		{a0, a1, false},
		{a1, a1, true},
		{a2, a2, true},
		{a0, n1, false},

		// Tuples
		{t0, t0, true},
		{t0, t1, false},
		{t1, t1, true},
		{t2, t2, true},
		{t0, n1, false},

		// Maps
		{m0, m0, true},
		{m0, m1, false},
		{m1, m1, true},
		{m2, m2, true},
		{m0, n1, false},
	}

	for _, pair := range cases {
		eq := pair.left.Equals(pair.right)

		if eq != pair.shouldEqual {
			if pair.shouldEqual {
				t.Errorf("%s doesn't equal %s", pair.left, pair.right)
			} else {
				t.Errorf("%s is equal to %s", pair.left, pair.right)
			}
		}
	}
}
