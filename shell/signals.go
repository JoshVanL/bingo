package shell

import (
	"os"
	"os/signal"
	"syscall"
)

func (s *Shell) signalHandler() <-chan os.Signal {
	sigCh := make(chan os.Signal, 10)
	signal.Notify(sigCh)
	ch := make(chan os.Signal, 10)

	go func() {
		for {
			sig := <-sigCh
			switch sig {
			case syscall.SIGTERM:
				s.die(0)

			default:
				ch <- sig
			}
		}
	}()

	return ch
}
