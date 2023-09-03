package main

import (
	"encoding/binary"
	"fmt"
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
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_IP)
	if err != nil {
		panic(err)
	}
	fmt.Printf("listen fd=%d\n", fd)
	if err = syscall.SetNonblock(fd, true); err != nil {
		panic(err)
	}

	if err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: 80}); err != nil {
		panic(err)
	}

	if err = syscall.Listen(fd, 1024); err != nil {
		panic(err)
	}

	e := epoll.EPoll{}
	if err = e.Init(true); err != nil {
		panic(err)
	}

	conn := Connection{Fd: fd, Typ: CONN_TYPE_LISTENER}

	if err = e.Add(fd, epoll.EPOLLIN, conn); err != nil {
		panic(err)
	}

	var (
		data     epoll.EpollData
		count, i int
		ptr      uintptr
		events   []epoll.EpollEvent = make([]epoll.EpollEvent, 10)
	)

	for {
		count, err = e.WaitEvent(events)
		if err != nil {
			fmt.Printf("err = %s\n", err)
		}

		for i = 0; i < count; i++ {
			ptr = uintptr(binary.BigEndian.Uint64(events[i].Ptr[:]))
			data = *(*epoll.EpollData)(unsafe.Pointer(ptr))

			fmt.Printf("fd is : %d\n", data.Fd)
			fmt.Printf("data is : %v\n", data.Data)

			// fmt.Printf("sizeof=%d\n", unsafe.Sizeof(events[i].Ptr))
			// fmt.Printf("after : %v\n", (*epoll.EpollData)(unsafe.Pointer(events[i].Ptr)))
			// (*(*Event)(events[ev].Ptr)
			// fmt.Printf("%v\n", (*(*epoll.EpollEvent)(events[i].Ptr)))

			// conn := (**(**epoll.Data)(unsafe.	Pointer(&events[i].Fd))).(Connection)
			// fmt.Printf("fd=%d, events=0x%x, conn.fd=%d, conn.typ=%d\n", events[i].Fd, events[i].Events, conn.Fd, conn.Typ)
			// data := (*epoll.EpollData)(events[i].Ptr)
			// // fmt.Printf("fd=%d, events=0x%x\n", data.Fd, data.Ptr)
			// fmt.Printf("%v\n", data)
		}
	}
}
