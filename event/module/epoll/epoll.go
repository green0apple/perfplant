package epoll

import (
	"errors"
	"fmt"
	"perfplant/event/module"
	"syscall"
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

func (e *EPoll) Add(fd int, events uint32) error {
	var ev syscall.EpollEvent
	ev.Fd = int32(fd)
	ev.Events = events
	if e.edgeMode {
		ev.Events |= EPOLLET
	}

	return syscall.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &ev)
}

func (e *EPoll) WaitEvent(events []syscall.EpollEvent) (int, error) {
	return syscall.EpollWait(e.fd, events, 0)
}

func (e *EPoll) WaitProcessLoop() error {
	if e.fd <= 0 {
		return ErrEPollNotInitialized
	}

	if !e.Callback.IsAllReady() {
		return module.ErrMissingCallback
	}

	var (
		events       []syscall.EpollEvent = make([]syscall.EpollEvent, 64)
		i, count, fd int
		err          error
	)
	for {
		count, err = e.WaitEvent(events)
		if err != nil {
			fmt.Printf("err : %s", err)
			return err
		}

		for i = 0; i < count; i++ {
			fd = int(events[i].Fd)

			if events[i].Events&syscall.EPOLLERR == 1 {
				e.Callback.DoProcessErr(fd)
				continue // no need to process error connection
			}

			if events[i].Events&syscall.EPOLLIN == 1 {
				e.Callback.DoRead(fd)
			}

			if events[i].Events&syscall.EPOLLOUT == 1 {
				e.Callback.DoWrite(fd)
			}

			if events[i].Events&(syscall.EPOLLHUP|syscall.EPOLLRDHUP) == 1 {
				e.Callback.DoClose(fd)
			}
		}
	}
}
