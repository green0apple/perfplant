package module

import "errors"

var (
	ErrMissingCallback = errors.New("missing callback function")
)

type Callback struct {
	DoAccept     func(fd int, ptr uintptr)
	DoRead       func(fd int, ptr uintptr)
	DoWrite      func(fd int, ptr uintptr)
	DoClose      func(fd int, ptr uintptr)
	DoProcessErr func(Fd int, ptr uintptr)
}

func (c *Callback) IsAllReady() bool {
	return !(c.DoAccept == nil || c.DoRead == nil || c.DoWrite == nil || c.DoClose == nil || c.DoProcessErr == nil)
}
