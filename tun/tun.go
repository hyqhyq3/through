package tun

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"
	"unsafe"
)

const (
	ICMP = 1
)

type IPFrame struct {
	Ver            uint8  // 4bits
	Length         uint16 // 16bits
	Identification uint16
	Flags          uint8  // 3bits. Do not fragment, more fragments
	FragmentOffset uint16 // 13bits.
	TTL            uint8  // 4bits
	Protocol       uint8  // 4bits
	Src            net.IP
	Dst            net.IP
	Payload        []byte
}

func readIPFrame(r io.Reader) (p *IPFrame, err error) {
	p = &IPFrame{}
	buf := make([]byte, 65535)
	_, err = r.Read(buf)
	if err != nil {
		return
	}
	p.Ver = buf[0] >> 4
	if p.Ver != 4 {
		log.Println("Only support ipv4")
		return nil, errors.New("Only support ipv4")
	}
	p.Length = binary.BigEndian.Uint16(buf[2:4])
	log.Println("packet length ", p.Length)

	p.Identification = binary.BigEndian.Uint16(buf[4:6])

	p.Flags = buf[6] >> 5

	p.FragmentOffset = (uint16(buf[6])&((1<<6)-1))<<8 + uint16(buf[7])

	p.TTL = buf[8]

	p.Protocol = buf[9]

	checksum0 := binary.BigEndian.Uint16(buf[10:12])
	buf[10] = 0
	buf[11] = 0

	var checksum uint32 = 0xffff
	for i := 0; i < 10; i++ {
		checksum += uint32(^binary.BigEndian.Uint16(buf[i*2 : i*2+2]))
		for checksum > 0xffff {
			checksum = (checksum >> 16) + (checksum & 0xffff)
		}

	}
	if uint16(checksum) != checksum0 {
		return nil, errors.New("Checksum error")
	}

	p.Src = buf[12:16]
	p.Dst = buf[16:20]

	p.Payload = buf[20:p.Length]
	return
}

const (
	Ping = 8
)

type ICMPHeader struct {
	Type     uint8
	Code     uint8
	Checksum uint16
}

type ICMPPing struct {
	ICMPHeader
	Identification uint16
	SequenceNumber uint16
}

func handleICMP(w io.Writer, frame *IPFrame) {
	bw := bufio.NewWriter(w)
	icmpType := frame.Payload[0]
	switch icmpType {
	case Ping:
		req := &ICMPPing{}
		copy((*[1 << 30]byte)(unsafe.Pointer(req))[:8], frame.Payload)
		ret := &ICMPPing{}
		ret.Type = 0
		ret.Code = 0
		ret.Identification = req.Identification
		ret.SequenceNumber = req.SequenceNumber
		bw.Write((*[1 << 30]byte)(unsafe.Pointer(ret))[:8])
		bw.Write(frame.Payload[8:])
		bw.Flush()
	}
}

func BringUp() {
	file, err := os.OpenFile("/dev/tun0", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	ifName := "tun0"
	_, err = createInterface(file.Fd(), ifName, syscall.IFF_UP|syscall.IFF_RUNNING)

	for {
		frame, err := readIPFrame(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if frame.Protocol == ICMP {
			handleICMP(file, frame)
		}
		fmt.Println(frame)
	}

}
