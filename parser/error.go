package parser

import (
	"fmt"
	"strings"

	"github.com/Zac-Garby/lang/token"
	"github.com/fatih/color"
)

// Error represents a parsing error
type Error struct {
	Message    string
	Start, End token.Position
}

// Error returns a string representation of an Error.
func (e Error) Error() string {
	return fmt.Sprintf("%s:%s ~ %s", e.Start.String(), e.End.String(), e.Message)
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

func (p *Parser) printError(index int) {
	err := p.Errors[index]

	fmt.Printf("%s → %s\t%s\n", err.Start.String(), err.End.String(), err.Message)
}

func (p *Parser) printErrorVerbose(index int) {
	err := p.Errors[index]

	fmt.Printf("in ")

	var (
		lines = strings.Split(p.text, "\n")
		grey  = color.New(color.FgHiWhite)
		red   = color.New(color.FgRed).Add(color.Bold)
	)

	grey.Printf("    %d| ", err.Start.Line)
	fmt.Printf("%s\n", lines[err.Start.Line-1])
	red.Printf(
		"    %s %s%s\n",
		strings.Repeat(" ", len(fmt.Sprintf("%d", err.Start.Line))),
		strings.Repeat(" ", err.Start.Column),
		strings.Repeat("^", err.End.Column-err.Start.Column+1),
	)

	red.Printf("%s → %s\t%s\n\n", err.Start.String(), err.End.String(), err.Message)
}

// PrintErrors prints all parser errors in a nice format
func (p *Parser) PrintErrors() {
	if len(p.Errors) == 0 {
		return
	}

	for i := range p.Errors {
		if i == 0 {
			p.printErrorVerbose(i)
		} else {
			p.printError(i)
		}
	}
}
