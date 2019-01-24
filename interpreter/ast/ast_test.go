package ast

import (
	"testing"
)

func Test_parseStatement(t *testing.T) {
	_, err := parseStatement("foo > foo")
	try(t, nil, err)

	_, err = parseStatement("foo foo > foo")
	try(t, nil, err)

	_, err = parseStatement("foo >  ")
	try(t, moreOperator, err)

	_, err = parseStatement(">  ")
	try(t, moreOperator, err)

	_, err = parseStatement("foo && foo fooo foo  ")
	try(t, nil, err)

	_, err = parseStatement("foo && foo fooo foo  && foo")
	try(t, nil, err)

	_, err = parseStatement("foo && foo fooo foo  && foo && ")
	try(t, moreOperator, err)

	_, err = parseStatement("foo && &&")
	try(t, badOperator, err)

	_, err = parseStatement("foo && && foo")
	try(t, badOperator, err)
}

func try(t *testing.T, exp, got error) {
	if exp == nil && got == nil {
		return
	}

	if exp != nil && got == nil {
		fatal(t, exp, got)
	}

	if exp == nil && got != nil {
		fatal(t, exp, got)
	}

	if exp.Error() != got.Error() {
		fatal(t, exp, got)
	}
}

func fatal(t *testing.T, exp, got error) {
	t.Fatalf("exp=%v, got=%v", exp, got)
}
