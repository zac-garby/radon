package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Zac-Garby/radon/lexer"
	"github.com/Zac-Garby/radon/parser"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("io error:", err)
			os.Exit(1)
		}

		line = strings.TrimSpace(line)

		if err := run(line); err != nil {
			fmt.Println(err)
		}
	}
}

func run(code string) error {
	var (
		l         = lexer.Lexer(code, "repl")
		p         = parser.New(l)
		prog, err = p.Parse()
	)

	if err != nil {
		return err
	}

	spew.Println(prog)

	return nil
}
