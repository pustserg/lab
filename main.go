package main

import (
	"fmt"
	"os"
	"github.com/pustserg/lab/commands"
)

func main() {
	args := os.Args[1:]
	commands.Handle(args)
}
