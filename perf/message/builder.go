package message

import "syscall"

const (
	REQUEST_MESSAGE_TYPE_FIXED = iota + 1
	REQUEST_MESSAGE_TYPE_RANDOM
)

const (
	RESPONSE_MESSAGE_TYPE_MATCHED = iota + 1
	RESPONSE_MESSAGE_TYPE_INCLUDED
)

type Message struct {
	From         syscall.SockaddrInet4
	To           syscall.SockaddrInet4
	Request      []byte
	Response     []byte
	ResponseType int
}

type MessageBuilder func() (Message, error)
