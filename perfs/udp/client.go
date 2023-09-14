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

	c.epoll.Callback.DoAccept = c.DoAccept
	c.epoll.Callback.DoRead = c.DoRead
	c.epoll.Callback.DoWrite = c.DoWrite
	c.epoll.Callback.DoClose = c.DoClose
	c.epoll.Callback.DoProcessErr = c.DoProcessErr

	if err = c.request(); err != nil {
		return err
	}

	return c.epoll.WaitProcessLoop()
}

func (c *Client) DoAccept(fd int, ptr uintptr) {
	fmt.Printf("DoAccept\n")
}

func (c *Client) DoRead(fd int, ptr uintptr) {
	fmt.Printf("DoRead\n")
}

func (c *Client) DoWrite(fd int, ptr uintptr) {
	fmt.Printf("DoWrite\n")
}

func (c *Client) DoClose(fd int, ptr uintptr) {
	fmt.Printf("DoClose\n")
}

func (c *Client) DoProcessErr(fd int, ptr uintptr) {
	fmt.Printf("DoProcessErr\n")
}

func (c *Client) request() error {
	var err error

	conn := event.NewUDPConn()
	if err = conn.Dial(&syscall.SockaddrInet4{Addr: [4]byte{192, 168, 0, 17}, Port: 4443}); err != nil {
		return err
	}

	if err = c.epoll.Add(conn.Fd(), epoll.EPOLLIN, uintptr(unsafe.Pointer(&conn))); err != nil {
		conn.Close()
		return err
	}

	if err = conn.Sendto(conn.DAddr(), []byte("1234")); err != nil {
		conn.Close()
		return err
	}

	return nil
}
