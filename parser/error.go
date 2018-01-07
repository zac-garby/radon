package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/Zac-Garby/radon/token"
)

// Error represents a parsing error
type Error struct {
	Message    string
	Start, End token.Position
}

// Error returns a string representation of an Error.
func (e Error) Error() string {
	return fmt.Sprintf("%s → %s ~ %s", e.Start.String(), e.End.String(), e.Message)
}

// Err creates an Error instance with the given arguments
func (p *Parser) Err(msg string, start, end token.Position) {
	err := Error{
		Message: msg,
		Start:   start,
		End:     end,
	}

	p.Errors = append(p.Errors, err)
}

func (p *Parser) defaultErr(msg string) {
	err := Error{
		Message: msg,
		Start:   p.cur.Start,
		End:     p.cur.End,
	}

	p.Errors = append(p.Errors, err)
}

func (p *Parser) peekErr(ts ...token.Type) {
	if len(ts) > 1 {
		msg := "expected either "

		for i, t := range ts {
			msg += string(t)

			if i+1 < len(ts) {
				msg += ", "
			} else if i < len(ts) {
				msg += ", or "
			}
		}

		msg += ", but got " + string(p.peek.Type)

		p.Err(msg, p.peek.Start, p.peek.End)
	} else if len(ts) == 1 {
		msg := fmt.Sprintf("expected %s, but got %s", ts[0], p.peek.Type)
		p.Err(msg, p.peek.Start, p.peek.End)
	}
}

func (p *Parser) curErr(ts ...token.Type) {
	if len(ts) > 1 {
		msg := "expected either "

		for i, t := range ts {
			msg += string(t)

			if i+1 < len(ts) {
				msg += ", "
			} else if i < len(ts) {
				msg += ", or "
			}
		}

		msg += ", but got " + string(p.cur.Type)

		p.Err(msg, p.cur.Start, p.cur.End)
	} else if len(ts) == 1 {
		msg := fmt.Sprintf("expected %s, but got %s", ts[0], p.cur.Type)
		p.Err(msg, p.cur.Start, p.cur.End)
	}
}

func (p *Parser) unexpectedTokenErr(t token.Type) {
	msg := fmt.Sprintf("unexpected token: %s", t)
	p.defaultErr(msg)
}

func (p *Parser) printError(index int, w io.Writer) {
	err := p.Errors[index]
	fmt.Fprintln(w, err.Error())
}

// PrintErrors prints all parser errors in a nice format
func (p *Parser) PrintErrors(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}

	if len(p.Errors) == 0 {
		return
	}

	for i := range p.Errors {
		p.printError(i, w)
	}
}
