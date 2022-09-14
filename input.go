//go:build linux || darwin

package game

import (
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
)

const (
	SysUp = 1000 + iota
	SysDown
	SysLeft
	SysRight
	SysParse
)

const (
	brkInt = unix.BRKINT
	ixon   = unix.IXON
	echo   = unix.ECHO
	icanon = unix.ICANON
	isig   = unix.ISIG
	iexten = unix.IEXTEN
	cSize  = unix.CSIZE
	parenb = unix.PARENB
	cs8    = unix.CS8
	vmin   = unix.VMIN
	vtime  = unix.VTIME
)

// GetCh get char from tty
func GetCh() (int, string) {
	var (
		in              int
		err             error
		sigio           = make(chan os.Signal)
		originalTermios *unix.Termios
	)

	if in, err = unix.Open("/dev/tty", unix.O_RDONLY, 0); err != nil {
		panic(err)
	}

	signal.Notify(sigio, unix.SIGIO)

	_, _ = unix.FcntlInt(uintptr(in), unix.F_SETFL, unix.O_ASYNC|unix.O_NONBLOCK)
	originalTermios, err = unix.IoctlGetTermios(in, tcGetRequest)

	tios := *originalTermios
	tios.Iflag &^= brkInt | ixon
	tios.Lflag &^= echo | icanon | isig | iexten
	tios.Cflag &^= cSize | parenb
	tios.Cflag |= cs8
	tios.Cc[vmin] = 1
	tios.Cc[vtime] = 0

	err = unix.IoctlSetTermios(in, tcSetRequest, &tios)
	defer func() {
		err = unix.IoctlSetTermios(in, tcSetRequest, originalTermios)
		_ = unix.Close(in)
	}()
	buf := make([]byte, 128)

LOOP:
	<-sigio
	n, err := unix.Read(in, buf)
	if err == unix.EAGAIN || err == unix.EWOULDBLOCK {
		goto LOOP
	}
	if n == 1 {
		return int(buf[0]), ""
	} else if n == 3 {
		switch buf[2] {
		case 65:
			return SysUp, ""
		case 66:
			return SysDown, ""
		case 67:
			return SysRight, ""
		case 68:
			return SysLeft, ""
		default:
			return SysParse, string(buf[0:n])
		}
	} else {
		return SysParse, string(buf[0:n])
	}
	return 0, ""
}
