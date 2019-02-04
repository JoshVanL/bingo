package command

import (
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/joshvanl/bingo/command/builtin"
)

type Cmd struct {
	f              func(ch <-chan os.Signal) error
	stop           chan struct{}
	stdout, stderr io.ReadCloser
}

func NewBin(cmd string, args []string, in io.ReadCloser) (*Cmd, error) {
	command := &Cmd{
		stop: make(chan struct{}),
	}

	switch cmd {
	case "cd":
		command.f = func(ch <-chan os.Signal) error {
			return builtin.Cd(args)
		}

	case "exit":
		command.f = func(ch <-chan os.Signal) error {
			return builtin.Exit(args)
		}

	default:

		cmd := exec.Command(cmd, args...)
		cmd.Stdin = in

		rerr, err := cmd.StderrPipe()
		if err != nil {
			return nil, err
		}

		rout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}

		command.stdout = rout
		command.stderr = rerr

		// TODO: do these options
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

		command.f = func(ch <-chan os.Signal) error {

			done := make(chan struct{})

			if err := cmd.Start(); err != nil {
				close(done)
				return err
			}

			go func() {
				for {
					select {
					case <-done:
						return

					case sig := <-ch:
						s, ok := sig.(syscall.Signal)
						if !ok {
							continue
						}

						if s == syscall.SIGCHLD {
							return
						}

						syscall.Kill(cmd.Process.Pid, s)
						continue
					}
				}
			}()

			err := cmd.Wait()
			//command.stderr.Close()
			//command.stdout.Close()

			close(done)

			if _, ok := err.(*exec.ExitError); ok {
				return nil
			}

			return err
		}
	}

	return command, nil
}

func (c *Cmd) Execute(ch <-chan os.Signal) error {
	return c.f(ch)
}

func (c *Cmd) Stop() {
	close(c.stop)
}

func (c *Cmd) Stdout() io.ReadCloser {
	return c.stdout
}

func (c *Cmd) Stderr() io.ReadCloser {
	return c.stderr
}
