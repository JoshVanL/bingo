package ast

import (
	"io"
	"os"

	"github.com/joshvanl/bingo/interpreter/ast/operators"
)

var _ Expression = &oPrint{}

type oPrint struct {
	args []string
	f    func() error
}

func (o *oPrint) Run(ch <-chan os.Signal) error {
	if err := o.f(); err != nil {
		panic(err)
	}
	return nil
}

func (o *oPrint) nextToken(token string) bool {
	o.args = append(o.args, token)
	return true
}

func (o *oPrint) prepare(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	f, err := operators.Print(in, o.args)
	if err != nil {
		return nil, nil, err
	}

	o.f = f
	return nil, nil, nil
}
