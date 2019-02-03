package ast

import (
	"io"
	"os"

	"github.com/joshvanl/bingo/interpreter/ast/operators"
)

var _ Expression = &oPrint{}
var _ Expression = &oAppend{}

type oPrint struct {
	args []string
	f    func(<-chan os.Signal) error
	stop chan struct{}
}

type oAppend struct {
	args []string
	f    func(<-chan os.Signal) error
	stop chan struct{}
}

func (o *oPrint) Run(ch <-chan os.Signal) error {
	return o.f(ch)
}

func (o *oPrint) Stop() {
	close(o.stop)
}

func (o *oAppend) Run(ch <-chan os.Signal) error {
	return o.f(ch)
}

func (o *oAppend) Stop() {
	close(o.stop)
}

func (o *oPrint) nextToken(token string) bool {
	o.args = append(o.args, token)
	return true
}

func (o *oAppend) nextToken(token string) bool {
	o.args = append(o.args, token)
	return true
}

func (o *oPrint) prepare(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	o.stop = make(chan struct{})
	f, err := operators.Print(in, o.args, o.stop)
	if err != nil {
		return nil, inerr, err
	}

	o.f = f
	return nil, inerr, nil
}

func (o *oAppend) prepare(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	o.stop = make(chan struct{})
	f, err := operators.Append(in, o.args, o.stop)
	if err != nil {
		return nil, inerr, err
	}

	o.f = f
	return nil, inerr, nil
}
