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

		if err := run(string(bytes), runtime.NewStore(nil)); err != nil {
			fmt.Print("\x1b[91m")
			fmt.Println("error:", err)
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

		if err := run(line, store); err != nil {
			fmt.Print("\x1b[91m")
			fmt.Println("error:", err)
			fmt.Print("\x1b[0m")
		}
	}
}

func run(code string, store *runtime.Store) error {
	var (
		l         = lexer.Lexer(code, "repl")
		p         = parser.New(l)
		prog, err = p.Parse()
	)

	if err != nil {
		return err
	}

	c := compiler.New()
	if err := c.Compile(prog); err != nil {
		return err
	}

	parsedCode, err := bytecode.Read(bytes.NewReader(c.Bytes))
	if err != nil {
		return err
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

	res, err := v.Run()
	if err != nil {
		return err
	}

	fmt.Print("\x1b[94m")
	fmt.Println("<", res)
	fmt.Print("\x1b[0m")

	return nil
}

func quit() {
	fmt.Println("quit")
	os.Exit(0)
}
