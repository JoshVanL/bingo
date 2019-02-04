package ast

import (
	"fmt"
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
	Stop()
	nextToken(string) bool
	prepare(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error)
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
		err := p.parseStmt(s)
		if err != nil {
			return nil, err
		}
	}

	return &Program{p.Statements}, nil
}

func (p *Program) parseStmt(stmtStr string) error {
	tokens := strings.Fields(stmtStr)
	s := new(Statement)

	if isBadStart(tokens[0]) {
		return fmt.Errorf("parse error near '%v'", tokens[0])
	}

	for len(tokens) != 0 {
		var exp Expression
		exp, tokens = s.parseExpression(tokens)

		s.Expressions = append(s.Expressions, exp)
	}

	p.Statements = append(p.Statements, s)

	return nil
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

func toOperator(token string) *operator {
	switch token {
	case ">":
		return newPrint()
	case "&>":
		return nil
	case ">>":
		return newAppend()
	case "&>>":
		return nil
	case "&&":
		return nil
	case "|":
		return newPipe()
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
		return true
	default:
		return false
	}
}

func isBadStart(token string) bool {
	switch token {
	case "|":
		return true
	default:
		return false
	}
}
