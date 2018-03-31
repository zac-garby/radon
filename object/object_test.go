package object_test

import (
	"fmt"
	"testing"

	. "github.com/Zac-Garby/radon/object"
)

func n(val float64) *Number {
	return &Number{Value: val}
}

func TestStringify(t *testing.T) {
	cases := map[Object]string{
		n(5):   "5",
		n(3.7): "3.7",
	}

	for o, s := range cases {
		if o.String() != s {
			fmt.Printf("%s != %s\n", o.String(), s)
			t.Fail()
		}
	}
}

func TestEquals(t *testing.T) {
	cases := []struct {
		a, b Object
		eq   bool
	}{
		{n(5), n(5), true},
		{n(10), n(11), false},
	}

	for _, c := range cases {
		if c.a.Equals(c.b) != c.eq {
			fmt.Printf("%v should equal %v: %t\n", c.a, c.b, c.eq)
			t.Fail()
		}
	}
}

func TestPrefix(t *testing.T) {
	cases := []struct {
		op      string
		in, out Object
	}{
		{"-", n(5), n(-5)},
		{"+", n(5), n(5)},
	}

	for _, c := range cases {
		got, ok := c.in.Prefix(c.op)
		if !ok {
			fmt.Printf("%v should be able to use prefix op %s\n", c.in, c.op)
			t.Fail()
			continue
		}

		if !got.Equals(c.out) {
			fmt.Printf("%s%v should equal %v\n", c.op, c.in, c.out)
			t.Fail()
		}
	}
}

func TestInfix(t *testing.T) {
	cases := []struct {
		left       Object
		op         string
		right, out Object
	}{
		{n(1), "+", n(2), n(3)},
		{n(1), "-", n(2), n(-1)},
		{n(1), "*", n(2), n(2)},
		{n(1), "/", n(2), n(0.5)},
		{n(1), "^", n(2), n(1)},
		{n(1), "//", n(2), n(0)},
		{n(1), "%", n(2), n(1)},
	}

	for _, c := range cases {
		got, ok := c.left.Infix(c.op, c.right)
		if !ok {
			fmt.Printf("%v should be able to use infix op %s\n", c.left, c.op)
			t.Fail()
			continue
		}

		if !got.Equals(c.out) {
			fmt.Printf("%v %s %v should equal %v\n", c.left, c.op, c.right, c.out)
			t.Fail()
		}
	}
}

func TestNumeric(t *testing.T) {
	cases := map[Object]float64{
		n(5): 5,
	}

	for in, out := range cases {
		val, ok := in.Numeric()
		if !ok {
			fmt.Printf("%v should be numeric\n", in)
			t.Fail()
			continue
		}

		if val != out {
			fmt.Printf("%v should have numeric value %v\n", in, out)
			t.Fail()
		}
	}
}
