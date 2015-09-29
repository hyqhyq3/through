package tun

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func BringUp() {
	file, err := os.OpenFile("/dev/tun0", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	ifName := "tun0"
	_, err = createInterface(file.Fd(), ifName, syscall.IFF_UP|syscall.IFF_RUNNING)

	for {
		buf := make([]byte, 100)
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		fmt.Println(buf[:n])
	}

}

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func createInterface(fd uintptr, ifName string, flags uint16) (createdIFName string, err error) {
	var req ifReq
	req.Flags = flags
	copy(req.Name[:], ifName)

	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s), uintptr(syscall.SIOCSIFFLAGS), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		err = errno
		return
	}
	createdIFName = strings.Trim(string(req.Name[:]), "\x00")
	return
}
