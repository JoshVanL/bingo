package builtin

import (
	"os"
)

// we need to fix this since we need to restore the terminal state properly
func Exit() {
	os.Exit(0)
}
