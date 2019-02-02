package ast

import (
	"io"
	"os"

	"github.com/joshvanl/bingo/interpreter/ast/operators"
)

var _ Expression = &oPrint{}

type oPrint struct {
	args []string
	err  io.WriteCloser
	f    func() error
}

func (o *oPrint) Run(ch <-chan os.Signal) error {
	return o.f()
}

func (o *oPrint) nextToken(token string) bool {
	o.args = append(o.args, token)
	return true
}

func (o *oPrint) prepare(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	f, err := operators.Print(in, o.args)
	if err != nil {
		return nil, inerr, err
	}

	o.f = f
	return nil, inerr, nil
}
