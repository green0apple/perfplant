package module

import (
	"errors"
)

var (
	ErrMissingCallback = errors.New("missing callback function")
)

type Callback struct {
	DoRead       func(fd int)
	DoWrite      func(fd int)
	DoClose      func(fd int)
	DoProcessErr func(fd int)
}

func (c *Callback) IsAllReady() bool {
	return !(c.DoRead == nil || c.DoWrite == nil || c.DoClose == nil || c.DoProcessErr == nil)
}
