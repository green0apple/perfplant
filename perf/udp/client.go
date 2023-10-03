package udp

import (
	"fmt"
	"perfplant/buffer/rbtree"
	"perfplant/event"
	"perfplant/event/module/epoll"
	"perfplant/perf/message"
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
	t.tree.Insert(uint32(conn.Fd()), unsafe.Pointer(conn))
}

func (t *clientConnectionsTree) DeleteConn(conn *event.UDPConn) {
	t.Lock()
	defer t.Unlock()
	t.tree.Delete(uint32(conn.Fd()))
}

func (t *clientConnectionsTree) LookupConn(fd int32) *event.UDPConn {
	t.Lock()
	ptr := t.tree.Lookup(uint32(fd))
	t.Unlock()
	if ptr == nil {
		return nil
	}

	return (*event.UDPConn)(ptr)
}

type Client struct {
	messageBuilder message.MessageBuilder

	tree  clientConnectionsTree
	epoll epoll.EPoll
}

func NewClient() *Client {
	return &Client{
		tree: clientConnectionsTree{tree: rbtree.NewRbtree()},
	}
}

func (c *Client) Run(builder message.MessageBuilder) error {
	var err error
	if err = c.epoll.Init(true); err != nil {
		return err
	}

	c.epoll.Callback.DoRead = c.DoRead
	c.epoll.Callback.DoWrite = c.DoWrite
	c.epoll.Callback.DoClose = c.DoClose
	c.epoll.Callback.DoProcessErr = c.DoProcessErr

	c.messageBuilder = builder

	if err = c.request(); err != nil {
		return err
	}

	return c.epoll.WaitProcessLoop()
}

func (c *Client) DoRead(fd int) {
	fmt.Printf("DoRead fd=%d\n", fd)

	conn := c.tree.LookupConn(int32(fd))
	if conn == nil {
		fmt.Printf("no conn found fd=%d\n", fd)
		syscall.Close(fd)
		return
	}

	var (
		b   []byte
		err error
	)
	for {
		_, b, err = conn.Recvmsg()
		if err != nil {
			if err != syscall.EAGAIN {
				c.CloseConn(conn)
				fmt.Printf("cannot read %s->%s fd=%d err=%s\n", conn.DAddrString(), conn.SAddrString(), fd, err)
			}
			return
		}

		c.response(conn, b)
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

func (c *Client) CloseConn(conn *event.UDPConn) {
	c.tree.DeleteConn(conn)
	conn.Close()
}

func (c *Client) request() error {
	msg, err := c.messageBuilder()
	if err != nil {
		return err
	}

	conn := event.NewUDPConn()
	if err = conn.Dial(&syscall.SockaddrInet4{Addr: [4]byte{127, 0, 0, 1}, Port: 4443}); err != nil {
		return err
	}

	if err = c.epoll.Add(int(conn.Fd()), epoll.EPOLLIN); err != nil {
		conn.Close()
		return err
	}

	c.tree.InsertConn(conn)

	if err = conn.Sendto(conn.DAddr(), msg.Request); err != nil {
		c.CloseConn(conn)
		return err
	}

	fmt.Printf("fd=%d\n", conn.Fd())

	return nil
}

func (c *Client) response(conn *event.UDPConn, data []byte) {
	fmt.Printf("%s->%s read %s\n", conn.DAddrString(), conn.SAddrString(), string(data))
}
