package ast

import (
	"io"
	"os"

	"github.com/joshvanl/bingo/command"
)

type cmd struct {
	args []string
	cmd  *command.Cmd
}

var _ Expression = &cmd{}

func (c *cmd) Run(ch <-chan os.Signal) error {
	err := c.cmd.Execute(ch)
	return err
}

func (c *cmd) prepare(in, inerr io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	c.cmd = command.NewBin(c.args[0], c.args[1:], in)
	return c.cmd.Stdout(), c.cmd.Stderr(), nil
}

func (c *cmd) nextToken(token string) bool {
	if isOperator(token) {
		return false
	}

	c.args = append(c.args, token)
	return true
}
