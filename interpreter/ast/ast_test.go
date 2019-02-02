package ast

import (
	"os"
	"testing"

	"github.com/joshvanl/bingo/interpreter/ast/errors"
)

func Test_Parse(t *testing.T) {
	try(t, "foo > foo", nil)
	try(t, "foo foo > foo", nil)
	try(t, "foo > ", errors.MissingExpression)
	try(t, "> ", errors.MissingExpression)
	try(t, "foo && foo foooo foo", nil)
	try(t, "foo && foo fooo foo  && foo", nil)

	//TODO:
	//try(t, "foo && &&", errors.MissingExpression)
	//try(t, "foo && && foo", errors.MissingExpression)
}

func try(t *testing.T, str string, exp error) {
	p, err := Parse(&str)
	if err != nil {
		if exp == nil {
			fatal(t, exp, err)
		}

		if exp.Error() != err.Error() {
			fatal(t, exp, err)
		}
	}

	for _, s := range p.Statements {
		err := s.Prepare(os.Stdin, os.Stdout, os.Stderr)

		if err != nil {
			if exp == nil {
				fatal(t, exp, err)
			}

			if exp.Error() != err.Error() {
				fatal(t, exp, err)
			} else {
				return
			}
		}
	}

	if exp != nil {
		fatal(t, exp, nil)
	}
}

func fatal(t *testing.T, exp, got error) {
	t.Fatalf("exp=%v, got=%v", exp, got)
}
