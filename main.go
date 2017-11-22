package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Zac-Garby/lang/parser"
)

const (
	prompt = "~ "

	load = ":load "
)

// The REPL
func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt)
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		text = strings.TrimRight(text, "\n")

		if strings.HasPrefix(text, load) {
			loadFile(strings.TrimPrefix(text, load))
		}

		parse := parser.New(text, "repl")
		prog := parse.Parse()

		if len(parse.Errors) > 0 {
			parse.PrintErrors()
		} else {
			fmt.Println(prog.Tree())
		}
	}
}

func loadFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}

	text, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	parse := parser.New(string(text), name)
	prog := parse.Parse()

	if len(parse.Errors) > 0 {
		parse.PrintErrors()
	} else {
		fmt.Println(prog.Tree())
	}

	return nil
}
