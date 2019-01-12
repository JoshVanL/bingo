package command

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/joshvanl/bingo/command/builtin"
)

type Command struct {
	cmdF func(ch <-chan os.Signal) error
}

func New(cmd string, args []string) *Command {
	var cmdF func(ch <-chan os.Signal) error

	switch cmd {
	case "cd":
		cmdF = func(ch <-chan os.Signal) error {
			return builtin.Cd(args)
		}

	default:

		cmd := exec.Command(cmd, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

		cmdF = func(ch <-chan os.Signal) error {
			if err := cmd.Start(); err != nil {
				return err
			}

			done := make(chan struct{})
			defer close(done)

			go func() {
				for {
					select {
					case sig := <-ch:
						s, ok := sig.(syscall.Signal)
						if !ok {
							continue
						}

						syscall.Kill(cmd.Process.Pid, s)
						continue

					case <-done:
						break
					}
				}
			}()

			err := cmd.Wait()
			_, ok := err.(*exec.ExitError)
			if !ok {
				return err
			}

			return nil
		}
	}

	return &Command{cmdF}
}

func (c *Command) Execute(ch <-chan os.Signal) error {
	return c.cmdF(ch)
}
