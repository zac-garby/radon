<h1 align="center">
  Radon
</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/zac-garby/radon)](https://goreportcard.com/report/github.com/zac-garby/radon)

Radon is a minimalistic, dynamic, compiled scripting language which runs on a virtual machine. In this sense it's in the same category as languages such as Python or Ruby. It's in very early development, and obviously isn't ready for production code yet, but it's quite fun to mess around with and see how much it can already do.

If you come across something which should work or should be considered as part of the language, make an issue. At this stage, I'll probably accept most things. If you think you can fix it yourself, still make an issue so I don't accidentally do the same thing, but I'm very happy to accept pull requests :). A number of things have been added to the TODO/ideas list at the bottom of this file, so if you're bored and want something to work on that'd be great!

This is a total rewrite, with more focus on code quality, testing, and some slight improvements to the actual language. The old version has [its own branch](https://github.com/zac-garby/radon/tree/old). This is still a work in progress, so I can't at all guarantee that it works. In fact, I can probably guarantee more that it _won't_.

For syntax highlighting and snippets in [TextMate](http://macromates.com), install [Radon.tmbundle](Radon.tmbundle) from this repository.

```
calc input = do
  result = 0

  for char in input do
    result = match char where
      | "+" -> result + 1,
      | "-" -> result - 1,
      | "*" -> result * 2,
      | "/" -> result / 2
  end
end
```

### What's implemented so far?

 - The AST
   - Tested
 - The lexer
   - Tested
 - The parser
   - Tested
 - All objects
   - Tested
 - Bytecode
   - Defined all instructions
   - Parsing instructions
   - Tested
 - Virtual Machine
   - All instructions
   - Scopes, stacks, etc...
   - Uh, everything else
   - Untested! :O

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

### TODO, or Some ideas
 - Some Haskell-style operators:
   - `|>` operator, e.g. `5 |> print`
   - `$` operator, e.g. `print $ 5 + 3`
 - Might be able to optimise tuple compilation by flattening the tree and calling `MakeTuple`
   - Probably only a very small performance increase though, but potentially worthwhile for large tuples
 - Make tuples _actually_ be stored using contiguous memory
 - Currying, so calling a function with too few arguments would make a new curried function
 - Empty stack after each statement. Store bytecode indices of the start of each statement.
 - Parse lists as a circumfix operator `[ ... ]` with a tuple inside
 - Parse maps as a circumfix operator `{ ... }` with a tuple of tuples inside, making the `:` operator the same as `,`, but possibly a different precedence
 - Max call-stack size
 - Builtins!