package game

import "golang.org/x/sys/unix"

const (
	tcGetRequest = unix.TIOCGETA
	tcSetRequest = unix.TIOCSETA
)
