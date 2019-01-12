package shell

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joshvanl/bingo/interpreter"
	"github.com/joshvanl/bingo/prompt"
	"github.com/joshvanl/bingo/utils"
)

type Shell struct {
	prompt *prompt.Prompt

	in, out, err *os.File
	reader       *bufio.Reader
}

func New() *Shell {
	return &Shell{
		prompt: prompt.New(),
		in:     os.Stdin,
		out:    os.Stdout,
		err:    os.Stderr,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (s *Shell) Prompt() {
	p, err := s.prompt.String()
	s.must(err)

	s.output(p)
}

func (s *Shell) Run() {
	i, err := s.reader.ReadString('\n')
	s.must(err)

	i = i[:len(i)-1]
	if err != nil || len(i) == 0 {
		return
	}

	s.must(interpreter.Run(&i))
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
