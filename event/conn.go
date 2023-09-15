package event

import (
	"errors"
	"hash/crc32"
	"perfplant/buffer/rbtree"
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

func Hash(saddr, daddr *syscall.SockaddrInet4) uint32 {
	var (
		b     []byte
		addrs [2]*syscall.SockaddrInet4
	)

	if saddr != nil && daddr != nil {
		if saddr.Port < daddr.Port {
			addrs[0] = saddr
			addrs[1] = daddr
		} else {
			addrs[0] = daddr
			addrs[1] = saddr
		}
	} else {
		if saddr != nil {
			addrs[0] = saddr
		}

		if daddr != nil {
			addrs[1] = daddr
		}
	}

	for _, a := range addrs {
		if a == nil {
			continue
		}

		b = append(b, a.Addr[:]...)
		b = append(b, rbtree.PortLittleEndian(a.Port)...)
	}

	return crc32.ChecksumIEEE(b)
}

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
	return Hash(c.saddr, c.daddr)
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

	uc.conn.SetFd(fd)
	uc.conn.SetSAddr(saddr)
	uc.conn.SetDAddr(daddr)

	return nil
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

	return syscall.Sendto(uc.fd, b, syscall.MSG_DONTWAIT, to)
}
