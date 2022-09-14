package game

import "golang.org/x/sys/unix"

const (
	tcGetRequest = unix.TCGETS
	tcSetRequest = unix.TCSETS
)
