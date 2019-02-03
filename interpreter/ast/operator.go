package ast

import (
	"io"
	"os"

	"github.com/joshvanl/bingo/interpreter/ast/operators"
)

var _ Expression = &operator{}

type operator struct {
	args []string
	stop chan struct{}

	f          func(<-chan os.Signal) error
	prepareF   func(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error)
	nextTokenF func(token string) bool
}

func newPrint() *operator {
	o := new(operator)
	o.prepareF = func(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
		o.stop = make(chan struct{})
		f, err := operators.Print(in, o.args, o.stop)
		if err != nil {
			return nil, inerr, err
		}

		o.f = f
		return nil, inerr, nil
	}

	o.nextTokenF = func(token string) bool {
		o.args = append(o.args, token)
		return true
	}

	return o
}

func newAppend() *operator {
	o := new(operator)
	o.prepareF = func(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
		o.stop = make(chan struct{})
		f, err := operators.Append(in, o.args, o.stop)
		if err != nil {
			return nil, inerr, err
		}

		o.f = f
		return nil, inerr, nil
	}

	o.nextTokenF = func(token string) bool {
		o.args = append(o.args, token)
		return true
	}

	return o
}

func (o *operator) Run(ch <-chan os.Signal) error {
	return o.f(ch)
}

func (o *operator) Stop() {
	close(o.stop)
}

func (o *operator) prepare(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	return o.prepareF(in, inerr)
}

func (o *operator) nextToken(token string) bool {
	return o.nextTokenF(token)
}
