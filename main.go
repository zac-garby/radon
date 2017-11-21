package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Zac-Garby/lang/parser"
)

// The REPL
func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		text = strings.TrimRight(text, "\n")

		parse := parser.New(text, "repl")
		prog := parse.Parse()

		if len(parse.Errors) > 0 {
			parse.PrintErrors()
		} else {
			fmt.Println(prog.Tree())
		}
	}
}
