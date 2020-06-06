package main

import (
	"fmt"
	"os"

	"github.com/ripta/mwnci/pkg/runner"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		for _, arg := range args {
			f, err := os.Open(arg)
			if err != nil {
				panic(err)
			}

			if err := runner.All(os.Stdout, f); err != nil {
				panic(err)
			}
		}
		return
	}

	fmt.Println("Croeso i mwnci. ^D i adael.")
	if err := runner.REPL(os.Stdout, os.Stdin); err != nil {
		panic(err)
	}
}
