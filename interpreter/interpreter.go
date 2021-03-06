package interpreter

import (
	"os"
	"strings"

	"github.com/joshvanl/bingo/interpreter/ast"
)

func Parse(pi *string) (func(ch <-chan os.Signal) error, error) {
	if pi == nil || len(strings.TrimSpace(*pi)) == 0 {
		return func(ch <-chan os.Signal) error {
			return nil
		}, nil
	}

	prog, err := ast.Parse(pi)
	if err != nil {
		return nil, err
	}

	for _, s := range prog.Statements {
		err := s.Prepare(os.Stdin, os.Stdout, os.Stderr)
		if err != nil {
			return nil, err
		}
	}

	return func(ch <-chan os.Signal) error {
		for _, s := range prog.Statements {
			if err := Run(s, ch); err != nil {
				return err
			}
		}

		return nil
	}, nil
}

func split(pi *string) (string, []string) {
	fields := strings.Fields(*pi)

	if len(fields) == 1 {
		return fields[0], nil
	}

	return fields[0], fields[1:]
}
