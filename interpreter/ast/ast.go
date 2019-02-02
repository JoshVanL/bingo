package ast

import (
	"io"
	"os"
	"strings"
)

type Program struct {
	Statements []*Statement
}

type Statement struct {
	Expressions []Expression
	Out, Err    io.ReadCloser
}

type Expression interface {
	Run(<-chan os.Signal) error
	nextToken(string) bool
	prepare(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error)
}

type operator interface {
	nextToken(string) bool
	prepare(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error)
	Run(<-chan os.Signal) error
}

func (s *Statement) Prepare(in io.ReadCloser, out, serr io.WriteCloser) error {
	var inerr io.ReadCloser

	for _, e := range s.Expressions {
		fin, ferr, err := e.prepare(in, inerr)
		if err != nil {
			return err
		}

		in = fin

		inerr = ferr
	}

	s.Out = in
	s.Err = inerr

	return nil
}

func Parse(in *string) (*Program, error) {
	p := new(Program)

	stmtsStr := strings.Split(*in, ";")
	for _, s := range stmtsStr {
		p.parseStmt(s)
	}

	return &Program{p.Statements}, nil
}

func (p *Program) parseStmt(stmtStr string) {
	tokens := strings.Fields(stmtStr)
	s := new(Statement)

	for len(tokens) != 0 {
		var exp Expression
		exp, tokens = s.parseExpression(tokens)

		s.Expressions = append(s.Expressions, exp)
	}

	p.Statements = append(p.Statements, s)
}

func (s *Statement) parseExpression(tokens []string) (Expression, []string) {
	var i int

	if o := toOperator(tokens[0]); o != nil {
		for i = range tokens[1:] {

			if !o.nextToken(tokens[i+1]) {
				return o, tokens[i+1:]
			}
		}

		return o, nil
	}

	cmd := new(cmd)
	for i = range tokens {
		if !cmd.nextToken(tokens[i]) {
			return cmd, tokens[i:]
		}
	}

	return cmd, nil
}

func toOperator(token string) operator {
	switch token {
	case ">":
		return new(oPrint)
	case "&>":
		return nil
	case ">>":
		return new(oAppend)
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

func isOperator(token string) bool {
	switch token {
	case ">":
		return true
	case "&>":
		return false
	case ">>":
		return true
	case "&>>":
		return false
	case "&&":
		return false
	case "|":
		return false
	default:
		return false
	}
}
