package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	curUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s, this is the Monkey programming language\n", curUser)
	fmt.Printf("Feel free to start typing commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
