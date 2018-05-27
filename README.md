<h1 align="center">Radon</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/zac-garby/radon)](https://goreportcard.com/report/github.com/zac-garby/radon)

This is a total rewrite of Radon, with more focus on code quality, testing, and some slight improvements to the actual language. The old version has [its own branch](https://github.com/zac-garby/radon/tree/old). This is still a work in progress, so I can't at all guarantee that it works. In fact, I can probably guarantee more that it _won't_.

For syntax highlighting and snippets in [TextMate](http://macromates.com), install [Radon.tmbundle](Radon.tmbundle) from this repository.

Here's a nice example of some of the features:

![](img/screenshot.png)

### What's implemented so far?

 - The AST
   - Tested
 - The lexer
   - Tested
 - The parser
   - Tested
 - All objects (except functions)
   - Tested
 - Bytecode
   - Defined all instructions
   - Parsing instructions
   - Tested

### Building

To use Radon, you'll have to build it yourself for now, since there's no point in my putting a new version on GitHub every time it changes. First, download the project and `cd` to it:

```
go get github.com/zac-garby/radon
cd $GOPATH/src/github.com/zac-garby/radon
```

Then, install the dependencies and install the `radon` command:

```
dep ensure
go install
```

If `$GOPATH/bin` is in your `$PATH` variable, you can start the REPL using the `radon` command. Otherwise, you'll have to use the actual path to the binary: `$GOPATH/bin/radon`, although I do recommend adding `$GOPATH/bin` to `$PATH`. You also might want to `mv $GOPATH/bin/radon /usr/local/bin`.

### Some ideas

 - Some Haskell-style operators:
   - `|>` operator, e.g. `5 |> print`
   - `$` operator, e.g. `print $ 5 + 3`
 - Might be able to optimise tuple compilation by flattening the tree and calling `MakeTuple`
   - Probably only a very small performance increase though, but potentially worthwhile for large tuples
 - Make tuples _actually_ be stored using contiguous memory