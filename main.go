package main

import (
	"github.com/pustserg/lab/commands"
	"os"
)

func main() {
	args := os.Args[1:]
	commands.Handle(args)
}
