package common

import "net"

type Dialer interface {
	Dial(network, addr string) (c net.Conn, err error)
	Name() string
}
