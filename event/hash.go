package event

import (
	"hash/crc32"
	"syscall"
)

// only 0~65535
func PortLittleEndian(port int) []byte {
	return []byte{byte(port & 255), byte(port >> 8 & 255)}
}

func FdLittleEndian(fd int32) []byte {
	return []byte{byte(fd & 255), byte(fd >> 8 & 255), byte(fd >> 16 & 255), byte(fd >> 32 & 255)}
}

func HashAddr(saddr, daddr *syscall.SockaddrInet4) uint32 {
	var (
		b     []byte
		addrs [2]*syscall.SockaddrInet4
	)

	if saddr != nil && daddr != nil {
		if saddr.Port < daddr.Port {
			addrs[0] = saddr
			addrs[1] = daddr
		} else {
			addrs[0] = daddr
			addrs[1] = saddr
		}
	} else {
		if saddr != nil {
			addrs[0] = saddr
		}

		if daddr != nil {
			addrs[1] = daddr
		}
	}

	for _, a := range addrs {
		if a == nil {
			continue
		}

		b = append(b, a.Addr[:]...)
		b = append(b, PortLittleEndian(a.Port)...)
	}

	return crc32.ChecksumIEEE(b)
}

func HashFd(fd int32) uint32 {
	return crc32.ChecksumIEEE(FdLittleEndian(fd))
}
