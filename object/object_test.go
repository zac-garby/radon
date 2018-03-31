package object_test

import (
	"fmt"
	"testing"

	. "github.com/Zac-Garby/radon/object"
)

func n(val float64) *Number {
	return &Number{Value: val}
}

func b(val bool) *Boolean {
	return &Boolean{Value: val}
}

func s(val string) *String {
	return &String{Value: val}
}

func l(vals ...Object) *List {
	return &List{Value: vals}
}

func tu(vals ...Object) *Tuple {
	return &Tuple{Value: vals}
}

func TestStringify(t *testing.T) {
	cases := map[Object]string{
		n(5):                 "5",
		n(3.7):               "3.7",
		b(true):              "true",
		b(false):             "false",
		s("foo"):             `"foo"`,
		&Nil{}:               "nil",
		l(n(1), n(2), n(3)):  "[1, 2, 3]",
		tu(n(1), n(2), n(3)): "(1, 2, 3)",
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
		{n(1), b(true), false},

		{b(true), b(true), true},
		{b(true), b(false), false},
		{b(false), n(5), false},

		{s("foo"), s("foo"), true},
		{s("foo"), s("bar"), false},
		{s("foo"), n(5), false},

		{&Nil{}, &Nil{}, true},
		{&Nil{}, n(5), false},

		{l(n(1), n(2)), l(n(1), n(2)), true},
		{l(n(1)), l(n(2)), false},
		{l(), l(), true},
		{l(n(1)), l(), false},
		{l(n(1)), l(n(1), n(2), n(3)), false},

		{tu(n(1), n(2)), tu(n(1), n(2)), true},
		{tu(n(1)), tu(n(2)), false},
		{tu(), tu(), true},
		{tu(n(1)), tu(), false},
		{tu(n(1)), tu(n(1), n(2), n(3)), false},
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

		{"!", b(true), b(false)},
		{"!", b(false), b(true)},
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
		{n(5), ">", n(1), b(true)},
		{n(5), ">=", n(5), b(true)},
		{n(1), "<", n(5), b(true)},
		{n(1), "<=", n(1), b(true)},
		{n(1), "|", n(2), n(3)},
		{n(1), "&", n(2), n(0)},
		{n(1), ",", n(2), tu(n(1), n(2))},

		{b(true), "&&", b(false), b(false)},
		{b(false), "||", b(true), b(true)},
		{b(true), "&", b(false), b(false)},
		{b(false), "|", b(true), b(true)},
		{b(true), ",", b(false), tu(b(true), b(false))},

		{s("foo"), "+", s("bar"), s("foobar")},
		{s("foo"), "<", s("bar"), b(false)},
		{s("foo"), ">", s("bar"), b(true)},
		{s("foo"), "<=", s("foo"), b(true)},
		{s("foo"), ">=", s("foo"), b(true)},
		{s("foo"), ",", s("bar"), tu(s("foo"), s("bar"))},

		{&Nil{}, ",", n(8), tu(&Nil{}, n(8))},

		{l(n(1), n(2)), "+", l(n(3), n(4)), l(n(1), n(2), n(3), n(4))},
		{l(n(1), n(2), n(3)), "[]", n(1), n(2)},
		{l(n(1), n(2)), ",", l(n(3), n(4)), tu(l(n(1), n(2)), l(n(3), n(4)))},

		{tu(n(1), n(2), n(3)), "[]", n(1), n(2)},
		{tu(n(1), n(2), n(3)), ".", n(2), n(3)},
		{tu(n(1)), ",", n(2), tu(n(1), n(2))},
		{tu(n(1)), ",", tu(n(2)), tu(n(1), tu(n(2)))},
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

		b(true):  1,
		b(false): 0,

		&Nil{}: 0,
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

func TestItems(t *testing.T) {
	cases := map[Object][]Object{
		l(n(1), n(2), n(3)):  []Object{n(1), n(2), n(3)},
		tu(n(1), n(2), n(3)): []Object{n(1), n(2), n(3)},
	}

	for in, out := range cases {
		items, ok := in.Items()
		if !ok {
			fmt.Printf("%v should be able to get items\n", in)
			t.Fail()
			continue
		}

		if len(items) != len(out) {
			fmt.Printf("got %v, expected %v\n", items, out)
			t.Fail()
			continue
		}

		for i, it := range items {
			if !it.Equals(out[i]) {
				fmt.Printf("%v != %v\n", it, out[i])
				t.Fail()
				continue
			}
		}
	}
}
