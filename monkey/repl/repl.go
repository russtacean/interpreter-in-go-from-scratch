package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">>"
const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func Start(in io.Reader, out io.Writer) {
	useVM := true // TODO parameterize this via command line
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		parse := parser.New(lex)

		program := parse.ParseProgram()
		if len(parse.Errors()) > 0 {
			printParserErrors(out, parse.Errors())
			continue
		}

		if useVM {
			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "Whoops, compile error:\n %s\n", err)
				continue
			}

			machine := vm.New(comp.Bytecode())
			err = machine.Run()
			if err != nil {
				fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
				continue
			}
			stackTop := machine.StackTop()
			io.WriteString(out, stackTop.Inspect())
			io.WriteString(out, "\n")

		} else {
			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Whoops, we ran into some monkey business\n")
	io.WriteString(out, "parser errors:\n")
	for _, errMsg := range errors {
		io.WriteString(out, "\t"+errMsg+"\n")
	}
}
