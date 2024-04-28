# interpreter-in-go-from-scratch
**Current status: complete**

This is a repo I'm devoting to the code I write as I follow 2 books:
- *Writing An Interpreter In Go* by Thorsten Ball: [link](https://interpreterbook.com/)
- *Writing A Compiler In Go* by Thorsten Ball: [link](https://compilerbook.com/)

## Monkey Language
These books walk you through how to create an interpreter for a language created for the books called Monkey. Here is some example Monkey code:
```
let x = 5;
let y = 6;
let z = x + y

let array = [1, 2, 3];
let newArray = push(array, 4);

let hash = {"foo": 1, "bar": 2}
hash[foo]

let newAdder = fn(a, b) {
    fn(c) { a + b + c };
};
let adder = newAdder(1, 2);
adder(8);
```
Monkey supports the following features:
- integers
- booleans
- strings
- arrays
 - hashes
- prefix-, infix- and index operators
- conditionals
- global and local bindings
- first-class functions
- return statements
- closures

## Commits
This repo is structured with commits I made as I went through each of the books. Commits have the chapter and section in them, implementing the contents of that section. Any commit with the words _extra credit_ were additional work I did that was left as an exercise for the reader or functionality I wanted to implement based on other languages (e.g. truthy/falsy values for some types)

## Interpreter Book
This books creates the base language syntax and a tree walking interpreter for it
- Chapter 1: Creates a lexer to let us tokenize Monkey code
- Chapter 2: Implements a Pratt parser to turn our tokenized code into an abstract syntax tree (AST)
- Chapter 3: Creates a tree-walking evaluator to let us evaluate our AST
- Chapter 4: Adds strings, composite data types (array and hashmap), and built-in functions (e.g. `len()`)

## Compiler Book
This books implements a stack based VM and a bytecode compiler for it using the lexer and parser from the first book. I treat the chapters in this book as continuations of the first book in my commits, but the first chapters are conceptual
- Chapter 6: Basic VM and compiler structure, support simple expression of `1 + 2`
- Chapter 7: Compiling and evaluating basic expressions with integers and bools (prefix and infix operations)
- Chapter 8: Support conditionals with jumps
- Chapter 9: Handle global variable bindings using a symbol table
- Chapter 10: Support strings and composite data types like arrays and hashes
- Chapter 11: Support functions and local variable bindings
- Chapter 12: Support built-in functions like `len()`, `append()`, etc
- Chapter 13: Support closures and recursive functions
- Chapter 14: Write benchmark for different engines and add support for REPL

