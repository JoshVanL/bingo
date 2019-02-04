package operators

import (
	"io"
	"os"
)

func Pipe(in io.ReadCloser, stop <-chan struct{}) (func(<-chan os.Signal) error, io.ReadCloser, error) {
	r, w := io.Pipe()
	done := make(chan struct{})

	return func(ch <-chan os.Signal) error {
		go func() {
			select {
			case <-done:
				w.Close()
			case <-stop:
			case <-ch:
			}
			//w.Close()
			//r.Close()
		}()

		io.Copy(w, in)

		close(done)

		return nil
	}, r, nil
}
