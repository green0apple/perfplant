package epoll

import (
	"encoding/binary"
	"errors"
	"perfplant/event/module"
	"syscall"
	"unsafe"
)

var (
	ErrEPollNotInitialized = errors.New("epoll does not be initialized")
)

const (
	EPOLLIN    = uint32(syscall.EPOLLIN)
	EPOLLOUT   = uint32(syscall.EPOLLOUT)
	EPOLLRDHUP = uint32(syscall.EPOLLRDHUP)
	EPOLLHUP   = uint32(syscall.EPOLLHUP)

	EPOLLET = (1 << 31)
)

type EPoll struct {
	Callback module.Callback
	fd       int
	edgeMode bool
}

func (e *EPoll) Init(edgeMode bool) error {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}
	e.fd = fd
	e.edgeMode = edgeMode

	return nil
}

func (e *EPoll) Close() {
	if e.fd > 0 {
		syscall.Close(e.fd)
	}
}

func (e *EPoll) Add(fd int, events uint32, data any) error {
	var (
		ev  EpollEvent
		ed  EpollData
		ptr uintptr
	)

	ed.Fd = int32(fd)
	ed.Data = data

	ptr = uintptr(unsafe.Pointer(&ed))
	binary.BigEndian.PutUint64(ev.Ptr[:], uint64(ptr))

	ev.Events = events
	if e.edgeMode {
		ev.Events |= EPOLLET
	}

	return EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &ev)
}

func (e *EPoll) WaitEvent(events []EpollEvent) (int, error) {
	return EpollWait(e.fd, events, 0)
}

func (e *EPoll) WaitProcessLoop() error {
	return nil
	// if e.fd <= 0 {
	// 	return ErrEPollNotInitialized
	// }

	// if !e.Callback.IsAllReady() {
	// 	return module.ErrMissingCallback
	// }

	// var (
	// 	events       []EpollEvent = make([]EpollEvent, 64)
	// 	data         *EpollData
	// 	isNewConn    bool
	// 	i, count, fd int
	// 	err          error
	// )
	// for {
	// 	count, err = e.WaitEvent(events)
	// 	if err != nil {
	// 		fmt.Printf("err : %s", err)
	// 		return err
	// 	}

	// 	for i = 0; i < count; i++ {
	// 		data = (*EpollData)(events[i].Ptr)
	// 		fd = int(data.Fd)
	// 		if data.Ptr == 0 {
	// 			isNewConn = true
	// 		} else {
	// 			isNewConn = false
	// 		}

	// 		if isNewConn {
	// 			fmt.Printf("new connection!!\n")
	// 		}

	// 		if events[i].Events&syscall.EPOLLERR == 1 {
	// 			e.Callback.DoProcessErr(fd, data.Ptr)
	// 			continue // no need to process error connection
	// 		}

	// 		if events[i].Events&syscall.EPOLLIN == 1 {
	// 			if isNewConn {
	// 				e.Callback.DoAccept(fd, data.Ptr)
	// 			} else {
	// 				e.Callback.DoRead(fd, data.Ptr)
	// 			}
	// 		}

	// 		if events[i].Events&syscall.EPOLLOUT == 1 {
	// 			e.Callback.DoWrite(fd, data.Ptr)
	// 		}

	// 		if events[i].Events&(syscall.EPOLLHUP|syscall.EPOLLRDHUP) == 1 {
	// 			e.Callback.DoClose(fd, data.Ptr)
	// 		}
	// 	}
	// }
}
