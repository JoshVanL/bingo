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
	stdout, stderr io.WriteCloser
}

func NewBin(cmd string, args []string, in io.ReadCloser) *Cmd {
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
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Noctty:  true,
		}

		command.f = func(ch <-chan os.Signal) error {

			if command.stdout != nil {
				cmd.Stdout = command.stdout
			} else {
				cmd.Stdout = os.Stdout
			}

			if command.stderr != nil {
				cmd.Stderr = command.stderr
			} else {
				cmd.Stderr = os.Stderr
			}

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

			if command.stdout != nil {
				command.stdout.Close()
			}

			if command.stderr != nil {
				command.stderr.Close()
			}

			close(done)

			return err
		}
	}

	return command
}

func (c *Cmd) Stdout() io.ReadCloser {
	r, w := io.Pipe()
	c.stdout = w
	return r
}

func (c *Cmd) Stderr() io.ReadCloser {
	r, w := io.Pipe()
	c.stderr = w
	return r
}

func (c *Cmd) Execute(ch <-chan os.Signal) error {
	return c.f(ch)
}

func (c *Cmd) Stop() {
	close(c.stop)
}
