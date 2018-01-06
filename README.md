This is a scripting language. It's currently unnamed. It's dynamically typed and
compiled. It compiles to [bytecode](./bytecode) which is executed on the [virtual
machine](vm).

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
 - Think of a name
 - Add more tests
 - Standard library
   - HTTP
   - String transformations
   - Look through Go std libraries
 - More builtins
   - Stack operations, `dup()`, `rot()`, etc...
   - Files, `open()`, `write()`, etc...
 - Go interop

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
   - see github.com/Zac-Garby/language
 - For loops over collections
   - e.g. `for item in [1, 2, 3]` or `for i, item in [1, 2, 3]`
 - Importing files
   - e.g. `import "path/to/file"`
   - e.g. `import "path/to/directory"` imports the files in the directory
     - Only in the top level, to allow package-private things
   - File importing will probably be done by inserting the AST of the
     imported file directly into the importing file's AST.
   - Global package store
     - Some directory
     - `package "package-name"` imports the directory relative to the global store
     - Bash command to import packages from GitHub by username and repo
