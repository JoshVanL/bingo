package interpreter

import (
	"strings"

	"github.com/joshvanl/bingo/command"
)

func Run(pi *string) error {
	if pi == nil {
		return nil
	}

	cmd, args := parse(pi)
	return command.Execute(cmd, args)
}

func parse(pi *string) (string, []string) {
	fields := strings.Fields(*pi)

	if len(fields) == 1 {
		return fields[0], nil
	}

	return fields[0], fields[1:]
}
