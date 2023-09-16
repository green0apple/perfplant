package event

import (
	"errors"
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

var (
	ErrInvalidConnFd    = errors.New("invalid fd")
	ErrOnlySupportInet4 = errors.New("only supports ipv4")
)

func resolveUDP() (int, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return -1, err
	}

	if err = syscall.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		syscall.Close(fd)
		return -1, err
	}

	if err = syscall.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		syscall.Close(fd)
		return -1, err
	}

	if err = syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_RECVERR, 1); err != nil {
		syscall.Close(fd)
		return -1, err
	}

	if err = syscall.SetNonblock(fd, true); err != nil {
		syscall.Close(fd)
		return -1, err
	}

	if err = syscall.Bind(fd, &syscall.SockaddrInet4{}); err != nil {
		syscall.Close(fd)
		return -1, err
	}

	return fd, nil
}

func Sockaddr4(sa syscall.Sockaddr) (*syscall.SockaddrInet4, error) {
	switch saddr := sa.(type) {
	case *syscall.SockaddrInet4:
		return saddr, nil
	default:
		return nil, ErrOnlySupportInet4
	}
}

type conn struct {
	fd int32

	typ int

	saddr *syscall.SockaddrInet4
	daddr *syscall.SockaddrInet4

	saddrString string
	daddrString string
}

func newConn() *conn {
	return &conn{}
}

func (c *conn) SetFd(fd int32)                { c.fd = fd }
func (c *conn) Fd() int32                     { return c.fd }
func (c *conn) IsValid() bool                 { return c.fd > 0 }
func (c *conn) Hash() uint32                  { return HashAddr(c.saddr, c.daddr) }
func (c *conn) SAddr() *syscall.SockaddrInet4 { return c.saddr }
func (c *conn) DAddr() *syscall.SockaddrInet4 { return c.daddr }
func (c *conn) SAddrString() string           { return c.saddrString }
func (c *conn) DAddrString() string           { return c.daddrString }

func (c *conn) SetSAddr(saddr *syscall.SockaddrInet4) {
	c.saddr = saddr
	c.saddrString = fmt.Sprintf("%d.%d.%d.%d:%d", saddr.Addr[0], saddr.Addr[1], saddr.Addr[2], saddr.Addr[3], saddr.Port)
}

func (c *conn) SetDAddr(daddr *syscall.SockaddrInet4) {
	c.daddr = daddr
	c.daddrString = fmt.Sprintf("%d.%d.%d.%d:%d", daddr.Addr[0], daddr.Addr[1], daddr.Addr[2], daddr.Addr[3], daddr.Port)
}

func (c *conn) Close() {
	if c.fd > 0 {
		syscall.Close(int(c.fd))
	}
}

type UDPConn struct {
	*conn
}

func NewUDPConn() *UDPConn {
	return &UDPConn{conn: newConn()}
}

func (uc *UDPConn) Dial(daddr *syscall.SockaddrInet4) error {
	fd, err := resolveUDP()
	if err != nil {
		return err
	}

	if err = syscall.Connect(fd, daddr); err != nil {
		syscall.Close(fd)
		return err
	}

	sa, err := syscall.Getsockname(fd)
	if err != nil {
		syscall.Close(fd)
		return err
	}

	saddr, err := Sockaddr4(sa)
	if err != nil {
		return err
	}

	uc.conn.SetFd(int32(fd))
	uc.conn.SetSAddr(saddr)
	uc.conn.SetDAddr(daddr)

	return nil
}

func (uc *UDPConn) Recvmsg() (*syscall.SockaddrInet4, []byte, error) {
	if !uc.IsValid() {
		return nil, nil, ErrInvalidConnFd
	}

	var b []byte = make([]byte, 32767)
	n, _, _, from, err := syscall.Recvmsg(int(uc.fd), b, nil, syscall.MSG_DONTWAIT)

	if err != nil {
		return nil, nil, err
	}

	saddr, err := Sockaddr4(from)
	if err != nil {
		return nil, nil, err
	}

	return saddr, b[:n], nil

}

func (uc *UDPConn) Sendto(to *syscall.SockaddrInet4, b []byte) error {
	if !uc.IsValid() {
		return ErrInvalidConnFd
	}

	return syscall.Sendto(int(uc.fd), b, syscall.MSG_DONTWAIT, to)
}
