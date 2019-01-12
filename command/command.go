package command

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/joshvanl/bingo/command/builtin"
)

func Execute(cmd string, args []string) error {
	switch cmd {
	case "cd":
		return builtin.Cd(args)
	default:

		cmd := exec.Command(cmd, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

		if err := cmd.Start(); err != nil {
			return err
		}

		return cmd.Wait()

	}
}
