package main

import (
	"github.com/joshvanl/bingo/shell"
)

func main() {
	shell := shell.New()

	defer func() {
		shell.TerminalOldState()
	}()

	for {
		shell.Run()
	}
}
