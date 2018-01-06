package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Zac-Garby/lang/bytecode"
	"github.com/Zac-Garby/lang/compiler"
	"github.com/Zac-Garby/lang/lexer"
	"github.com/Zac-Garby/lang/object"
	"github.com/Zac-Garby/lang/parser"
	"github.com/Zac-Garby/lang/token"
	"github.com/Zac-Garby/lang/vm"
	"github.com/carmark/pseudo-terminal-go/terminal"
	"github.com/fatih/color"
)

type command func(arg string)

const prompt = "~> "

var (
	store     = vm.NewStore()
	printToks = false
	printTree = false
	printCode = false
	execCode  = true

	commands = map[string]command{
		"toks": func(arg string) { printToks = arg == "on" },
		"tree": func(arg string) { printTree = arg == "on" },
		"code": func(arg string) { printCode = arg == "on" },
		"exec": func(arg string) { execCode = arg == "on" },
		"quit": func(_ string) { os.Exit(0) },
	}
)

// The REPL
func main() {
	term, err := terminal.NewWithStdInOut()
	if err != nil {
		panic(err)
	}
	defer term.ReleaseFromStdInOut()
	term.SetPrompt(prompt)

outer:
	for {
		text, err := term.ReadLine()
		if err != nil {
			break
		}

		for name, fn := range commands {
			if !strings.HasPrefix(text, ":"+name) {
				continue
			}

			arg := strings.TrimSpace(strings.TrimLeft(text, ":"+name))
			fn(arg)
			continue outer
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
	if printToks && print {
		lex := lexer.Lexer(input, filename)

		for tok := lex(); tok.Type != token.EOF; tok = lex() {
			fmt.Println(tok.String())
		}
	}

	parse := parser.New(input, filename)
	prog := parse.Parse()

	if len(parse.Errors) > 0 {
		parse.PrintErrors(os.Stderr)
		return nil
	}

	if printTree && print {
		fmt.Println(prog.Tree())
	}

	cmp := compiler.New()
	if err := cmp.Compile(prog); err != nil {
		return err
	}

	code, err := bytecode.Read(cmp.Bytes)
	if err != nil {
		return err
	}

	if printCode && print {
		for _, instr := range code {
			fmt.Printf("%20s (%d)\n", instr.Name, instr.Arg)
		}
	}

	if !execCode {
		return nil
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
		} else if val != nil && !(val.Type() == object.TupleType && len(val.(*object.Tuple).Value) == 0) && print {
			color.Cyan(" %s", val.String())
		}
	}

	return nil
}
