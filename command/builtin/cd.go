package builtin

import (
	"errors"
	"os"
)

func Cd(args []string) error {

	if len(args) == 0 {
		return os.Chdir(os.Getenv("HOME"))
	}

	if len(args) > 2 {
		return errors.New("cd: expecting single directory argument")
	}

	return os.Chdir(reduce(args[0]))
}

func reduce(arg string) string {
	if len(arg) < 2 {
		return arg
	}

	n := 0
	i := rune(arg[0])
	for m, a := range arg[1:] {
		if i == '/' && a == '/' {
			n = m
		}

		i = a
	}

	return arg[n:]
}
