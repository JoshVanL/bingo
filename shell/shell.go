package shell

import (
	"bufio"
	"io"
	"os"
	"syscall"

	"golang.org/x/sys/unix"

	"github.com/joshvanl/bingo/interpreter"
	"github.com/joshvanl/bingo/prompt"
	"github.com/joshvanl/bingo/shell/terminal"
)

const (
	keyEscape = 27

	ioctlReadTermios  = unix.TCGETS
	ioctlWriteTermios = unix.TCSETS
)

type Shell struct {
	prompt *prompt.Prompt

	in       *os.File
	out, err *bufio.Writer

	currState, oldState unix.Termios
	fd                  int

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

	s.fd = int(os.Stdin.Fd())
	termios, err := unix.IoctlGetTermios(s.fd, ioctlReadTermios)
	s.must(err)

	s.readCh = s.listenStdin()

	s.oldState = *termios

	termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	s.currState = *termios
	err = unix.IoctlSetTermios(s.fd, ioctlWriteTermios, termios)
	s.must(err)

	return s
}

//func (s *Shell) Prompt() {
//	p, err := s.prompt.String()
//	s.must(err)
//
//	s.output('\r', p...)
//}

func (s *Shell) Run() {
	var i string

LOOP:
	for {
		select {
		case sig := <-s.sig:

			if sig == syscall.SIGCHLD {
				continue LOOP
			}

			//s.output("\r\n")
			return

		case pi := <-s.readCh:
			i = *pi
			break LOOP
		}
	}

	f := interpreter.Parse(&i)
	s.must(unix.IoctlSetTermios(s.fd, ioctlWriteTermios, &s.oldState))
	s.must(f(s.sig))
	s.must(unix.IoctlSetTermios(s.fd, ioctlWriteTermios, &s.currState))

	for len(s.sig) > 0 {
		<-s.sig
	}

	s.holdReader <- struct{}{}
}

func (s *Shell) listenStdin() chan *string {
	ch := make(chan *string)
	//buff := make([]byte, 0, 1024)

	term := terminal.NewTerminal(s.in, s.prompt.String)

	go func() {
		for {
			line, err := term.ReadLine()
			if err != nil {
				if err == io.EOF {
					s.die(0)
				}
				s.must(err)
			}

			ch <- &line
			<-s.holdReader
		}
	}()

	//go func() {
	//	for {
	//		b := make([]byte, 1)
	//		_, err := s.in.Read(b)
	//		if err != nil {
	//			if err == io.EOF {
	//				s.die(0)
	//			}
	//			s.must(err)
	//		}

	//		os.Stdout.Write(b)
	//		buff = append(buff, b[0])

	//		switch b[0] {
	//		case 4:
	//			s.die(0)

	//		case 127:
	//			os.Stdout.Write([]byte{keyEscape, '[', 'D'})

	//		case '\r':
	//			os.Stdout.Write([]byte{'\n', '\r'})
	//			i := string(buff[:len(buff)])
	//			ch <- &i
	//			buff = make([]byte, 0, 1024)
	//			<-s.holdReader
	//		}
	//	}
	//}()

	return ch
}

func (s *Shell) must(err error) {
	if err != nil {
		s.err.WriteString("bingo: " + err.Error() + "\n\r")
		s.err.Flush()
	}
}

func (s *Shell) output(os ...rune) {
	_, err := s.out.WriteString(string(os))
	s.must(err)
	s.must(s.out.Flush())
}

func (s *Shell) die(exitCode int) {
	unix.IoctlSetTermios(int(os.Stdin.Fd()), ioctlWriteTermios, &s.oldState)
	os.Exit(exitCode)
}
