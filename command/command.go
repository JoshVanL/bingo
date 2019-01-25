package command

import (
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/joshvanl/bingo/command/builtin"
)

type Cmd struct {
	f      func(ch <-chan os.Signal) error
	stdout io.WriteCloser
}

func NewBin(cmd string, args []string) *Cmd {
	command := new(Cmd)

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

		if command.stdout != nil {
			cmd.Stdout = command.stdout
		} else {
			cmd.Stdout = os.Stdout
		}

		cmd.Stderr = os.Stderr

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Noctty:  true,
		}

		command.f = func(ch <-chan os.Signal) error {

			//if err := cmd.Process.Release(); err != nil {
			//	return err
			//}

			done := make(chan struct{})
			defer close(done)

			if err := cmd.Start(); err != nil {
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
			if err == nil {
				return err
			}

			_, ok := err.(*exec.ExitError)
			if !ok {
				return err
			}

			return nil
		}
	}

	return command
}

func (c *Cmd) Stdout() io.ReadCloser {
	r, w := io.Pipe()
	c.stdout = w
	return r
}

func (c *Cmd) Execute(ch <-chan os.Signal) error {
	return c.f(ch)
}
