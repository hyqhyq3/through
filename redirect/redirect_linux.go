// +build linux
// +build arm
package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

type sockAddrIn struct {
	family uint16
	port   [2]byte
	addr   [4]byte
	zero   [8]byte
}

func getOriginDst(c net.Conn) (addr net.Addr, err error) {
	switch conn := c.(type) {
	case *net.TCPConn:
		file, err := conn.File()
		if err != nil {
			return nil, err
		}

		var value sockAddrIn
		var vlen = unsafe.Sizeof(value)
		_, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT,
			file.Fd(), syscall.SOL_IP, 80,
			uintptr(unsafe.Pointer(&value)), uintptr(unsafe.Pointer(&vlen)), 0)
		if errno != 0 {
			return nil, errno
		}
		myaddr := net.TCPAddr{}
		myaddr.IP = net.IP(value.addr[:])
		myaddr.Port = int(binary.BigEndian.Uint16(value.port[:]))
		fmt.Println(value)
		addr = &myaddr
	default:
		return nil, errors.New("only support tcp")
	}
	return
}
