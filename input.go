package cli

import (
	"golang.org/x/term"
	"os"
)

type Input interface {
	ReadPassword() (string, error)
}

type input struct {
}

func (i *input) ReadPassword() (string, error) {
	// based on https://github.com/golang/go/issues/19909
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		tty, err := os.Open("/dev/tty")
		if err != nil {
			return "", err
		}

		defer tty.Close()
		fd = int(tty.Fd())
	}

	pass, err := term.ReadPassword(fd)
	return string(pass), err
}
