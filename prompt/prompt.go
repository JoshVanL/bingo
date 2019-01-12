package prompt

import (
	"fmt"
	"os"
)

type Prompt struct {
	fs []func() (string, error)

	outputF func(os ...string)
}

func New() *Prompt {
	return &Prompt{
		fs: []func() (string, error){cwd, dollar},
	}
}

func (p *Prompt) String() (string, error) {
	var out string
	var errs error

	for _, f := range p.fs {
		s, err := f()
		if err != nil {
			errs = fmt.Errorf("%s\n%s", errs, err)
			continue
		}

		out = out + s
	}

	return out, errs
}

func dollar() (string, error) {
	return ":$ ", nil
}

func cwd() (string, error) {
	return os.Getwd()
}
