package event

import "syscall"

type Conn struct {
	fd    int
	saddr *syscall.SockaddrInet4
	daddr *syscall.SockaddrInet4
}
