package event

import (
	"errors"
	"hash/crc32"
	"perfplant/buffer/rbtree"
	"syscall"
)

var (
	ErrInvalidConnFd    = errors.New("invalid fd")
	ErrOnlySupportInet4 = errors.New("only supports ipv4")
)

type conn struct {
	fd int

	typ int

	saddr *syscall.SockaddrInet4
	daddr *syscall.SockaddrInet4
}

func newConn() *conn {
	return &conn{}
}

func (c *conn) SetFd(fd int)                          { c.fd = fd }
func (c *conn) Fd() int                               { return c.fd }
func (c *conn) IsValid() bool                         { return c.fd > 0 }
func (c *conn) SetSAddr(saddr *syscall.SockaddrInet4) { c.saddr = saddr }
func (c *conn) SetDAddr(daddr *syscall.SockaddrInet4) { c.daddr = daddr }
func (c *conn) SAddr() *syscall.SockaddrInet4         { return c.saddr }
func (c *conn) DAddr() *syscall.SockaddrInet4         { return c.daddr }

func (c *conn) Hash() uint32 {
	var b []byte
	if c.saddr != nil {
		b = append(b, c.saddr.Addr[:]...)
		b = append(b, rbtree.PortLittleEndian(c.saddr.Port)...)
	}

	if c.daddr != nil {
		b = append(b, c.daddr.Addr[:]...)
		b = append(b, rbtree.PortLittleEndian(c.daddr.Port)...)
	}

	return crc32.ChecksumIEEE(b)
}

func (c *conn) Close() {
	if c.fd > 0 {
		syscall.Close(c.fd)
	}
}

type UDPConn struct {
	*conn
}

func NewUDPConn() *UDPConn {
	return &UDPConn{conn: newConn()}
}

func (uc *UDPConn) Recvmsg() (*syscall.SockaddrInet4, []byte, error) {
	if !uc.IsValid() {
		return nil, nil, ErrInvalidConnFd
	}

	var b []byte = make([]byte, 32767)
	n, _, _, from, err := syscall.Recvmsg(uc.fd, b, nil, syscall.MSG_DONTWAIT)

	if err != nil {
		return nil, nil, err
	}

	switch saddr := from.(type) {
	case *syscall.SockaddrInet4:
		return saddr, b[:n], nil
	default:
		err = ErrOnlySupportInet4
		return nil, nil, err
	}
}

func (uc *UDPConn) Sendto(to *syscall.SockaddrInet4, b []byte) error {
	if !uc.IsValid() {
		return ErrInvalidConnFd
	}

	return syscall.Sendto(uc.fd, b, syscall.MSG_DONTWAIT, to)
}
