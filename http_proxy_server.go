// http_proxy_server.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"regexp"
	"strings"
)

func Splice(r io.Reader, w io.WriteCloser, exit chan<- bool) {
	buf := make([]byte, 1024)
	for {
		rn, err := r.Read(buf)
		if rn > 0 {
			_, err := w.Write(buf[:rn])
			if err != nil {
				w.Close()
				break
			}
		}
		if err != nil {
			w.Close()
			break
		}
	}
	exit <- true
}

type Request struct {
	Path    string
	Addr    string
	Headers map[string][]string
}

type HttpProxyServer struct {
	lsn net.Listener
}

func NewHttpProxyServer() *HttpProxyServer {
	return &HttpProxyServer{}
}

func handleProxyClient(c net.Conn) {
	var closeConn net.Conn = c
	defer func() {
		if closeConn != nil {
			closeConn.Close()
		}
	}()

	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)
	defer w.Flush()
	requestLine, err := r.ReadString(byte('\n'))
	if err != nil {
		log.Println(err)
		return
	}

	req := &Request{
		Headers: make(map[string][]string),
	}
	for {
		line, err := r.ReadString(byte('\n'))
		if err != nil {
			log.Println(err)
			return
		}
		if len(line) <= 2 {
			break
		}
		line = strings.TrimSpace(line)
		arr := strings.SplitN(line, ":", 2)
		if len(arr) == 2 {
			req.Headers[arr[0]] = append(req.Headers[arr[0]], strings.TrimSpace(arr[1]))
		}
	}
	// todo readAuth
	re, err := regexp.Compile(`([A-Z]+) ([^\s]*) (HTTP/\d\.\d)`)
	if err != nil {
		log.Fatal(err)
	}
	m := re.FindStringSubmatch(requestLine)
	if len(m) < 3 {
		w.WriteString("HTTP/1.1 404 Host Not Found\r\n\r\nHost not found")
		return
	}
	req.Path = m[2]

	var rc net.Conn
	switch m[1] {
	case "CONNECT":
		req.Addr = req.Path
		host, _, err := net.SplitHostPort(req.Addr)
		if err != nil {
			log.Println("cannot split host " + req.Addr)
			return
		}
		d := Route(host)
		if d == nil {
			log.Println("cannot route to " + req.Addr)
			return
		}

		rc, err = d.Dial("tcp", req.Addr)
		if err != nil {
			w.WriteString("HTTP/1.1 500 Server Error\r\n\r\nCannot connect to " + req.Addr)
			return
		}
		w.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n")
		w.Flush()
	default:
		if host, ok := req.Headers["Host"]; ok {
			req.Addr = host[0]
			if strings.IndexByte(req.Addr, byte(':')) == -1 {
				req.Addr = req.Addr + ":80"
			}
			host, _, _ := net.SplitHostPort(req.Addr)
			d := Route(host)
			rc, err = d.Dial("tcp", req.Addr)
			if err != nil {
				w.WriteString("HTTP/1.1 500 Cannot Connect to Host\r\n\r\n")
				return
			}
			wrc := bufio.NewWriter(rc)
			u, _ := url.Parse(m[2])
			wrc.WriteString(fmt.Sprintf("%s %s %s\r\n", m[1], u.Path, m[3]))
			for key, headers := range req.Headers {
				if len(key) > 6 && key[:6] == "Proxy-" {
					continue
				}
				for _, header := range headers {
					wrc.WriteString(key + ": " + header + "\r\n")
				}
			}
			wrc.WriteString("\r\n")
			wrc.Flush()

		} else {
			w.WriteString("HTTP/1.1 400 Bad Request\r\n\r\n")
			return
		}

	}

	exit := make(chan bool)
	go Splice(c, rc, exit)
	go Splice(rc, c, exit)
	<-exit
	<-exit
	log.Println("connection closed")
}

func (h *HttpProxyServer) ListenAndServe(addr string, exit <-chan bool) (err error) {
	h.lsn, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	conn := make(chan net.Conn)
	go func() {
		for {
			c, err := h.lsn.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			conn <- c
		}
	}()
M:
	for {
		select {
		case <-exit:
			h.lsn.Close()
			break M
		case c, ok := <-conn:
			if !ok {
				h.lsn.Close()
				break M
			}
			go handleProxyClient(c)
		}
	}
	return
}
