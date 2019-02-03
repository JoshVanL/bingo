package operators

import (
	"io"
	"os"
	"strings"

	"github.com/joshvanl/bingo/interpreter/ast/errors"
)

func Print(in io.ReadCloser, args []string) (func() error, error) {
	return writeToFile(in, args, os.O_CREATE|os.O_WRONLY, true)
}

func Append(in io.ReadCloser, args []string) (func() error, error) {
	return writeToFile(in, args, os.O_APPEND|os.O_CREATE|os.O_WRONLY, false)
}

func writeToFile(in io.ReadCloser, args []string, flags int, remove bool) (func() error, error) {
	if len(args) < 1 {
		return nil, errors.MissingExpression
	}

	return func() error {
		perm := os.FileMode(0644)
		info, err := os.Stat(args[0])
		if err == nil {
			perm = info.Mode()

			if remove {
				err = os.Remove(args[0])
				if err != nil {
					return err
				}
			}
		}

		f, err := os.OpenFile(args[0], flags, perm)
		if err != nil {
			return err
		}

		if _, err := io.Copy(f, in); err != nil {
			f.Close()
			return err
		}

		if len(args) > 1 {
			_, err = f.Write([]byte(strings.Join(args[1:], " ")))
			if err != nil {
				f.Close()
				return err
			}
		}

		return f.Close()
	}, nil
}
