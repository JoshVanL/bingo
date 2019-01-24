package interpreter

import (
	"os"
	"strings"

	"github.com/joshvanl/bingo/command"
)

func Parse(pi *string) func(ch <-chan os.Signal) error {
	if pi == nil || len(strings.TrimSpace(*pi)) == 0 {
		return func(ch <-chan os.Signal) error {
			return nil
		}
	}

	cmd, args := split(pi)
	c := command.NewBin(cmd, args)
	return func(ch <-chan os.Signal) error {
		return c.Execute(ch)
	}
}

func split(pi *string) (string, []string) {
	fields := strings.Fields(*pi)

	if len(fields) == 1 {
		return fields[0], nil
	}

	return fields[0], fields[1:]
}
