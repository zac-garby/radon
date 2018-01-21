<h1 align="center">Radon</h1>

[![Join the chat at https://gitter.im/radon-lang/Lobby](https://badges.gitter.im/radon-lang/Lobby.svg)](https://gitter.im/radon-lang/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Build Status](https://travis-ci.org/Zac-Garby/radon.svg?branch=master)](https://travis-ci.org/Zac-Garby/radon)

Radon is a dynamically typed, compiled scripting language which runs on a virtual
machine.

This is what it looks like:

Function definitions look a lot like how you'd define a function in maths.
```
# Function definition
fib(n) = match n where
    | 0 -> 1,
    | 1 -> 1,
    | _ -> fib(n - 1) + fib(n - 2)

v = fib(20)
print(v)
```

The object system is based on models, which are basically like structs or
classes.
```
vector = model(x, y)

# A method
vector.translate(x, y) = {
    self.x = self.x + x
    self.y = self.y + y
}

# Instantiation
pos = vector(101, 38)
```

# TODO

**Improvements**
 - Refactor virtual machine
   - Have the actual fetch-execute cycle on the VM, not the frame
   - Have a call stack in the VM
   - Every cycle, execute an instruction from the top frame in the call stack
   - Store names in the frame, not the store
     - The store should just store variables
   - Other minor things
   - Maybe rewrite it entirely
 - Improve REPL
   - Don't use `pseudo-terminal`, because it doesn't allow unicode input
   - Enable multi-line input
 - Add more tests
 - Standard library
   - HTTP
   - String transformations
   - Look through Go std libraries
 - More builtins
   - Files, `open()`, `write()`, etc...
 - Builtin models
   - `file`
   - Networking stuff
     - Sockets
     - HTTP
 - Go interop
 - Remove `Token()` method from AST nodes
 - Add better stack memory management
   - Maybe clear the stack after every statement. However, that would probably
     mean an instruction would have to be added after each statement, bloating
     the bytecode.

**New language features**
 - Make a Radon bytecode file format (maybe `.rex` for Radon executable)
   - Contains bytecode, constants, names, etc...
   - Either in an unreadable binary format, or a readable, editable (JSON?) format
   - Can be executed directly using `radon -x <filename>`
 - Implement the operator-assignment operators (e.g. `+=`)
   - They already parse properly, they just need compilation
 - Allow models to be named
   - `model vector(x, y)` instead of `vector = model(x, y)`
   - Will be printed as `<model vector>` instead of just `<model>`
 - Allow functions to be named
   - `a(x, y) = x + y` will be called `a`
   - Will be printed as `<function a>` instead of `<function>`
 - Ruby-like code blocks after functions
   - `f(x) { print("hello") }` is the same as `f(x, => print("hello"))`
   - `f(x) { |a, b| a + b }` is the same as `f(x, (a, b) => a + b)`
 - String interpolation
   - `x = "world"; "hello ${x}"` is the same as `"hello world"`
   - Only occurs in double-quoted strings
 - Variadic parameters
   - e.g. `f(x, ...y) = x + sum(y)`
   - `f(1, 2, 3)` &rarr; `6`
 - Model inheritance
   - `model (x) : parent (x, 0)`
   - see github.com/Zac-Garby/language
 - Special model methods
   - `init` is called when the model is instantiated
   - Various ones for operators, such as `plus`
 - For loops over collections
   - e.g. `for item in [1, 2, 3]` or `for i, item in [1, 2, 3]`
 - Global package store
   - Some directory, maybe `$RADONPATH`
   - `import "* package-name"` imports the directory relative to the global store
   - Bash command to import packages from GitHub by username and repo
