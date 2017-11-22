package token

import "fmt"

// A Position represents the position of a token
// in the source code.
type Position struct {
	Line, Column int
	Filename     string
}

func (p *Position) String() string {
	return fmt.Sprintf("%s/%d:%d", p.Filename, p.Line, p.Column)
}
