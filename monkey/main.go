package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s, this is the Monkey programming language\n", user)
	fmt.Printf("Feel free to start typing commands")
	repl.Start(os.Stdin, os.Stdout)
}
