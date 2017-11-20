package lexer

import (
	"regexp"
	"unicode"

	"github.com/Zac-Garby/lang/token"
)

var lineEndings = []token.Type{
	token.ID,
	token.String,
	token.Number,
	token.True,
	token.False,
	token.Nil,
	token.Break,
	token.Next,
	token.Return,
	token.RightParen,
	token.RightSquare,
	token.RightBrace,
}

// Lexer takes a string and returns a stream of tokens
// The stream of tokens is in the form of a function
// which returns the next token.
func Lexer(str, file string) func() token.Token {
	var (
		index = 0
		col   = 1
		line  = 1
		ch    = make(chan token.Token)
	)

	go func() {
		for {
			if index < len(str) {
				foundSpace := false

				for index < len(str) && (unicode.IsSpace(rune(str[index])) || str[index] == '#') {
					if unicode.IsSpace(rune(str[index])) {
						index++
						col++

						if str[index-1] == '\n' {
							col = 1
							line++
						}

						foundSpace = true
					} else {
						for index < len(str) && str[index] != '\n' {
							index++
						}

						col = 1
					}
				}

				if foundSpace {
					continue
				}

				found := false

				remainingSubstring := str[index:]

				for _, pair := range lexicalDictionary {
					var (
						regex   = pair.regex
						handler = pair.handler
						pattern = regexp.MustCompile(regex)
						match   = pattern.FindStringSubmatch(remainingSubstring)
					)

					if len(match) > 0 {
						found = true
						t, literal, whole := handler(match)
						l := len(whole)

						ch <- token.Token{
							Type:    t,
							Literal: literal,
							Start:   token.Position{Line: line, Column: col, Filename: file},
							End:     token.Position{Line: line, Column: col + l - 1, Filename: file},
						}

						index += l
						col += l

						for index < len(str) && unicode.IsSpace(rune(str[index])) && str[index] != '\n' {
							index++
							col++
						}

						if index < len(str) && str[index] == '#' {
							for index < len(str) && str[index] != '\n' {
								index++
							}
						}

						isLineEnding := false

						for _, ending := range lineEndings {
							if t == ending {
								isLineEnding = true
							}
						}

						if (isLineEnding && index < len(str) && (str[index] == '\n' || str[index] == '}')) || index >= len(str) {
							ch <- token.Token{
								Type:    token.Semi,
								Literal: ";",
								Start:   token.Position{Line: line, Column: col, Filename: file},
								End:     token.Position{Line: line, Column: col, Filename: file},
							}
						}

						break
					}
				}

				if !found {
					ch <- token.Token{
						Type:    token.Illegal,
						Literal: string(str[index]),
						Start:   token.Position{Line: line, Column: col, Filename: file},
						End:     token.Position{Line: line, Column: col, Filename: file},
					}

					index++
					col++
				}
			} else {
				index++
				col++

				ch <- token.Token{
					Type:    token.EOF,
					Literal: "",
					Start:   token.Position{Line: line, Column: col, Filename: file},
					End:     token.Position{Line: line, Column: col, Filename: file},
				}
			}
		}
	}()

	return func() token.Token {
		return <-ch
	}
}
