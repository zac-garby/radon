package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/compiler"
	"github.com/Zac-Garby/lang/parser"
	"github.com/Zac-Garby/lang/vm"
)

const (
	prompt = "~> "
	load   = ":load "
)

var (
	store = vm.NewStore()
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
			if err := loadFile(strings.TrimPrefix(text, load)); err != nil {
				os.Stderr.WriteString(err.Error() + "\n")
			}

			continue
		}

		if err := execute(text, "repl", true, store); err != nil {
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

	return execute(string(text), name, false, store)
}

func execute(input, filename string, print bool, sto *vm.Store) error {
	parse := parser.New(input, filename)
	prog := parse.Parse()

	if len(parse.Errors) > 0 {
		parse.PrintErrors(os.Stderr)
		return nil
	}

	cmp := compiler.New()
	if err := cmp.Compile(prog); err != nil {
		return err
	}

	code, err := bytecode.Read(cmp.Bytes)
	if err != nil {
		return err
	}

	sto.Names = cmp.Names

	v := vm.New()
	v.Run(code, sto, cmp.Constants)

	if err := v.Error(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	} else {
		val, err := v.ExtractValue()
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
		} else if val != nil && print {
			fmt.Println(val)
		}
	}

	return nil
}
