package ast

import (
	"io"
	"os"
	"strings"

	"github.com/joshvanl/bingo/command"
	"github.com/joshvanl/bingo/interpreter/ast/errors"
	"github.com/joshvanl/bingo/interpreter/ast/operators"
)

type Program struct {
	Statements []*Statement
}

type Statement struct {
	Expressions []Expression

	// should be len(Expressions)-1
	Operators []Operator
}

type Expression interface {
	Execute(ch <-chan os.Signal) error
	Stdout() io.ReadCloser
}

type Operator func() error
type stringParseOp func(Expression, []string) (func() error, error)

func Parse(in *string) (*Program, error) {
	stmtsStr := strings.Split(*in, ";")

	var stmts []*Statement
	for _, s := range stmtsStr {
		stmt, err := parseStatement(s)
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}

	return &Program{stmts}, nil
}

func parseStatement(stmtStr string) (*Statement, error) {
	ss := strings.Fields(stmtStr)

	var opsS []stringParseOp
	expsS := make([][]string, 1)

	expI := 0

	for _, s := range ss {
		o := toOperator(s)

		if o != nil {
			if len(opsS) >= len(expsS) {
				return nil, errors.BadOperator
			}

			opsS = append(opsS, o)
			expI++

			continue
		}

		if expI >= len(expsS) {
			expsS = append(expsS, []string{})
		}

		expsS[expI] = append(expsS[expI], s)
	}

	if len(opsS) >= len(expsS) {
		return nil, errors.MissingExpression
	}

	exps := make([]Expression, len(expsS))
	for i, e := range expsS {
		exps[i] = command.NewBin(e[0], e[1:])
	}

	var err error
	ops := make([]Operator, len(opsS))
	for i, o := range opsS {
		ops[i], err = o(exps[i], expsS[i+1])
		if err != nil {
			return nil, err
		}
	}

	return &Statement{
		exps,
		ops,
	}, nil
}

func toOperator(o string) stringParseOp {
	switch o {
	case ">":
		return func(cmd Expression, args []string) (func() error, error) {
			return operators.Print(cmd.Stdout(), args)
		}
		//r, w := io.Pipe()
		//f, err := operators.Print(r, "")
		//return &Operator{
		//	F:  f,
		//	WC: w,
		//}, err
	case "&>":
		return nil
	case ">>":
		return nil
	case "&>>":
		return nil
	case "&&":
		return nil
	case "|":
		return nil
	default:
		return nil
	}
}
