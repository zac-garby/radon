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

		if err := run(string(bytes)); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}

func startRepl() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("io error:", err)
			os.Exit(1)
		}

		line = strings.TrimSpace(line)

		if err := run(line); err != nil {
			fmt.Println("error:", err)
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
		runtime.NewStore(nil),
		runtime.NewStore(nil),
		c.Constants,
		c.Names,
		c.Jumps,
	)

	v.PushFrame(frame)

	res, err := v.Run()
	if err != nil {
		return err
	}

	fmt.Println("<", res)

	return nil
}

func quit() {
	fmt.Println("quit")
	os.Exit(0)
}
