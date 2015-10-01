package tun

import (
	"encoding/binary"
	"net"
	"strings"
	"syscall"
	"unsafe"
)

/*
#include <net/if.h>
#include <netinet/in.h>
#include <sys/sockio.h>
#include <memory.h>
#include <sys/ioctl.h>
#include <errno.h>
#include <stdlib.h>

#define SIN(x) ((struct sockaddr_in *)&(x))
#define SAD(x) ((struct sockaddr *)&(x))

static void
setsin (struct sockaddr_in *sa, int family, u_long addr)
{
    bzero (sa, sizeof (*sa));
    sa->sin_len = sizeof (*sa);
    sa->sin_family = family;
    sa->sin_addr.s_addr = addr;
}

int setInterfaceIP(char* ifName, unsigned long local, unsigned long remote)
{
	int s = socket(AF_INET, SOCK_DGRAM, 0);
	struct ifaliasreq req;
	memset(&req, 0, sizeof(req));
	strcpy(req.ifra_name, ifName);


	setsin(SIN(req.ifra_addr), AF_INET, htonl(local));
	setsin(SIN(req.ifra_broadaddr), AF_INET, htonl(remote));
	setsin(SIN(req.ifra_mask), AF_INET, htonl(0xffffff00));
	if(ioctl(s, SIOCAIFADDR, &req) == -1) {
		return errno;
	}

//	if(ioctl(s, SIOCSIFDSTADDR, &req) == -1) {
//		return errno;
//	}

	return 0;
}
*/
import "C"

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

func setInterfaceIP(ifName string, local, remote net.IP) (err error) {
	cIfName := C.CString(ifName)
	cLocal := C.ulong(binary.BigEndian.Uint32(local.To4()))
	cRemote := C.ulong(binary.BigEndian.Uint32(remote.To4()))
	errno := C.setInterfaceIP(cIfName, cLocal, cRemote)
	C.free(unsafe.Pointer(cIfName))
	if errno != 0 {
		err = syscall.Errno(errno)
	}
	return
}
