package parser

import (
	"github.com/Zac-Garby/radon/token"
)

const (
	lowest = iota
	assign
	lambda
	or
	and
	bitOr
	bitAnd
	equals
	compare
	sum
	product
	exp
	prefix
	call
	index
)

var precedences = map[token.Type]int{
	token.Assign:         assign,
	token.Declare:        assign,
	token.AndEquals:      assign,
	token.BitAndEquals:   assign,
	token.BitOrEquals:    assign,
	token.ExpEquals:      assign,
	token.FloorDivEquals: assign,
	token.MinusEquals:    assign,
	token.ModEquals:      assign,
	token.OrEquals:       assign,
	token.PlusEquals:     assign,
	token.SlashEquals:    assign,
	token.StarEquals:     assign,
	token.Or:             or,
	token.And:            and,
	token.BitOr:          bitOr,
	token.BitAnd:         bitAnd,
	token.Equal:          equals,
	token.NotEqual:       equals,
	token.LessThan:       compare,
	token.GreaterThan:    compare,
	token.LessThanEq:     compare,
	token.GreaterThanEq:  compare,
	token.Plus:           sum,
	token.Minus:          sum,
	token.Star:           product,
	token.Slash:          product,
	token.Mod:            product,
	token.Exp:            exp,
	token.FloorDiv:       exp,
	token.Bang:           prefix,
	token.Dot:            index,
	token.LeftSquare:     index,
	token.LeftParen:      call,
	token.LambdaArrow:    lambda,
}
