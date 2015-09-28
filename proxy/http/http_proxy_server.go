package http_proxy

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/hyqhyq3/through/common"
	"github.com/hyqhyq3/through/proxy"
)

type request struct {
	Path    string
	Addr    string
	Headers map[string][]string
}

type HttpProxyServer struct {
	lsn net.Listener
}

func (h *HttpProxyServer) handleHttpProxyClient(c net.Conn, d common.Dialer) {
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

	req := &request{
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
	go proxy.Splice(c, rc, exit)
	go proxy.Splice(rc, c, exit)
	<-exit
	<-exit
	log.Println("connection closed")
}

func Listen(network, addr string) (h *HttpProxyServer, err error) {
	lsn, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	return &HttpProxyServer{
		lsn: lsn,
	}, nil
}

func (h *HttpProxyServer) Serve(exit <-chan bool, d common.Dialer) (err error) {
	conn := make(chan net.Conn)
	go func() {
		for {
			c, err := h.lsn.Accept()
			if err != nil {
				log.Println(err)
				break
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
			go h.handleHttpProxyClient(c, d)
		}
	}
	return
}

func ListenAndServe(addr string, d common.Dialer, exit <-chan bool) error {
	server, err := Listen("tcp", addr)
	if err != nil {
		return err
	}
	server.Serve(exit, d)
	return nil
}
