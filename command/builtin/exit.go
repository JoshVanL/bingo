package builtin

import (
	"fmt"
	"os"
	"strconv"
)

// we need to fix this since we need to restore the terminal state properly
func Exit(args []string) error {
	if len(args) == 0 {
		os.Exit(0)
	}

	if len(args) > 1 {
		return fmt.Errorf("expecting single integer exit code, got=%d", len(args))
	}

	n, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	os.Exit(n)

	return nil
}
