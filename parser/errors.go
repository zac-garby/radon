package parser

import (
	"fmt"

	"github.com/Zac-Garby/radon/token"
)

// An Error represents a parsing error. Contains a message and two positions for
// the start and end of the error.
type Error struct {
	Message    string
	Start, End token.Position
}

// Error returns a string representation of an Error, to comply with the error
// interface.
func (e *Error) Error() string {
	return fmt.Sprintf("** Parse error ~ [%s-%s] %s", e.Start.String(), e.End.String(), e.Message)
}

// err creates a new Error. Calls fmt.Sprintf on the message with ...args.
func (p *Parser) err(msg string, start, end token.Position, args ...interface{}) {
	err := &Error{
		Message: fmt.Sprintf(msg, args...),
		Start:   start,
		End:     end,
	}

	p.Errors = append(p.Errors, err)
}

// defaultErr is the same as `p.err`, but assumes the start and end positions to
// be the start and end of the current token.
func (p *Parser) defaultErr(msg string, args ...interface{}) {
	p.err(msg, p.cur.Start, p.cur.End, args...)
}

func (p *Parser) peekErr(t token.Type) {
	if p.peek.Type == token.Semi {
		p.err("unexpected end of line, wanted '%s'", p.peek.Start, p.peek.End, t)
	} else {
		p.err("expected '%s' but got '%s'", p.peek.Start, p.peek.End, t, p.peek.Type)
	}
}

func (p *Parser) unexpected(t token.Type) {
	if t == token.Semi {
		p.defaultErr("unexpected end of line")
	} else {
		p.defaultErr("unexpected token: %s", t)
	}
}
