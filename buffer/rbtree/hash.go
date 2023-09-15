package rbtree

import (
	"hash/crc32"
)

func Hash(b ...byte) uint32 {
	return crc32.ChecksumIEEE(b)
}
