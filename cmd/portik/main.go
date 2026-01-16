package main

import (
	"os"

	"portik/internal/cli"
)

func main() {
	// takes the arg and passes to cli RUn
	os.Exit(cli.Run(os.Args[1:]))
}
