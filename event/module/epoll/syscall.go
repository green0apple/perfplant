package epoll

import (
	"syscall"
	"unsafe"
)

var (
	_zero uintptr
)

// Do the interface allocations only once for common
// Errno values.
var (
	errEAGAIN error = syscall.EAGAIN
	errEINVAL error = syscall.EINVAL
	errENOENT error = syscall.ENOENT
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case syscall.EAGAIN:
		return errEAGAIN
	case syscall.EINVAL:
		return errEINVAL
	case syscall.ENOENT:
		return errENOENT
	}
	return e
}

/*
   struct epoll_event {
        uint32_t      events;
        epoll_data_t  data;
    };

    union epoll_data {
        void     *ptr;
        int       fd;
        uint32_t  u32;
        uint64_t  u64;
    };

*/

type EpollData struct {
	Data any
	Fd   int32
}

type EpollEvent struct {
	Events uint32
	Ptr    [8]byte
}

func EpollWait(epfd int, events []EpollEvent, msec int) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(events) > 0 {
		_p0 = unsafe.Pointer(&events[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}

	r0, _, e1 := syscall.Syscall6(syscall.SYS_EPOLL_WAIT, uintptr(epfd), uintptr(_p0), uintptr(len(events)), uintptr(msec), 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}

	return
}

func EpollCtl(epfd int, op int, fd int, event *EpollEvent) (err error) {
	_, _, e1 := syscall.RawSyscall6(syscall.SYS_EPOLL_CTL, uintptr(epfd), uintptr(op), uintptr(fd), uintptr(unsafe.Pointer(event)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}
