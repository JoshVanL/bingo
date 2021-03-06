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
	if len(arg) < 1 {
		return arg
	}

	if len(arg) > 1 && arg[0] == '~' &&
		arg[1] == '/' {
		arg = os.Getenv("HOME") + arg[1:]
	}

	n := 0
	i := rune(arg[0])
	for m, a := range arg[1:] {
		if i == '/' && a == '/' {
			n = m + 1
		}

		i = a
	}

	return arg[n:]
}
