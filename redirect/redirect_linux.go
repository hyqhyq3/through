package main

import (
	"encoding/binary"
	"errors"
	"net"
	"syscall"
	"unsafe"
)

//#include "netinet/in.h"
//#include "stdlib.h"
//#include "linux/netfilter_ipv4.h"
import "C"

func getOriginDst(c net.Conn) (addr net.Addr, err error) {
	switch conn := c.(type) {
	case *net.TCPConn:
		file, err := conn.File()
		if err != nil {
			return nil, err
		}

		var value C.struct_sockaddr_in
		var vlen = unsafe.Sizeof(value)
		_, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT,
			file.Fd(), syscall.SOL_IP, C.SO_ORIGINAL_DST,
			uintptr(unsafe.Pointer(&value)), uintptr(unsafe.Pointer(&vlen)), 0)
		if errno != 0 {
			return nil, errno
		}
		var ip [4]byte
		binary.LittleEndian.PutUint32(ip[:], uint32(value.sin_addr.s_addr))
		myaddr := net.TCPAddr{}
		myaddr.IP = net.IP(ip[:])
		myaddr.Port = int(C.ntohs(C.uint16_t(value.sin_port)))
		addr = &myaddr
	default:
		return nil, errors.New("only support tcp")
	}
	return
}
