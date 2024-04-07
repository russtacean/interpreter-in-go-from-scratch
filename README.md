# interpreter-in-go-from-scratch

This is a repo I'm devoting to the code I write as I follow 2 books:
- *Writing An Interpreter In Go* by Thorsten Ball: [link](https://interpreterbook.com/)
- *Writing A Compiler In Go* by Thorsten Ball: [link](https://compilerbook.com/)

## Monkey Language
These books walk you through how to create an interpreter for a language created for the books called Monkey. Here is some example Monkey code:
```
let x = 5;
let addFive = fn(y) {
  return x + y;
}
let twelve = addFive(7);

let array = [1, 2, 3];
let newArray = push(array, 4);

let hash = {"foo": 1, "bar": 2}
hash[foo]
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
- Chapter 2: Creates a parser to turn our tokenized code into an abstract syntax tree (AST)
- Chapter 3: Creates a tree-walking evaluator to let us evaluate our AST
- Chapter 4: Adds strings, data structures (array and hashmap), and built-in functions (e.g. len())

