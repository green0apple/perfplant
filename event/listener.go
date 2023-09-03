package event

import (
	"errors"
	"syscall"
)

type ListenOption int

const (
	LISTEN_OPT_NONBLOCK = iota + 1
	LISTEN_OPT_REUSEADDR
	LISTEN_OPT_REUSEPORT
)

var (
	ErrInvalidListenerOption = errors.New("invalid listener option")
)

func setListenerOpt(fd int, options ...ListenOption) error {
	var err error

	for _, opt := range options {
		switch opt {
		case LISTEN_OPT_NONBLOCK:
			err = syscall.SetNonblock(fd, true)
		case LISTEN_OPT_REUSEADDR:
			err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.SO_REUSEADDR, 1)
		case LISTEN_OPT_REUSEPORT:
			// TODO :: REUSEPORT?
			//err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.SO_PORT, 1)

		default:
			err = ErrInvalidListenerOption
		}

		if err != nil {
			return err
		}
	}

	return nil
}

type Listener Conn

func NewListener() *Listener {
	return (*Listener)(NewConn())
}

func (l *Listener) ListenUDP(addr syscall.SockaddrInet4, backlog int, options ...ListenOption) error {
	return l.listen(addr, backlog, syscall.SOCK_DGRAM, options...)
}

func (l *Listener) listen(addr syscall.SockaddrInet4, backlog int, sockType int, options ...ListenOption) error {
	fd, err := syscall.Socket(syscall.AF_INET, sockType, syscall.IPPROTO_IP)
	if err != nil {
		return err
	}

	if err = setListenerOpt(fd, options...); err != nil {
		syscall.Close(fd)
		return err
	}

	if err = syscall.Bind(fd, &addr); err != nil {
		syscall.Close(fd)
		return err
	}

	if err = syscall.Listen(fd, backlog); err != nil {
		syscall.Close(fd)
		return err
	}

	(*l).SetFd(fd)
	l.SetType(CONN_TYPE_LISTENER)

	// TODO :: SetAddr
	l.conn.saddr = &addr
	return nil
}
