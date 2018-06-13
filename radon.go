package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"

	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/compiler"
	"github.com/Zac-Garby/radon/lexer"
	"github.com/Zac-Garby/radon/object"
	"github.com/Zac-Garby/radon/parser"
	"github.com/Zac-Garby/radon/runtime"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(c chan os.Signal) {
		<-c
		quit()
	}(c)

	if len(os.Args) < 2 {
		startRepl()
	} else {
		filename := os.Args[1]
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("couldn't open", filename)
			os.Exit(2)
		}

		if _, err = run(string(bytes), runtime.NewStore(nil)); err != nil {
			fmt.Print("\x1b[91m")
			fmt.Println(err)
			fmt.Print("\x1b[0m")

			os.Exit(1)
		}
	}
}

func startRepl() {
	reader := bufio.NewReader(os.Stdin)
	store := runtime.NewStore(nil)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("io error:", err)
			os.Exit(1)
		}

		line = strings.TrimSpace(line)

		res, err := run(line, store)
		if err != nil {
			fmt.Print("\x1b[91m") // red
			fmt.Println(" ", err)
			fmt.Print("\x1b[0m")
		} else if res != nil {
			fmt.Print("\x1b[94m") // blue
			fmt.Println(" ", res)
			fmt.Print("\x1b[0m")
		}
	}
}

func run(code string, store *runtime.Store) (object.Object, error) {
	var (
		l         = lexer.Lexer(code, "repl")
		p         = parser.New(l)
		prog, err = p.Parse()
	)

	if err != nil {
		return nil, err
	}

	c := compiler.New()
	if err := c.Compile(prog); err != nil {
		return nil, err
	}

	parsedCode, err := bytecode.Read(bytes.NewReader(c.Bytes))
	if err != nil {
		return nil, err
	}

	v := runtime.New()
	frame := v.MakeFrame(
		parsedCode,
		nil,
		store,
		c.Constants,
		c.Names,
		c.Jumps,
	)

	v.PushFrame(frame)

	return v.Run()
}

func quit() {
	fmt.Println("quit")
	os.Exit(0)
}
