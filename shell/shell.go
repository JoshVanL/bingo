package shell

import (
	"bufio"
	"io"
	"os"
	"syscall"

	"github.com/joshvanl/bingo/interpreter"
	"github.com/joshvanl/bingo/prompt"
	"github.com/joshvanl/bingo/utils"
)

type Shell struct {
	prompt *prompt.Prompt

	in       *os.File
	out, err *bufio.Writer

	sig        <-chan os.Signal
	readCh     <-chan *string
	holdReader chan struct{}
}

func New() *Shell {
	s := &Shell{
		prompt:     prompt.New(),
		in:         os.Stdin,
		out:        bufio.NewWriter(os.Stdout),
		err:        bufio.NewWriter(os.Stderr),
		holdReader: make(chan struct{}),
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

LOOP:
	for {
		select {
		case sig := <-s.sig:

			if sig == syscall.SIGCHLD {
				continue LOOP
			}

			s.output("\n")
			return

		case pi := <-s.readCh:
			i = *pi
			break LOOP
		}
	}

	if i[0] == '\n' {
		i = i[1:]
	}

	if len(i) == 0 {
		s.holdReader <- struct{}{}
		return
	}

	if i[len(i)-1] == '\n' {
		i = i[:len(i)-1]
	}

	f := interpreter.Parse(&i)
	s.must(f(s.sig))

	for len(s.sig) > 0 {
		<-s.sig
	}

	s.holdReader <- struct{}{}
}

func (s *Shell) listenStdin() chan *string {
	ch := make(chan *string)

	go func() {
		for {
			reader := bufio.NewReader(s.in)
			i, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					os.Exit(0)
				}

				s.must(err)
				continue
			}

			ch <- &i
			<-s.holdReader
		}
	}()

	return ch
}

func (s *Shell) must(err error) {
	if err != nil {
		s.err.WriteString("bingo error: " + err.Error() + "\n")
		s.err.Flush()
	}
}

func (s *Shell) output(os ...string) {
	_, err := s.out.WriteString(utils.Join(os))
	s.must(err)
	s.must(s.out.Flush())
}
