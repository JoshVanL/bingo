package ast

import (
	"errors"
	"strings"
)

var (
	badOperator  = errors.New("could not parse operator")
	moreOperator = errors.New("parse error near '\\n'")
)

var operators = []string{
	">",
	"&>",
	"&>>",
	"&&",
	"|",
}

type Program struct {
	Statements []*Statement
}

type Statement struct {
	Expressions []Expression

	// should be len(Expressions)-1
	Operators []Operator
}

//type Expression struct {
//	Command string
//}

type Expression []string

type Operator string

const Print Operator = ">"
const PrintStderr Operator = "&>"
const Append Operator = ">>"
const AppendStderr Operator = "&>>"
const And Operator = "&&"
const Pipe Operator = "|"
const OperatorInValid Operator = ""

func Parse(in string) (*Program, error) {
	stmtsStr := strings.Split(in, ";")

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

	var ops []Operator
	exps := make([]Expression, 1)

	expI := 0

	for _, s := range ss {
		o := toOpoerator(s)
		if o != OperatorInValid {
			if len(ops) >= len(exps) {
				return nil, badOperator
			}

			ops = append(ops, o)
			expI++

			continue
		}

		if expI >= len(exps) {
			exps = append(exps, Expression{})
		}

		exps[expI] = append(exps[expI], s)
	}

	if len(ops) >= len(exps) {
		return nil, moreOperator
	}

	return &Statement{
		exps,
		ops,
	}, nil
}

func toOpoerator(o string) Operator {
	switch o {
	case ">":
		return Print
	case "&>":
		return PrintStderr
	case ">>":
		return Append
	case "&>>":
		return AppendStderr
	case "&&":
		return And
	case "|":
		return Pipe
	default:
		return OperatorInValid
	}
}
