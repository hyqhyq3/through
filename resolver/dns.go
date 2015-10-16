package resolver

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

func Serve() error {
	conn, err := net.ListenPacket("udp", ":8024")
	if err != nil {
		return err
	}
	for {
		buf := make([]byte, 4096)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}
		fmt.Println(addr, buf[:n], string(buf[:n]))
		msg := new(dns.Msg)
		err = msg.Unpack(buf[:n])
		if err != nil {
			return err
		}

		dc, err := net.Dial("udp", "114.114.114.114:53")
		if err != nil {
			return err
		}
		_, err = dc.Write(buf[:n])
		if err != nil {
			return err
		}
		n, err = dc.Read(buf)
		if err != nil {
			return err
		}
		conn.WriteTo(buf[:n], addr)
		break
	}
	return nil
}
