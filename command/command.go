package command

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/joshvanl/bingo/command/builtin"
)

type Cmd struct {
	f              func(ch <-chan os.Signal) error
	stop           chan struct{}
	stdout, stderr *os.File
}

func NewBin(cmd string, args []string, in *os.File) (*Cmd, error) {
	rout, wout, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	rerr, werr, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	command := &Cmd{
		stop:   make(chan struct{}),
		stdout: rout,
		stderr: rerr,
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
		cmd.Stdout = wout
		cmd.Stderr = werr

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Noctty:  true,
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
			wout.Close()
			werr.Close()

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

func (c *Cmd) Stdout() *os.File {
	return c.stdout
}

func (c *Cmd) Stderr() *os.File {
	return c.stderr
}
