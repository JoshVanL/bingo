package prompt

import (
	"fmt"
	"os"
)

type Prompt struct {
	fs []func() ([]rune, error)

	outputF func(os ...string)
}

func New() *Prompt {
	return &Prompt{
		fs: []func() ([]rune, error){cwd, dollar},
	}
}

func (p *Prompt) String() []rune {
	var out []rune
	var errs error

	for _, f := range p.fs {
		s, err := f()
		if err != nil {
			errs = fmt.Errorf("%s\n%s", errs, err)
			continue
		}

		out = append(out, s...)
	}

	if errs != nil {
		fmt.Fprint(os.Stderr, "prompt error: ", errs.Error())
	}

	return out
}

func dollar() ([]rune, error) {
	return []rune(":$ "), nil
}

func cwd() ([]rune, error) {
	wd, err := os.Getwd()
	return []rune(wd), err
}
