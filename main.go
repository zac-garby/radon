package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Zac-Garby/lang/compiler"

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

		if err := execute(text, "repl"); err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
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

	return execute(string(text), name)
}

func execute(input, filename string) error {
	parse := parser.New(input, filename)
	prog := parse.Parse()

	if len(parse.Errors) > 0 {
		parse.PrintErrors()
		return nil
	}

	cmp := compiler.New()
	if err := cmp.Compile(prog); err != nil {
		return err
	}

	fmt.Println(cmp.Bytes)

	return nil
}
