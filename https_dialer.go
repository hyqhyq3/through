package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/url"
)

type Dialer interface {
	Dial(network, addr string) (c net.Conn, err error)
	Name() string
}

type AuthInfo struct {
	Username string
	Password string
}

type HTTPSDialer struct {
	name string
	Host string
	Auth *AuthInfo
}

func CreateHTTPSProxy(name string, u *url.URL) (dailer Dialer, err error) {
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
		name: name,
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
	w.WriteString(fmt.Sprintf("CONNECT %s HTTP/1.0\r\n", addr))
	w.WriteString(fmt.Sprintf("Host: %s", addr))
	w.WriteString("Proxy-Agent: through/1.0\r\n")
	if h.hasAuth() {
		w.WriteString(fmt.Sprintf("Proxy-Authorization: %s\r\n", h.basicAuth()))
	}
	w.WriteString("\r\n")
	w.Flush()

	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadString(byte('\n'))
		if err != nil {
			return nil, err
		}
		if line == "\r\n" {
			break
		}
	}
	return conn, nil
}

func (h *HTTPSDialer) hasAuth() bool {
	return h.Auth != nil
}

func (h *HTTPSDialer) basicAuth() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(h.Auth.Username+":"+h.Auth.Password)))
}

func (h *HTTPSDialer) Name() string {
	return h.name
}
