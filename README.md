[![Join the chat at https://gitter.im/radon-lang/Lobby](https://badges.gitter.im/radon-lang/Lobby.svg)](https://gitter.im/radon-lang/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

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

Since the VM is stack-based, there are builtin functions to manipulate the stack.
When an expression is evaluated, it's implicitly pushed to the stack.
```
5
pop() + 1   # 6

1
2
3
list(3)     # [1, 2, 3]
```

# TODO

**Improvements**
 - Add more tests
 - Standard library
   - HTTP
   - String transformations
   - Look through Go std libraries
 - More builtins
   - Stack operations, `dup()`, `rot()`, etc...
   - Files, `open()`, `write()`, etc...
 - Go interop
 - Don't print the REPL result if it's `()`
 - Remove `Token()` method from AST nodes
 - Add better stack memory management
   - Maybe clear the stack after every statement. However, that would probably
     mean an instruction would have to be added after each statement, bloating
     the bytecode.

**New language features**
 - Optional/non-optional variables
   - A variable ending in `?` can be set to `nil`
   - Other variables cannot
   - This will be checked at compile time
 - Variadic parameters
   - e.g. `f(x, ...y) = x + sum(y)`
   - `f(1, 2, 3)` &rarr; `6`
 - Model inheritance
   - `model (x) : parent (x, 0)`
   - see github.com/Zac-Garby/radonuage
 - For loops over collections
   - e.g. `for item in [1, 2, 3]` or `for i, item in [1, 2, 3]`
 - Importing things
   - e.g. `import "path/to/directory"` imports the files in the directory
     - Only in the top level, to allow package-private things
   - Global package store
     - Some directory
     - `package "package-name"` imports the directory relative to the global store
     - Bash command to import packages from GitHub by username and repo
