package main

import (
	"errors"
	"net"
	"net/url"

	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

type SSDialer struct {
	name     string
	Host     string
	Method   string
	Password string
}

func NewSSDialer(name string, u *url.URL) (Dialer, error) {
	host := u.Host
	if u.User == nil {
		return nil, errors.New("Must have encrypt method and password")
	}
	method := u.User.Username()
	password, _ := u.User.Password()
	return &SSDialer{name, host, method, password}, nil
}

func init() {
	RegisterProxyType("ss", NewSSDialer)
}

func (s *SSDialer) Dial(network, addr string) (c net.Conn, err error) {
	cipher, err := ss.NewCipher(s.Method, s.Password)
	if err != nil {
		return nil, err
	}
	ssConn, err := ss.Dial(addr, s.Host, cipher)
	if err != nil {
		return nil, err
	}
	return ssConn, nil
}

func (s *SSDialer) Name() string {
	return s.name
}
