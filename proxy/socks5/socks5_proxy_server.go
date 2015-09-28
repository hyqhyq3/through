package socks5_proxy

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/hyqhyq3/through/common"
	"github.com/hyqhyq3/through/proxy"
)

const SOCKS_VERSION byte = 0x5

const (
	_ = iota
	CONNECT
	BIND
	UDP_ASSOCIATE
)

const (
	IPv4   = 0x1
	Domain = 0x3
	IPv6   = 0x4
)

type Socks5ProxyServer struct {
	lsn net.Listener
}

func readVer(r *bufio.Reader) error {
	ver, err := r.ReadByte()
	if err != nil {
		return errors.New("cannot read ver")
	}
	if ver != SOCKS_VERSION {
		return errors.New("socks version error")
	}
	return nil
}

type request struct {
	Cmd  uint8
	Dst  string // ipv4 or ipv6 or domain
	Port int
}

func readRequest(r *bufio.Reader) (*request, error) {
	err := readVer(r)
	if err != nil {
		return nil, err
	}
	req := &request{}
	req.Cmd, err = r.ReadByte()
	if err != nil {
		return nil, err
	}

	r.ReadByte()

	atype, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch atype {
	case IPv4:
		ip := make([]byte, net.IPv4len)
		_, err := io.ReadFull(r, ip)
		if err != nil {
			return nil, err
		}
		req.Dst = net.IP(ip).String()
	case Domain:
		l, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if l > 0 {
			domain := make([]byte, l)
			_, err := io.ReadFull(r, domain)
			if err != nil {
				return nil, err
			}
			req.Dst = string(domain)
		}
	case IPv6:
		ip := make([]byte, net.IPv6len)
		_, err := io.ReadFull(r, ip)
		if err != nil {
			return nil, err
		}
		req.Dst = net.IP(ip).String()
	}

	port := make([]byte, 2)
	_, err = io.ReadFull(r, port)
	if err != nil {
		return nil, err
	}
	req.Port = int(binary.BigEndian.Uint16(port))
	return req, nil
}

func handleProxyClient(d common.Dialer, c net.Conn) {
	defer c.Close()

	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)

	err := readVer(r)
	if err != nil {
		log.Println(err)
		return
	}

	nmethods, err := r.ReadByte()
	if err != nil {
		log.Println("cannot read nmethods")
		return
	}

	if nmethods > 0 {
		methods := make([]byte, nmethods)
		_, err := io.ReadFull(r, methods)
		if err != nil {
			log.Println("cannot read methods")
			return
		}
	}

	//TODO auth
	w.Write([]byte{SOCKS_VERSION, 0x0})
	w.Flush()

	req, err := readRequest(r)
	if err != nil {
		log.Println("read request error")
		return
	}
	fmt.Println(req)
	if req.Cmd == CONNECT {
		rc, err := d.Dial("tcp", req.Dst+":"+strconv.Itoa(req.Port))
		if err != nil {
			return
		}
		defer rc.Close()
		w.Write([]byte{SOCKS_VERSION, 0, 0, 1})
		addr := rc.LocalAddr().(*net.TCPAddr)
		w.Write(addr.IP[:4])
		port := make([]byte, 2)
		binary.BigEndian.PutUint16(port, uint16(addr.Port))
		w.Write(port)
		w.Flush()

		exit := make(chan bool)
		go proxy.Splice(c, rc, exit)
		go proxy.Splice(rc, c, exit)
		<-exit
		<-exit
	}

}

func Listen(addr string) (*Socks5ProxyServer, error) {
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Socks5ProxyServer{lsn}, nil
}

func (s *Socks5ProxyServer) Serve(d common.Dialer, exit <-chan bool) (err error) {
	connChan := make(chan net.Conn)
	go func() {
		for {
			conn, err := s.lsn.Accept()
			if err != nil {
				log.Println(err)
				break
			}
			connChan <- conn
		}
	}()

M:
	for {
		select {
		case conn, ok := <-connChan:
			if !ok {
				s.lsn.Close()
				break M
			}
			go handleProxyClient(d, conn)

		case <-exit:
			s.lsn.Close()
			break M
		}
	}
	return nil
}

func ListenAndServeSocks5(addr string, d common.Dialer, exit <-chan bool) (err error) {
	server, err := Listen(addr)
	if err != nil {
		return err
	}

	return server.Serve(d, exit)
}
