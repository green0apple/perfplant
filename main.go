package main

import (
	"encoding/binary"
	"fmt"
	"perfplant/event"
	"perfplant/event/module/epoll"
	"syscall"
	"unsafe"
)

const (
	CONN_TYPE_LISTENER = iota + 1
	CONN_TYPE_TCP
	CONN_TYPE_UDP
)

type Connection struct {
	Fd  int
	Typ int
}

func main() {
	listener, err := event.ListenUDP(syscall.SockaddrInet4{Port: 80}, 1024, event.LISTEN_OPT_NONBLOCK)
	if err != nil {
		panic(err)
	}

	e := epoll.EPoll{}
	if err = e.Init(true); err != nil {
		panic(err)
	}

	if err = e.Add(listener.Fd(), epoll.EPOLLIN, &listener); err != nil {
		panic(err)
	}

	var (
		data     epoll.EpollData
		count, i int
		ptr      uintptr
		events   []epoll.EpollEvent = make([]epoll.EpollEvent, 10)

		n    int
		from syscall.SockaddrInet4
	)

	for {
		count, err = e.WaitEvent(events)
		if err != nil {
			fmt.Printf("err = %s\n", err)
		}

		for i = 0; i < count; i++ {
			ptr = uintptr(binary.BigEndian.Uint64(events[i].Ptr[:]))
			data = *(*epoll.EpollData)(unsafe.Pointer(ptr))

			n, from, err = syscall.Recvfrom(int(data.Fd))

			fmt.Printf("fd is : %d\n", data.Fd)
			fmt.Printf("data is : %v\n", data.Data)
		}
	}
}
