package errors

import (
	"errors"
	"fmt"
)

var (
	BadOperator       = errors.New("could not parse operator")
	MissingExpression = errors.New("parse error near '\\n'")
)

func Join(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	if len(errs) == 1 {
		return errs[0]
	}

	err := errs[0]
	for _, e := range errs[1:] {
		err = fmt.Errorf("%s\n%s", err, e)
	}

	return err
}
