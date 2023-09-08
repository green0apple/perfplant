package rbtree

import (
	"hash/crc32"
)

func Hash(b ...byte) uint {
	return uint(crc32.ChecksumIEEE(b))
}

// only 0~65535
func PortLittleEndian(port int) []byte {
	return []byte{byte(port & 255), byte(port >> 8 & 255)}
}
