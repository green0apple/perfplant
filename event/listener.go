package event

import (
	"errors"
	"perfplant/buffer/rbtree"
	"sync"
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

type ListenerConnectionsTree struct {
	sync.RWMutex
	rbtree.Tree
}

func (this *ListenerConnectionsTree) InsertConn(conn *UDPConn) {
	this.Lock()
	defer this.Unlock()

}

func (this *ListenerConnectionsTree) lookupConn(raddr *syscall.SockaddrInet4) *UDPConn {

}

type Listener interface {
	Listen(addr syscall.Sockaddr, options ...ListenOption) error
	Fd() int
	Close() error
}

type UDPListener struct {
	conn        *UDPConn
	connections *ListenerConnectionsTree
}

func (this *UDPListener) Fd() int { return this.conn.fd }
func (this *UDPListener) Close()  { this.conn.Close() }

func (this *UDPListener) Listen(addr syscall.SockaddrInet4, backlog int, options ...ListenOption) error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
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

	if this.conn == nil {
		this.conn = NewUDPConn()
	}

	this.conn.SetFd(fd)
	this.conn.SetSAddr(&addr)
	return nil
}

func (this *UDPListener) Recvmsg() (*UDPConn, []byte, error) {
	if this.conn.IsValid() {
		return nil, nil, ErrInvalidConnFd
	}

	from, b, err := this.conn.Recvmsg()
	if err != nil {
		return nil, nil, err
	}

	conn, err := this.connections.Search(this.conn.saddr, from)
	if err != nil {
		return nil, nil, err
	}

	if conn == nil {
		conn = NewUDPConn()
	}

	return conn, b, nil
}

func ListenUDP(addr syscall.SockaddrInet4, backlog int, options ...ListenOption) (*UDPListener, error) {
	u := UDPListener{}

	if err := u.Listen(addr, backlog, options...); err != nil {
		return nil, err
	}

	return &u, nil
}
