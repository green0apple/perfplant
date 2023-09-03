package event

import "syscall"

const (
	CONN_TYPE_CONN = iota + 1
	CONN_TYPE_LISTENER
)

type Conn struct {
	fd int

	typ int

	saddr *syscall.SockaddrInet4
	daddr *syscall.SockaddrInet4
}

func NewConn() *Conn {
	return &Conn{}
}

func (c *Conn) SetFd(fd int) { c.fd = fd }
func (c *Conn) GetFd() int   { return c.fd }

func (c *Conn) SetType(typ int)      { c.typ = typ }
func (c *Conn) GetType() int         { return c.typ }
func (c *Conn) IsDefaultConn() bool  { return c.typ == CONN_TYPE_CONN }
func (c *Conn) IsListenerConn() bool { return c.typ == CONN_TYPE_LISTENER }
