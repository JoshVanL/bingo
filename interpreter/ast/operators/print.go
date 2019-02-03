package operators

import (
	"io"
	"os"
	"strings"

	"github.com/joshvanl/bingo/interpreter/ast/errors"
)

func Print(in *os.File, args []string, stop <-chan struct{}) (func(<-chan os.Signal) error, error) {
	return writeToFile(in, args, os.O_CREATE|os.O_WRONLY, true, stop)
}

func Append(in *os.File, args []string, stop <-chan struct{}) (func(<-chan os.Signal) error, error) {
	return writeToFile(in, args, os.O_APPEND|os.O_CREATE|os.O_WRONLY, false, stop)
}

func writeToFile(in *os.File, args []string, flags int, remove bool, stop <-chan struct{}) (func(<-chan os.Signal) error, error) {
	if len(args) < 1 {
		return nil, errors.MissingExpression
	}

	return func(ch <-chan os.Signal) error {
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

		done := make(chan struct{})
		go func() {
			select {
			case <-ch:
				f.Close()
			case <-stop:
				f.Close()
			case <-done:
			}
		}()

		io.Copy(f, in)
		close(done)

		if len(args) > 1 {
			_, err = f.Write([]byte(strings.Join(args[1:], " ")))
		}

		f.Close()

		return err
	}, nil
}
