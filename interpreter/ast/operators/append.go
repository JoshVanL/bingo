package operators

import (
	"io"
	"os"
	"strings"

	"github.com/joshvanl/bingo/interpreter/ast/errors"
)

func Append(in io.ReadCloser, args []string) (func() error, error) {
	if len(args) < 1 {
		return nil, errors.MissingExpression
	}

	return func() error {
		f, err := os.OpenFile(args[0], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		if _, err := io.Copy(f, in); err != nil {
			return err
		}

		if len(args) > 1 {
			_, err = f.Write([]byte(strings.Join(args[1:], " ")))
			if err != nil {
				return err
			}
		}

		return f.Close()
	}, nil
}
