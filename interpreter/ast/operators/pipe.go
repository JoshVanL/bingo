package operators

import (
	"io"
	"os"
)

func Pipe(in *os.File, stop <-chan struct{}) (func(<-chan os.Signal) error, *os.File, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	rp, wp := io.Pipe()
	done := make(chan struct{})

	return func(ch <-chan os.Signal) error {
		go func() {
			select {
			case <-done:
			case <-stop:
			case <-ch:
			}
			wp.Close()
			rp.Close()
			r.Close()
			w.Close()
		}()

		go func() {
			io.Copy(wp, in)
		}()

		io.Copy(w, rp)

		close(done)

		return nil
	}, r, nil
}
