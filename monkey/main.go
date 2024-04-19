package main

import (
	"flag"
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

var useVM = flag.Bool("useVM", true, "Set to use bytecode VM for better performance")

func main() {
	flag.Parse()

	curUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s, this is the Monkey programming language\n", curUser)
	fmt.Printf("Feel free to start typing commands\n")
	repl.Start(os.Stdin, os.Stdout, *useVM)
}
