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

	wg.Add(len(stmt.Expressions))

	go func() {
		io.Copy(os.Stdout, stmt.Out)
	}()

	go func() {
		io.Copy(os.Stderr, stmt.Err)
	}()

	for i := 0; i < len(stmt.Expressions); i++ {
		go func(i int) {
			err := stmt.Expressions[i].Run(ch)
			if err != nil {

				errLock.Lock()
				result = append(result, err)
				errLock.Unlock()

			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	return errors.Join(result)
}
