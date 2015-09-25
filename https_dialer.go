package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
)

type AuthInfo struct {
	Username string
	Password string
}

type HTTPSDialer struct {
	Host string
	Port string
	Auth *AuthInfo
}

func CreateHTTPSProxy(u *url.URL) (dailer Dialer, err error) {
	auth := &AuthInfo{}
	auth.Username = u.User.Username()
	auth.Password, _ = u.User.Password()
	var host, port string
	host, port, err0 := net.SplitHostPort(u.Host)
	if err0 != nil {
		host = u.Host
		port = "443"
	}
	return &HTTPSDialer{
		Host: host + ":" + port,
		Auth: auth,
	}, nil
}

func init() {
	RegisterProxyType("https", CreateHTTPSProxy)
}

func (h *HTTPSDialer) Dial(network, addr string) (c net.Conn, err error) {
	if network != "tcp" {
		return nil, errors.New("https proxy only support tcp")
	}

	conn, err := tls.Dial("tcp", h.Host, nil)
	if err != nil {
		return
	}
	w := bufio.NewWriter(conn)
	w.WriteString("CONNECT %s HTTP/1.1\r\n")
	w.WriteString("Proxy-Agent: through/1.0\r\n")
	w.WriteString("\r\n")
	w.Flush()

	r := bufio.NewReader(conn)
	line, _, err := r.ReadLine()
	if err != nil {
		return
	}
	fmt.Println(string(line))

	return
}
