package main

import (
	"github.com/joshvanl/bingo/shell"
)

func main() {
	shell := shell.New()

	for {
		shell.Run()
	}
}
