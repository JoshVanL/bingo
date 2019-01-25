package interpreter

import (
	"os"
	"sync"

	"github.com/joshvanl/bingo/interpreter/ast"
	"github.com/joshvanl/bingo/interpreter/ast/errors"
)

func Run(stmt *ast.Statement, ch <-chan os.Signal) error {
	var errLock sync.Mutex
	var wg sync.WaitGroup
	var result []error

	for i := 0; i < len(stmt.Expressions); i++ {

		if i < len(stmt.Operators) {
			wg.Add(1)
			go func() {
				err := stmt.Operators[i]()
				if err != nil {

					errLock.Lock()
					result = append(result, err)
					errLock.Unlock()

				}

				wg.Done()
			}()
		}

		err := stmt.Expressions[i].Execute(ch)
		if err != nil {

			errLock.Lock()
			result = append(result, err)
			errLock.Unlock()

		}

		wg.Wait()
		if result != nil {
			return errors.Join(result)
		}
	}

	return nil
}
