package interpreter

import (
	"io"
	"os"
	"sync"

	"github.com/joshvanl/bingo/interpreter/ast"
	"github.com/joshvanl/bingo/interpreter/ast/errors"
)

func Run(stmt *ast.Statement, ch <-chan os.Signal) error {
	var errLock sync.Mutex
	var wg sync.WaitGroup
	var result []error

	wg.Add(len(stmt.Expressions) + 2)

	go func() {
		if stmt.Out != nil {
			io.Copy(os.Stdout, stmt.Out)
			stmt.Out.Close()
		}
		wg.Done()
	}()

	go func() {
		if stmt.Err != nil {
			io.Copy(os.Stderr, stmt.Err)
			stmt.Err.Close()
		}
		wg.Done()
	}()

	var stop bool
	stopAll := func() {
		if stop {
			return
		}

		stop = true

		for i := 0; i < len(stmt.Expressions); i++ {
			stmt.Expressions[i].Stop()
		}
	}

	for i := 0; i < len(stmt.Expressions); i++ {
		go func(i int) {
			err := stmt.Expressions[i].Run(ch)
			if err != nil {

				errLock.Lock()
				result = append(result, err)
				stopAll()
				errLock.Unlock()

				if stmt.Err != nil {
					stmt.Err.Close()
				}

				if stmt.Out != nil {
					stmt.Out.Close()
				}

			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	return errors.Join(result)
}
