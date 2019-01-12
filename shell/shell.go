package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/joshvanl/bingo/interpreter"
	"github.com/joshvanl/bingo/prompt"
	"github.com/joshvanl/bingo/utils"
)

type Shell struct {
	prompt *prompt.Prompt

	in, out, err *os.File

	sig    <-chan os.Signal
	readCh <-chan *string

	reader *bufio.Reader
}

func New() *Shell {
	s := &Shell{
		prompt: prompt.New(),
		in:     os.Stdin,
		out:    os.Stdout,
		err:    os.Stderr,
		reader: bufio.NewReader(os.Stdin),
	}

	s.sig = s.signalHandler()
	s.readCh = s.listenStdin()

	return s
}

func (s *Shell) Prompt() {
	p, err := s.prompt.String()
	s.must(err)

	s.output(p)
}

func (s *Shell) Run() {
	var i string

	select {
	case <-s.sig:
		s.output("\n")
		return
	case pi := <-s.readCh:
		i = *pi
	}

	if i[0] == '\n' {
		i = i[:len(i)-1]
	}

	if len(i) == 0 {
		return
	}

	f := interpreter.Parse(&i)
	s.must(f(s.sig))
}

func (s *Shell) listenStdin() chan *string {
	ch := make(chan *string)

	go func() {
		for {
			i, err := s.reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					os.Exit(0)
				}

				s.must(err)
				continue
			}

			ch <- &i
		}
	}()

	return ch
}

func (s *Shell) must(err error) {
	if err != nil {
		fmt.Fprint(s.err, "bingo error: ", err.Error(), "\n")
	}
}

func (s *Shell) output(os ...string) {
	_, err := fmt.Fprint(s.out, utils.Join(os))
	s.must(err)
}
