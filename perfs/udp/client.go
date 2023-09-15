package udp

import (
	"fmt"
	"perfplant/buffer/rbtree"
	"perfplant/event"
	"perfplant/event/module/epoll"
	"sync"
	"syscall"
	"unsafe"
)

type clientConnectionsTree struct {
	sync.RWMutex

	tree *rbtree.Rbtree
}

func (t *clientConnectionsTree) InsertConn(conn *event.UDPConn) {
	t.Lock()
	defer t.Unlock()
	t.tree.Insert(conn.Hash(), unsafe.Pointer(conn))
}

func (t *clientConnectionsTree) DeleteConn(conn *event.UDPConn) {
	t.Lock()
	defer t.Unlock()
	t.tree.Delete(conn.Hash())
}

func (t *clientConnectionsTree) LookupConn(saddr, daddr *syscall.SockaddrInet4) *event.UDPConn {
	t.Lock()
	ptr := t.tree.Lookup(event.Hash(saddr, daddr))
	t.Unlock()
	if ptr == nil {
		return nil
	}

	return (*event.UDPConn)(ptr)
}

type Client struct {
	tree  clientConnectionsTree
	epoll epoll.EPoll
}

func NewClient() *Client {
	return &Client{
		tree: clientConnectionsTree{tree: rbtree.NewRbtree()},
	}
}

func (c *Client) Run() error {
	var err error
	if err = c.epoll.Init(true); err != nil {
		return err
	}
	c.epoll.Callback.DoRead = c.DoRead
	c.epoll.Callback.DoWrite = c.DoWrite
	c.epoll.Callback.DoClose = c.DoClose
	c.epoll.Callback.DoProcessErr = c.DoProcessErr

	if err = c.request(); err != nil {
		return err
	}

	return c.epoll.WaitProcessLoop()
}

func (c *Client) DoRead(fd int) {
	fmt.Printf("DoRead fd=%d\n", fd)

	sa, err := syscall.Getsockname(fd)
	if err != nil {
		fmt.Printf("cannot get sockname fd=%d err=%s\n", fd, err)
		syscall.Close(fd)
		return
	}

	var saddr *syscall.SockaddrInet4
	saddr, err = event.Sockaddr4(sa)
	if err != nil {
		fmt.Printf("cannot get sockaddr4 addr=%v fd=%d err=%s\n", sa, fd, err)
		syscall.Close(fd)
		return
	}

	var (
		b     []byte = make([]byte, 32767)
		from  syscall.Sockaddr
		daddr *syscall.SockaddrInet4
		n     int
	)
	for {
		n, from, err = syscall.Recvfrom(fd, b, syscall.MSG_DONTWAIT)
		if err != nil {
			if err != syscall.EAGAIN {
				fmt.Printf("cannot read sockaddr4=%v fd=%d err=%s\n", saddr, fd, err)
			}
			return
		}

		daddr, err = event.Sockaddr4(from)
		if err != nil {
			fmt.Printf("cannot get sockaddr addr=%v fd=%d err=%s\n", from, fd, err)
			syscall.Close(fd)
			continue
		}

		c.response(saddr, daddr, b[:n])
	}
}

func (c *Client) DoWrite(fd int) {
	fmt.Printf("DoWrite\n")
}

func (c *Client) DoClose(fd int) {
	fmt.Printf("DoClose\n")
}

func (c *Client) DoProcessErr(fd int) {
	fmt.Printf("DoProcessErr\n")
}

func (c *Client) request() error {
	var err error

	conn := event.NewUDPConn()
	if err = conn.Dial(&syscall.SockaddrInet4{Addr: [4]byte{192, 168, 0, 17}, Port: 4443}); err != nil {
		return err
	}

	if err = c.epoll.Add(conn.Fd(), epoll.EPOLLIN); err != nil {
		conn.Close()
		return err
	}

	c.tree.InsertConn(conn)

	if err = conn.Sendto(conn.DAddr(), []byte("1234")); err != nil {
		c.tree.DeleteConn(conn)
		conn.Close()
		return err
	}

	fmt.Printf("fd=%d\n", conn.Fd())

	return nil
}

func (c *Client) response(saddr, daddr *syscall.SockaddrInet4, data []byte) {
	conn := c.tree.LookupConn(saddr, daddr)
	if conn == nil {
		fmt.Printf("no conn found saddr=%v daddr=%v", saddr, daddr)
		return
	}

	fmt.Printf("conn=%v\n", conn)
}
