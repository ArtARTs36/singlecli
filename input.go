package cli

import (
	"os"
	"syscall"

	"golang.org/x/term"
)

type Input interface {
	ReadPassword() (string, error)
}

type input struct {
}

func (i *input) ReadPassword() (string, error) {
	// based on https://github.com/golang/go/issues/19909
	var fd int
	if term.IsTerminal(syscall.Stdin) {
		fd = syscall.Stdin
	} else {
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
