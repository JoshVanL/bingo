package shell

import (
	"bufio"
	"io"
	"os"
	"syscall"

	"github.com/joshvanl/bingo/interpreter"
	"github.com/joshvanl/bingo/prompt"
	"golang.org/x/sys/unix"

	"github.com/joshvanl/bingo/utils"
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

	oldState *unix.Termios

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

	fd := int(os.Stdin.Fd())
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	s.must(err)

	oldState := termios
	s.oldState = oldState

	termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	err = unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)
	s.must(err)

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

	for len(i) > 0 && (i[0] == '\r' || i[0] == '\n') {
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
	buff := make([]byte, 0, 1024)

	go func() {
		for {
			b := make([]byte, 1)
			_, err := s.in.Read(b)
			if err != nil {
				if err == io.EOF {
					s.die(0)
				}
				s.must(err)
			}

			os.Stdout.Write(b)
			buff = append(buff, b[0])

			//if b[0] == '\r' {
			//	continue
			//}

			if b[0] == 13 {

				os.Stdout.Write([]byte{keyEscape, '[', 'B'})
				//os.Stdout.Write([]byte{'\r', '\n'})
				i := string(buff[:len(buff)])
				ch <- &i
				buff = make([]byte, 0, 1024)
				<-s.holdReader

				continue
			}

			//io.Copy(s.out, s.in)
			//r, _, err := reader.ReadRune()
			//s.must(err)
			//fmt.Printf("here %s\n", r)

			//if r != 27 {
			//	s.out.WriteRune(r)
			//}

			//fmt.Printf(">%s\n", r)
			//i, err := reader.ReadString('\n')
			//if err != nil {
			//	if err == io.EOF {
			//		os.Exit(0)
			//	}

			//	s.must(err)
			//	continue
			//}

			//ch <- &i
			//<-s.holdReader
		}
	}()

	return ch
}

func (s *Shell) must(err error) {
	if err != nil {
		s.err.WriteString("bingo: " + err.Error() + "\n")
		s.err.Flush()
	}
}

func (s *Shell) output(os ...string) {
	_, err := s.out.WriteString(utils.Join(os))
	s.must(err)
	s.must(s.out.Flush())
}

func (s *Shell) die(exitCode int) {
	unix.IoctlSetTermios(int(os.Stdin.Fd()), ioctlWriteTermios, s.oldState)
	os.Exit(exitCode)
}
